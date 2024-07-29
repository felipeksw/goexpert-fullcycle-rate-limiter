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

type rateLimit struct {
	configuration *config.Config
	repository    any
}

func NewRateLimit(configuration *config.Config, repository any) *rateLimit {
	return &rateLimit{
		configuration: configuration,
		repository:    repository,
	}
}

type erroDto struct {
	Msg string `json:"msg"`
}

func (rl *rateLimit) RateLimiter(next http.Handler) http.Handler {
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
			rps = rl.configuration.RequestPerSecondPerIp
			blk = rl.configuration.BlockingTimeIp
		}

		if apiKey != "" && rl.configuration.RrequestPerSecondPerApiKey >= rl.configuration.RequestPerSecondPerIp {
			key = apiKey
			rps = rl.configuration.RrequestPerSecondPerApiKey
			blk = rl.configuration.BlockingTimeApiKey
		}

		if key != "" {
			slog.Debug("[RateLimiter]", "key", key, "ReqPerSec", rps, "BlockingTime", blk)
			sts, err := usecase.RateLimitByKey(rl.repository.(repository.RateLimiterStorage), key, rps, blk)
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
