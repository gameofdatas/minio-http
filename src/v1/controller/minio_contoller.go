package controller

import (
	"encoding/base64"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/minio-rest/src/minio"
	"github.com/minio/minio-go/pkg/s3utils"
)

type Request struct {
	Bucket      string `json:"bucket"`                 // name of the bucket where file should be uploaded
	FilePath    string `json:"file_path"`              // path where file should be uploaded
	FileName    string `json:"file_name"`              // name of the file to be uploaded
	FileContent string `json:"file_content,omitempty"` // base64 encoded string of file
	RetentionMs int    `json:"retention_ms,omitempty"` // milliseconds after the object should be deleted from MinIO
}

func DownloadFromMinio(svc minio.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := &Request{}
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
		if err := s3utils.CheckValidBucketName(req.Bucket); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
		if err := s3utils.CheckValidObjectName(req.FileName); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
		fileContent, err := svc.GetObject(req.FilePath, req.FileName, req.Bucket)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
		encoded := base64.StdEncoding.EncodeToString(fileContent)
		return c.SendString(encoded)
	}
}

func UploadToMinio(svc minio.Service) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := &Request{}
		if err := c.BodyParser(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
		decodedBase64, err := base64.StdEncoding.DecodeString(req.FileContent)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
		err = svc.PutObject(decodedBase64, req.FilePath, req.FileName, req.Bucket, req.RetentionMs)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": true,
				"msg":   err.Error(),
			})
		}
		return c.SendStatus(http.StatusCreated)
	}
}
