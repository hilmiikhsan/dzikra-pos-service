package config

import (
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/config"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
)

var (
	Envs *Config // Envs is global vars Config.
	once sync.Once
)

type Config struct {
	App struct {
		Name                    string `env:"APP_NAME" env-default:"dzikra-user-service"`
		Environtment            string `env:"APP_ENV" env-default:"development"`
		BaseURL                 string `env:"APP_BASE_URL" env-default:"http://localhost:9090"`
		Port                    string `env:"APP_PORT" env-default:"9090"`
		GrpcPort                string `env:"APP_GRPC_PORT" env-default:"7000"`
		LogLevel                string `env:"APP_LOG_LEVEL" env-default:"debug"`
		LogFile                 string `env:"APP_LOG_FILE" env-default:"./logs/app.log"`
		LogFileWs               string `env:"APP_LOG_FILE_WS" env-default:"./logs/ws.log"`
		Domain                  string `env:"APP_DOMAIN" env-default:"localhost:9090"`
		LocalStoragePublicPath  string `env:"LOCAL_STORAGE_PUBLIC_PATH" env-default:"./storage/public"`
		LocalStoragePrivatePath string `env:"LOCAL_STORAGE_PRIVATE_PATH" env-default:"./storage/private"`
	}
	DB struct {
		ConnectionTimeout int `env:"DB_CONN_TIMEOUT" env-default:"30" env-description:"database timeout in seconds"`
		MaxOpenCons       int `env:"DB_MAX_OPEN_CONS" env-default:"20" env-description:"database max open conn in seconds"`
		MaxIdleCons       int `env:"DB_MAX_IdLE_CONS" env-default:"20" env-description:"database max idle conn in seconds"`
		ConnMaxLifetime   int `env:"DB_CONN_MAX_LIFETIME" env-default:"0" env-description:"database conn max lifetime in seconds"`
	}
	Guard struct {
		JwtPrivateKey             string `env:"JWT_PRIVATE_KEY" env-default:""`
		JwtTokenExpiration        string `env:"JWT_TOKEN_EXPIRATION" env-default:"15m"`
		JwtRefreshTokenExpiration string `env:"JWT_REFRESH_TOKEN_EXPIRATION" env-default:"720h"`
	}
	DzikraPostgres struct {
		Host     string `env:"DZIKRA_POSTGRES_HOST" env-default:"localhost"`
		Port     string `env:"DZIKRA_POSTGRES_PORT" env-default:"5432"`
		Username string `env:"DZIKRA_POSTGRES_USER" env-default:"postgres"`
		Password string `env:"DZIKRA_POSTGRES_PASSWORD" env-default:"postgres"`
		Database string `env:"DZIKRA_POSTGRES_DB" env-default:"dzikra_user"`
		SslMode  string `env:"DZIKRA_POSTGRES_SSL_MODE" env-default:"disable"`
	}
	RedisDB struct {
		Host     string `env:"DZIKRA_REDIS_HOST" env-default:"redis"`
		Port     string `env:"DZIKRA_REDIS_PORT" env-default:"6379"`
		Password string `env:"DZIKRA_REDIS_PASSWORD" env-default:"password"`
		Database int    `env:"DZIKRA_REDIS_DB" env-default:"0"`
	}
	MinioStorage struct {
		Endpoint  string `env:"DZIKRA_MINIO_ENDPOINT" env-default:"localhost:9000"`
		AccessKey string `env:"DZIKRA_MINIO_ACCESS_KEY" env-default:""`
		SecretKey string `env:"DZIKRA_MINIO_SECRET_KEY" env-default:""`
		Bucket    string `env:"DZIKRA_MINIO_BUCKET" env-default:""`
		UseSSL    bool   `env:"DZIKRA_MINIO_USE_SSL" env-default:"false"`
		PublicURL string `env:"DZIKRA_MINIO_PUBLIC_URL" env-default:"http://localhost:9000"`
	}
	Auth struct {
		AuthGrpcHost string `env:"AUTH_GRPC_HOST" env-default:"localhost:7000"`
	}
	Notification struct {
		NotificationGrpcHost string `env:"NOTIFICATION_GRPC_HOST" env-default:"localhost:7001"`
	}
	Order struct {
		OrderGrpcHost string `env:"ORDER_GRPC_HOST" env-default:"localhost:7002"`
	}
}

// Option is Configure type return func.
type Option = func(c *Configure) error

// Configure is the data struct.
type Configure struct {
	path     string
	filename string
}

// Configuration create instance.
func Configuration(opts ...Option) *Configure {
	c := &Configure{}

	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			log.Fatal().Err(err).Msg("error while configuring")
			panic(err)
		}
	}
	return c
}

