package main

import (
	"context"
	"dev-tool/log"
	"github.com/pkg/errors"

	"github.com/rs/zerolog"
)

func main() {
	log.InitBergLog("blogs/test.log", zerolog.InfoLevel)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceId", "13185999712")
	log.Info(ctx, "Test")
	log.Debug(ctx, "Test")
	log.Warn(ctx, "Test")
	log.Error(ctx, errors.New("Test"))
}
