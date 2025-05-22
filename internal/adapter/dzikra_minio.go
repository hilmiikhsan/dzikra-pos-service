package adapter

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/infrastructure/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func WithDzikraMinio() (Option, error) {
	endpoint := config.Envs.MinioStorage.Endpoint
	accessKey := config.Envs.MinioStorage.AccessKey
	secretKey := config.Envs.MinioStorage.SecretKey
	bucketName := config.Envs.MinioStorage.Bucket
	useSSL := config.Envs.MinioStorage.UseSSL

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Error connecting to Dzikra Minio")
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	exists, err := client.BucketExists(context.Background(), bucketName)
	if err != nil {
		log.Fatal().Err(err).Msg("Error checking if bucket exists")
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		if err := client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{}); err != nil {
			log.Fatal().Err(err).Msg("Error creating bucket")
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}

		log.Info().Msgf("Bucket %s successfully created", bucketName)
	}

	log.Info().Msg("Dzikra Minio connected")

	return func(a *Adapter) {
		a.DzikraMinio = client
	}, nil
}
