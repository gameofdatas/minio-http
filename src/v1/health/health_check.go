package health

import (
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio-rest/version"
	"github.com/rs/zerolog/log"
)

var upSince = time.Now()

type healthCheck struct {
	Alive     bool   `json:"alive"`
	Since     string `json:"since"`
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	Commit    string `json:"commit"`
}

// HealthCheckHandler : A very simple health check.
// @Summary      HealthCheck
// @Description  Returns the status of the running application
// @Tags         HealthCheck
// @Produce      json
// @Success      200  {object}  healthCheck  "Global health of the application"
// @Failure      500  {object}  object       "Error description formated as {"msg":"string"}"
// @Router       /api/status [get]
func HealthCheckHandler(c *fiber.Ctx) error {
	log.Debug().Msg("handle health check")
	healthCheck := &healthCheck{
		Alive:     true,
		Since:     time.Since(upSince).String(),
		Version:   version.Version,
		BuildDate: version.BuildDate,
		GoVersion: version.GoVersion,
		Commit:    version.GitCommit,
	}
	responseBytes, err := json.Marshal(healthCheck)
	if err != nil {
		log.Error().Msgf("Error marshalling healthCheck response: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": true,
			"msg":   err,
		})
	}
	return c.Send(responseBytes)
}
