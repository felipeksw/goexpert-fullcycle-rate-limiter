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

		randKey := rand.Intn(9)

		// Reinicia a posição do slice para reutilização
		if uuidKeyList[randKey].ttl < time.Now().UnixMilli() {
			uuidKeyList[randKey].key = ""
			uuidKeyList[randKey].count = 0
			uuidKeyList[randKey].timestsamp = 0
			uuidKeyList[randKey].ttl = 0
		}

		// Adiciona ao slice a chave, caso a posição esteja vazia
		if uuidKeyList[randKey].key == "" {
			uuidKeyList[randKey].key = uuid.New().String()
			uuidKeyList[randKey].count = 1
			uuidKeyList[randKey].timestsamp = time.Now().UnixMilli()
			uuidKeyList[randKey].ttl = 0
		}

		sts, err := usecase.RateLimitByKey(redisRepo, uuidKeyList[randKey].key, requestPerSecond, blockingTime)
		assert.Nil(t, err)

		if time.Now().UnixMilli()-uuidKeyList[randKey].timestsamp <= 1000 {
			if uuidKeyList[randKey].count+1 >= requestPerSecond {
				uuidKeyList[randKey].count = -1
				uuidKeyList[randKey].ttl = time.Now().UnixMilli() + (blockingTime * int64(time.Duration(time.Millisecond)))
				slog.Debug("", "millisecond", int64(time.Duration(time.Millisecond)))
			}
			uuidKeyList[randKey].count = uuidKeyList[randKey].count + 1
		} else {
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

		time.Sleep(time.Duration(int(time.Millisecond) * rand.Intn(100)))
	}

	fmt.Printf("%+v\n", uuidKeyList)
}
