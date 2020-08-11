package rpacontext

import (
	"context"
	"os"
	"os/signal"

	"github.com/rs/zerolog"
)

func NewGracefulContext() (ctx context.Context) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())

	go func(cancel context.CancelFunc) {
		s := make(chan os.Signal, 1)
		signal.Notify(s, os.Interrupt)
		<-s
		println()

		cancel()
	}(cancel)

	return ctx
}

func WithLogger(ctx context.Context, log zerolog.Logger) (newctx context.Context) {
	return log.WithContext(ctx)
}

func Logger(ctx context.Context) (logger *zerolog.Logger) {
	return zerolog.Ctx(ctx)
}
