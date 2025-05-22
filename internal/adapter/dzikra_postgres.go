package adapter

import (
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/infrastructure/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func WithDzikraPostgres() Option {
	return func(a *Adapter) {
		dbUser := config.Envs.DzikraPostgres.Username
		dbPassword := config.Envs.DzikraPostgres.Password
		dbName := config.Envs.DzikraPostgres.Database
		dbHost := config.Envs.DzikraPostgres.Host
		dbSSLMode := config.Envs.DzikraPostgres.SslMode
		dbPort := config.Envs.DzikraPostgres.Port

		dbMaxPoolSize := config.Envs.DB.MaxOpenCons
		dbMaxIdleConns := config.Envs.DB.MaxIdleCons
		dbConnMaxLifetime := config.Envs.DB.ConnMaxLifetime

		connectionString := "user=" + dbUser + " password=" + dbPassword + " host=" + dbHost + " port=" + dbPort + " dbname=" + dbName + " sslmode=" + dbSSLMode + " TimeZone=UTC"
		db, err := sqlx.Connect("postgres", connectionString)
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to Postgres")
		}

		db.SetMaxOpenConns(dbMaxPoolSize)
		db.SetMaxIdleConns(dbMaxIdleConns)
		db.SetConnMaxLifetime(time.Duration(dbConnMaxLifetime) * time.Second)

		// check connection
		err = db.Ping()
		if err != nil {
			log.Fatal().Err(err).Msg("Error connecting to Dzikra Postgres")
		}

		a.DzikraPostgres = db
		log.Info().Msg("Dzikra Postgres connected")
	}
}
