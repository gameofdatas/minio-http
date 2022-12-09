package log

import (
	"encoding/json"
	"errors"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio-rest/src/v1/controller"
	"github.com/rs/zerolog"
)

type ErrMessage struct {
	Err     bool   `json:"error"`
	Message string `json:"msg"`
}

type logFields struct {
	srvName     string
	ID          string
	RemoteIP    string
	Host        string
	Method      string
	Path        string
	Protocol    string
	StatusCode  int
	Latency     float64
	Error       error
	Stack       []byte
	RequestBody []byte
}

func (lf *logFields) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("srv", "minio-rest").
		Str("remote_ip", lf.RemoteIP).
		Str("host", lf.Host).
		Str("method", lf.Method).
		Str("path", lf.Path).
		Str("protocol", lf.Protocol).
		Int("status_code", lf.StatusCode).
		Float64("latency", lf.Latency).
		Str("tag", "request").
		Bytes("request_body", lf.RequestBody)

	if lf.Error != nil {
		e.Err(lf.Error)
	}

	if lf.Stack != nil {
		e.Bytes("stack", lf.Stack)
	}
}

// Middleware requestid + logger + recover for request traceability
func Middleware(log zerolog.Logger, filter func(*fiber.Ctx) bool) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if filter != nil && filter(c) {
			return c.Next()
		}

		start := time.Now()

		requestBody := c.Request().Body()
		req := new(controller.Request)
		json.Unmarshal(requestBody, req)
		rid := req.FilePath

		fields := &logFields{
			ID:          rid,
			RemoteIP:    c.IP(),
			Method:      c.Method(),
			RequestBody: c.Request().Body(),
			Host:        c.Hostname(),
			Path:        c.Path(),
			Protocol:    c.Protocol(),
		}

		defer func() {

			fields.StatusCode = c.Response().StatusCode()
			fields.Latency = time.Since(start).Seconds()
			body := c.Response().Body()
			errMsg := new(ErrMessage)
			json.Unmarshal(body, errMsg)
			if errMsg.Err {
				fields.Error = errors.New(errMsg.Message)
				fields.Stack = debug.Stack()
			}
			switch {
			case fields.StatusCode >= 500:
				log.Error().EmbedObject(fields).Msg("server_error")
			case fields.StatusCode >= 400:
				log.Error().EmbedObject(fields).Msg("client_error")
			case fields.StatusCode >= 300:
				log.Warn().EmbedObject(fields).Msg("redirect")
			case fields.StatusCode >= 200:
				log.Info().EmbedObject(fields).Msg("success")
			case fields.StatusCode >= 100:
				log.Info().EmbedObject(fields).Msg("informative")
			default:
				log.Warn().EmbedObject(fields).Msg("unknown_status")
			}
		}()

		return c.Next()
	}
}
