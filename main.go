package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minio-rest/config"
	"github.com/minio-rest/log"
	"github.com/minio-rest/src/minio"
	routes "github.com/minio-rest/src/v1"
	zlog "github.com/rs/zerolog/log"
)

func main() {
	conf := config.Config()
	zlogger, err := log.NewLogger()
	if err != nil {
		zlog.Fatal().Err(err)
	}
	objectStorer, err := minio.NewClient(conf)
	if err != nil {
		zlog.Fatal().Err(err)
	}
	app := fiber.New(fiber.Config{
		ErrorHandler: routes.StandardErrorHandler,
	})
	app.Use(log.Middleware(zlogger.Logger, nil))
	routes.AddRoutes(app, minio.Service{
		MinioClient: &minio.ObjectStoreClient{
			Client: objectStorer,
		},
	})

	err = app.Listen(conf.GetString("API_PORT"))
	if err != nil {
		zlog.Fatal().Err(err)
	}
}
