package adapter

import (
	"context"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/infrastructure/config"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

func WithDzikraRedis() Option {
	return func(a *Adapter) {
		redisHost := config.Envs.RedisDB.Host
		redisPort := config.Envs.RedisDB.Port
		redisPassword := config.Envs.RedisDB.Password
		redisDB := config.Envs.RedisDB.Database

		// Create Redis client
		rdb := redis.NewClient(&redis.Options{
			Addr:     redisHost + ":" + redisPort,
			Password: redisPassword,
			DB:       redisDB,
		})

		// Test the connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err := rdb.Ping(ctx).Result()
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to Dzikra Redis")
		}

		a.DzikraRedis = rdb
		log.Info().Msg("Dzikra Redis connected")
	}
}
