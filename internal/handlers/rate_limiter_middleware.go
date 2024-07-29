package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/felipeksw/goexpert-fullcycle-rate-limiter/internal/infra/repository"
	"github.com/felipeksw/goexpert-fullcycle-rate-limiter/internal/usecase"
	"github.com/felipeksw/goexpert-fullcycle-rate-limiter/pgk/config"
)

/*
type rateLimit struct {
	configuration     *config.Config
	repository *repository.RateLimiterStorage
}

func NewRateLimit(cnf *config.Config, rep *repository.RateLimiterStorage) *rateLimit {
	return &rateLimit{
		configuration:     cnf,
		repository: rep,
	}
}
*/

type RateLimit struct {
	Configuration *config.Config
	Repository    *repository.RedisServer
}

type erroDto struct {
	Msg string `json:"msg"`
}

func (rl *RateLimit) RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		v := strings.Split(r.RemoteAddr, ":")
		ip := v[0]
		for i := 1; i < len(v)-1; i++ {
			ip = ip + ":" + v[i]
		}

		apiKey := strings.TrimSpace(r.Header.Get("API_KEY"))

		var key string
		var rps int64
		var blk int64

		if ip != "" {
			key = ip
			rps = rl.Configuration.RequestPerSecondPerIp
			blk = rl.Configuration.BlockingTimeIp
		}

		if apiKey != "" && rl.Configuration.RrequestPerSecondPerApiKey >= rl.Configuration.RequestPerSecondPerIp {
			key = apiKey
			rps = rl.Configuration.RrequestPerSecondPerApiKey
			blk = rl.Configuration.BlockingTimeApiKey
		}

		if key != "" {
			slog.Debug("[RateLimiter]", "key", key, "ReqPerSec", rps, "BlockingTime", blk)
			sts, err := usecase.RateLimitByKey(rl.Repository, key, rps, blk)
			if err != nil {
				slog.Error("[RateLimiter]", "msg", err.Error())
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(&erroDto{Msg: err.Error()})
				return
			}

			if sts {
				slog.Debug("[RateLimiter]", "msg", "requisicao descartada pelo rate militer")
				w.Header().Add("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				json.NewEncoder(w).Encode(&erroDto{Msg: "you have reached the maximum number of requests or actions allowed within a certain time frame"})
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
