package usecase

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/felipeksw/goexpert-fullcycle-rate-limiter/internal/infra/repository"
)

func RateLimitByKey(repo repository.RateLimiterStorage, key string, requestPerSecond int64, blockingTime int64) (bool, error) {

	slog.Debug("[RateLimitByKey]", "msg", "entrou no RateLimitByKey")
	slog.Debug("[RateLimitByKey]", "msg", repo.Debug("repositorio"))

	// Recupera o contador da chave no cache
	cntS, err := repo.Get(key + ":cnt")
	if err != nil {
		slog.Error("[RateLimitByKey]", "GET", "cnt", "key", key, "msg", err.Error())
		return false, err
	}
	slog.Debug("[RateLimitByKey]", "cntS", cntS, "key", key)

	// Se não localizar o contador, adiciona a chave um contador e timestamp, contando o acesso atual
	if cntS == "" {
		slog.Debug("[RateLimitByKey]", "mgs", "entrou para criar o controle da key")

		_, err := repo.Set(key+":cnt", 1, 0)
		if err != nil {
			slog.Error("[RateLimitByKey]", "SET", "cnt:1", "key", key, "msg", err.Error())
			return false, err
		}
		_, err = repo.Set(key+":timestamp", time.Now().UnixMilli(), 0)
		if err != nil {
			slog.Error("[RateLimitByKey]", "SET", "timestamp", "key", key, "msg", err.Error())
			return false, err
		}
		return false, nil
	}

	// Converte o contador para inteiro
	cnt, err := strconv.ParseInt(cntS, 10, 64)
	if err != nil {
		slog.Error("[RateLimitByKey]", "cntS", cntS, "key", key, "msg", err.Error())
		return false, err
	}

	// Verifica se a chave já está bloqueada
	// Se sim:
	//  retorna TRUE informando que o limite de acessos por segundo foi atingido
	if cnt < 1 {
		slog.Debug("[RateLimitByKey]", "mgs", "chave já bloqueada")
		return true, nil
	}

	// Recupera o timestamp da chave no cahe
	tspS, err := repo.Get(key + ":timestamp")
	if err != nil {
		slog.Error("[RateLimitByKey]", "GET", "timestamp", "key", key, "msg", err.Error())
		return false, err
	}

	// Se não localizar o timestamp, descarta a chave
	if tspS == "" {
		slog.Error("[RateLimitByKey]", "tspS", "", "key", key, "msg", "key discarded: could not find timestamp referrence for key ["+key+"] in the cache")
		_, err := repo.Del(key + ":cnt")
		if err != nil {
			slog.Error("[RateLimitByKey]", "DEL", "cnt", "key", key, "msg", err.Error())
			return false, err
		}
		return false, nil
	}

	// Converte o timestamp para inteiro
	tsp, err := strconv.ParseInt(tspS, 10, 64)
	if err != nil {
		slog.Error("[RateLimitByKey]", "tspS", tspS, "key", key, "msg", err.Error())
		return false, err
	}

	// Verifica "agora" está dentor de 1 segundo desde o primeiro acesso da chave
	if time.Now().UnixMilli()-tsp <= 1000 {
		slog.Debug("[RateLimitByKey]", "Now", time.Now().UnixMilli(), "KeyTimestamp", tsp)

		// Verifica se o contador, atingiu o limite configurado
		// Se sim:
		//   bloqueia a chave pelo tempo de bloqueio configurado
		//   retorna TRUE informando que o limite de acessos por segundo foi atingido
		if cnt+1 >= requestPerSecond {
			slog.Debug("[RateLimitByKey]", "KeyCounter", cnt+1, "RequestPerSecond", requestPerSecond)

			_, err := repo.Set(key+":cnt", -1, int32(blockingTime))
			if err != nil {
				slog.Error("[RateLimitByKey]", "SET", "cnt:-1", "key", key, "msg", err.Error())
				return false, err
			}
			_, err = repo.Set(key+":timestamp", tsp, int32(blockingTime))
			if err != nil {
				slog.Error("[RateLimitByKey]", "SET", "timestamp", "key", key, "msg", err.Error())
				return false, err
			}
			return true, nil
		}

		// Incrementa o contador da chave
		_, err = repo.Increment(key+":cnt", 0)
		if err != nil {
			slog.Error("[RateLimitByKey]", "INCR", "", "key", key, "msg", err.Error())
			return false, err
		}
		return false, nil
	}

	// Apaga o contador da chave
	_, err = repo.Del(key + ":cnt")
	if err != nil {
		slog.Error("[RateLimitByKey]", "DEL", "cnt", "key", key, "msg", err.Error())
		return false, err
	}

	// Apaga o timestamp da chave
	_, err = repo.Del(key + ":timestamp")
	if err != nil {
		slog.Error("[RateLimitByKey]", "DEL", "timestamp", "key", key, "msg", err.Error())
		return false, err
	}

	return false, nil
}
