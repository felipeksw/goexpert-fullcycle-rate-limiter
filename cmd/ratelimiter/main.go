package main

import (
	"log/slog"
	"net/http"

	"github.com/felipeksw/goexpert-fullcycle-rate-limiter/internal/handlers"
	"github.com/felipeksw/goexpert-fullcycle-rate-limiter/internal/infra/repository"
	"github.com/felipeksw/goexpert-fullcycle-rate-limiter/pgk/config"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func main() {

	slog.SetLogLoggerLevel(slog.LevelDebug)

	configs, err := config.LoadConfig([]string{"."})
	if err != nil {
		panic(err)
	}
	slog.Debug("VIPER", "configs", configs)

	redisRepo := repository.NewRedisCahe(
		redis.NewClient(&redis.Options{
			Addr: configs.RedisHost + ":" + configs.RedisPort,
		}),
	)

	rl := &handlers.RateLimit{
		Configuration: configs,
		Repository:    redisRepo,
	}

	mux := mux.NewRouter()
	mux.Use(rl.RateLimiter)
	mux.HandleFunc("/", handlers.HelloWord).Methods("GET")

	slog.Info("[webserver listening]", "port", configs.AppPort)
	err = http.ListenAndServe(":"+configs.AppPort, mux)
	slog.Error("could not start the webserver:" + err.Error())
}