// Initialize will create instance of Configure.
func (c *Configure) Initialize() {
	once.Do(func() {
		Envs = &Config{}
		if err := config.Load(config.Opts{
			Config:    Envs,
			Paths:     []string{c.path},
			Filenames: []string{c.filename},
		}); err != nil {
			log.Fatal().Err(err).Msg("get config error")
		}

		Envs.App.Name = utils.GetEnv("APP_NAME", Envs.App.Name)
		Envs.App.Port = utils.GetEnv("APP_PORT", Envs.App.Port)
		Envs.App.GrpcPort = utils.GetEnv("APP_GRPC_PORT", Envs.App.GrpcPort)
		Envs.App.LogLevel = utils.GetEnv("APP_LOG_LEVEL", Envs.App.LogLevel)
		Envs.App.LogFile = utils.GetEnv("APP_LOG_FILE", Envs.App.LogFile)
		Envs.App.LogFileWs = utils.GetEnv("APP_LOG_FILE_WS", Envs.App.LogFileWs)
		Envs.App.Domain = utils.GetEnv("APP_DOMAIN", Envs.App.Domain)
		Envs.App.LocalStoragePublicPath = utils.GetEnv("LOCAL_STORAGE_PUBLIC_PATH", Envs.App.LocalStoragePublicPath)
		Envs.App.LocalStoragePrivatePath = utils.GetEnv("LOCAL_STORAGE_PRIVATE_PATH", Envs.App.LocalStoragePrivatePath)
		Envs.DB.ConnectionTimeout = utils.GetIntEnv("DB_CONN_TIMEOUT", Envs.DB.ConnectionTimeout)
		Envs.DB.MaxOpenCons = utils.GetIntEnv("DB_MAX_OPEN_CONS", Envs.DB.MaxOpenCons)
		Envs.DB.MaxIdleCons = utils.GetIntEnv("DB_MAX_IdLE_CONS", Envs.DB.MaxIdleCons)
		Envs.DB.ConnMaxLifetime = utils.GetIntEnv("DB_CONN_MAX_LIFETIME", Envs.DB.ConnMaxLifetime)
		Envs.Guard.JwtPrivateKey = utils.GetEnv("JWT_PRIVATE_KEY", Envs.Guard.JwtPrivateKey)
		Envs.Guard.JwtTokenExpiration = utils.GetEnv("JWT_TOKEN_EXPIRATION", Envs.Guard.JwtTokenExpiration)
		Envs.Guard.JwtRefreshTokenExpiration = utils.GetEnv("JWT_REFRESH_TOKEN_EXPIRATION", Envs.Guard.JwtRefreshTokenExpiration)
		Envs.DzikraPostgres.Host = utils.GetEnv("DZIKRA_POSTGRES_HOST", Envs.DzikraPostgres.Host)
		Envs.DzikraPostgres.Port = utils.GetEnv("DZIKRA_POSTGRES_PORT", Envs.DzikraPostgres.Port)
		Envs.DzikraPostgres.Username = utils.GetEnv("DZIKRA_POSTGRES_USER", Envs.DzikraPostgres.Username)
		Envs.DzikraPostgres.Password = utils.GetEnv("DZIKRA_POSTGRES_PASSWORD", Envs.DzikraPostgres.Password)
		Envs.DzikraPostgres.Database = utils.GetEnv("DZIKRA_POSTGRES_DB", Envs.DzikraPostgres.Database)
		Envs.DzikraPostgres.SslMode = utils.GetEnv("DZIKRA_POSTGRES_SSL_MODE", Envs.DzikraPostgres.SslMode)
		Envs.RedisDB.Host = utils.GetEnv("DZIKRA_REDIS_HOST", Envs.RedisDB.Host)
		Envs.RedisDB.Port = utils.GetEnv("DZIKRA_REDIS_PORT", Envs.RedisDB.Port)
		Envs.RedisDB.Password = utils.GetEnv("DZIKRA_REDIS_PASSWORD", Envs.RedisDB.Password)
		Envs.RedisDB.Database = utils.GetIntEnv("DZIKRA_REDIS_DB", Envs.RedisDB.Database)
		Envs.MinioStorage.Endpoint = utils.GetEnv("DZIKRA_MINIO_ENDPOINT", Envs.MinioStorage.Endpoint)
		Envs.MinioStorage.AccessKey = utils.GetEnv("DZIKRA_MINIO_ACCESS_KEY", Envs.MinioStorage.AccessKey)
		Envs.MinioStorage.SecretKey = utils.GetEnv("DZIKRA_MINIO_SECRET_KEY", Envs.MinioStorage.SecretKey)
		Envs.MinioStorage.Bucket = utils.GetEnv("DZIKRA_MINIO_BUCKET", Envs.MinioStorage.Bucket)
		Envs.MinioStorage.UseSSL = utils.GetBoolEnv("DZIKRA_MINIO_USE_SSL", Envs.MinioStorage.UseSSL)
		Envs.MinioStorage.PublicURL = utils.GetEnv("DZIKRA_MINIO_PUBLIC_URL", Envs.MinioStorage.PublicURL)
		Envs.Auth.AuthGrpcHost = utils.GetEnv("AUTH_GRPC_HOST", Envs.Auth.AuthGrpcHost)
		Envs.Notification.NotificationGrpcHost = utils.GetEnv("NOTIFICATION_GRPC_HOST", Envs.Notification.NotificationGrpcHost)
		Envs.Order.OrderGrpcHost = utils.GetEnv("ORDER_GRPC_HOST", Envs.Order.OrderGrpcHost)
	})
}

// WithPath will assign to field path Configure.
func WithPath(path string) Option {
	return func(c *Configure) error {
		c.path = path
		return nil
	}
}

// WithFilename will assign to field name Configure.
func WithFilename(name string) Option {
	return func(c *Configure) error {
		c.filename = name
		return nil
	}
}
