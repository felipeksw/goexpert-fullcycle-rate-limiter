package usecase_test

import (
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/felipeksw/goexpert-fullcycle-rate-limiter/internal/infra/repository"
	"github.com/felipeksw/goexpert-fullcycle-rate-limiter/internal/usecase"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/rand"
)

func TestRateLimitByKey(t *testing.T) {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	redisRepo := repository.NewRedisCahe(
		redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   3,
		}),
	)
	var requestPerSecond int64 = 3
	var blockingTime int64 = 2

	type uuidKey struct {
		key        string
		count      int64
		timestsamp int64
		ttl        int64
	}
	uuidKeyList := [10]uuidKey{}

	for i := 0; i < 10; i++ {
		slog.Info("[]", "i", i)
		uuidKeyList[i].key = ""
		uuidKeyList[i].count = 0
		uuidKeyList[i].timestsamp = 0
		uuidKeyList[i].ttl = 0
	}

	for i := 0; i < 100; i++ {

		slog.Debug("-------")

		randKey := rand.Intn(9)

		// Reinicia a posição do slice para reutilização
		if uuidKeyList[randKey].ttl > 0 && time.Now().UnixMilli() > uuidKeyList[randKey].ttl {
			slog.Debug("[TEST]", "msg", "clear por ttl", "key", uuidKeyList[randKey].key, "ttl", uuidKeyList[randKey].ttl, "now", time.Now().UnixMilli())
			uuidKeyList[randKey].key = ""
			uuidKeyList[randKey].count = 0
			uuidKeyList[randKey].timestsamp = 0
			uuidKeyList[randKey].ttl = 0
		}

		// Adiciona ao slice a chave, caso a posição esteja vazia
		if uuidKeyList[randKey].key == "" {
			slog.Debug("[TEST]", "randKey", randKey, "msg", "posição vazia")
			uuidKeyList[randKey].key = uuid.New().String()
			uuidKeyList[randKey].count = 0
			uuidKeyList[randKey].timestsamp = time.Now().UnixMilli()
			uuidKeyList[randKey].ttl = 0
		}

		slog.Debug("[TEST]", "key", uuidKeyList[randKey].key)
		slog.Debug("[TEST]", "count", uuidKeyList[randKey].count)
		slog.Debug("[TEST]", "timestsamp", uuidKeyList[randKey].timestsamp)
		slog.Debug("[TEST]", "ttl", uuidKeyList[randKey].ttl)

		sts, err := usecase.RateLimitByKey(redisRepo, uuidKeyList[randKey].key, requestPerSecond, blockingTime)
		assert.Nil(t, err)

		slog.Debug("[TEST]", "msg", "NOW", "timestsamp", time.Now().UnixMilli())

		if time.Now().UnixMilli()-uuidKeyList[randKey].timestsamp <= 1000 {
			slog.Debug("[TEST]", "msg", "entrou dentro de 1 segundo", "requestPerSecond", requestPerSecond, "uuidKeyList[randKey].count+1", uuidKeyList[randKey].count+1)
			if uuidKeyList[randKey].count+1 >= requestPerSecond {
				uuidKeyList[randKey].count = -1
				uuidKeyList[randKey].ttl = time.Now().UnixMilli() + (blockingTime * int64(time.Duration(time.Millisecond)))
				slog.Debug("[TEST]", "msg", "ttl configurado", "ttl", uuidKeyList[randKey].ttl, "count", uuidKeyList[randKey].count)
				continue
			}

			if uuidKeyList[randKey].count >= 0 {
				uuidKeyList[randKey].count = uuidKeyList[randKey].count + 1
				slog.Debug("[TEST]", "msg", "aumento do contador", "key", uuidKeyList[randKey].key, "count", uuidKeyList[randKey].count)
			}

		} else {
			slog.Debug("[TEST]", "msg", "clear por passar de 1 segundo entre as chamadas", "key", uuidKeyList[randKey].key)
			uuidKeyList[randKey].key = ""
			uuidKeyList[randKey].count = 0
			uuidKeyList[randKey].timestsamp = 0
			uuidKeyList[randKey].ttl = 0
		}

		if uuidKeyList[randKey].count < 0 {
			assert.True(t, sts)
		}

		if uuidKeyList[randKey].count > 0 {
			assert.False(t, sts)
		}

		time.Sleep(time.Duration(int(time.Millisecond) * rand.Intn(300)))
	}

	fmt.Printf("%+v\n", uuidKeyList)
}
