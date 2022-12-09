package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/minio-rest/src/minio"
	"github.com/minio-rest/src/v1/controller"
	"github.com/minio-rest/src/v1/health"
)

func AddRoutes(app *fiber.App, service minio.Service) {
	v1 := app.Group("api")
	// Health
	v1.Get("/health", health.HealthCheckHandler)

	v1.Post("/upload", controller.UploadToMinio(service))
	v1.Post("/download", controller.DownloadFromMinio(service))
}
