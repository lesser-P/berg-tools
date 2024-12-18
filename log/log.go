package log

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"path/filepath"
	"time"
)

type BergLog struct {
	log zerolog.Logger
}

var blog *BergLog

func InitBergLog(logPath string, level zerolog.Level) {
	var logger zerolog.Logger
	if logPath == "stdout" {
		cw := zerolog.ConsoleWriter{
			Out:          os.Stdout,
			TimeFormat:   time.DateTime,
			TimeLocation: time.Local,
			PartsOrder:   []string{"time", "level", "message"},
		}
		logger = zerolog.New(cw).With().Timestamp().Logger()
	} else {
		dir := filepath.Dir(logPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				panic("InitBerg fail:" + err.Error())
			}
		}
		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic("init fail:" + err.Error())
		}
		logger = zerolog.New(file).With().Timestamp().Logger()
	}
	logger.Level(level)
	blog = &BergLog{
		logger,
	}
}
func Info(ctx context.Context, msg string) {
	var traceId string
	traceId = ctx.Value("traceId").(string)
	blog.log.Info().Str("traceId", traceId).Msg(msg)
}
func Warn(ctx context.Context, msg string) {
	var traceId string
	traceId = ctx.Value("traceId").(string)
	blog.log.Warn().Str("traceId", traceId).Msg(msg)
}
func Error(ctx context.Context, err error) {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	var traceId string
	traceId = ctx.Value("traceId").(string)
	blog.log.Error().Stack().Err(err).Str("traceId", traceId).Msg(err.Error())
}
func Debug(ctx context.Context, msg string) {
	var traceId string
	traceId = ctx.Value("traceId").(string)
	blog.log.Debug().Str("traceId", traceId).Msg(msg)
}
