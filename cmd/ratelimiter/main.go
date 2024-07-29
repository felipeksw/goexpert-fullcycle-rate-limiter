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

	/*
		// Requisito: Crie uma “strategy” que permita trocar facilmente o Redis por outro mecanismo de persistência
		// Solução: Uso de interfaces para o repositório
		// Caso seja deseje trocar o mecanismo de persistênica, basta implementar seguinto o contrato da
		//  interface repository.RateLimiterStorage e posteriormente trocar o parâmetro repository
		//  na criação do rate limit
		// Para fins didáticos, como exemplo, foi construído um repositório utilizando a biblioteca
		//  fs-cache (funções não implementadas)
		//
		ops := fscache.New()
		fscacheClient := repository.NewFsCahe(ops.Memdis())
		rateLimit := handlers.NewRateLimit(configs, fscacheClient)
	*/

	redisClient := repository.NewRedisCahe(
		redis.NewClient(&redis.Options{
			Addr: configs.RedisHost + ":" + configs.RedisPort,
		}),
	)
	rateLimit := handlers.NewRateLimit(configs, redisClient)

	mux := mux.NewRouter()
	mux.Use(rateLimit.RateLimiter)
	mux.HandleFunc("/", handlers.HelloWord).Methods("GET")

	slog.Info("[webserver listening]", "port", configs.AppPort)
	err = http.ListenAndServe(":"+configs.AppPort, mux)
	slog.Error("could not start the webserver:" + err.Error())
}
