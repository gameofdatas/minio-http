package log

import (
	"bytes"
	"os"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type ZLogger struct {
	Logger        zerolog.Logger
	LogContent    *bytes.Buffer
	LogStartTime  time.Time
	LoggingServer *LoggingServer
}

func NewLogger() (*ZLogger, error) {
	zLogger := new(ZLogger)
	zLogger.LogContent = new(bytes.Buffer)
	zLogger.LogStartTime = time.Now()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	multiWriter := zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout})
	zLogger.LoggingServer = NewLoggerClient()
	zLogger.Logger = zerolog.New(multiWriter).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	zlog.Logger = zLogger.Logger
	return zLogger, nil
}
