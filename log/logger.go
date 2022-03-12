package log

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"os"
	"time"
)

var (
	Logger zerolog.Logger
)

func SetUpLogger() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	//zerolog.TimeFieldFormat = zerolog.TimeFormatUnix // Readable or epoch?

	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(time.UTC)
	}

	Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
}
