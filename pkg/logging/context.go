package logging

import (
	"context"
)

type sinkKey struct{}

type loggerKey struct{}

func WithSink(ctx context.Context, sink Sink) context.Context {
	return context.WithValue(ctx, sinkKey{}, sink)
}

func getSink(ctx context.Context) Sink {
	if ctx != nil {
		if v, ok := ctx.Value(sinkKey{}).(Sink); ok {
			return v
		}
	}
	return nil
}

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

func getLogger(ctx context.Context) Logger {
	if ctx != nil {
		if v, ok := ctx.Value(loggerKey{}).(Logger); ok {
			return v
		}
	}
	return nil
}
