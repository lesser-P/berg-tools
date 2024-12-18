package example

import (
	"context"
	"dev-tool/log"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"testing"
)

func TestExampleBergLog(t *testing.T) {
	log.InitBergLog("stdout", zerolog.InfoLevel)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceId", "0x000000000000berg")
	log.Info(ctx, "Test")
	log.Debug(ctx, "Test")
	log.Warn(ctx, "Test")
	log.Error(ctx, errors.New("Test"))
	log.InitBergLog("test.log", zerolog.InfoLevel)
	log.Info(ctx, "Test")
	log.Debug(ctx, "Test")
	log.Warn(ctx, "Test")
	log.Error(ctx, errors.New("Test"))
}
