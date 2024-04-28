package database

import (
	"context"
	"log"
	"os"
	"time"
)

var started int

const (
	reset  = "\033[0m"
	red    = "\033[31m"
	green  = "\033[32m"
	yellow = "\033[33m"
	blue   = "\033[34m"
)

type logger interface {
	Printf(string, ...interface{})
}

type hookLogger struct {
	log logger
}

func QueryLogger() *hookLogger {
	return &hookLogger{
		log: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (h *hookLogger) Before(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, &started, time.Now()), nil
}

func (h *hookLogger) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	duration := time.Since(ctx.Value(&started).(time.Time))
	h.log.Printf("%s%s Args: `%q`. Took: %s%s", blue, query, args, duration, reset)
	return ctx, nil
}

func (h *hookLogger) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	duration := time.Since(ctx.Value(&started).(time.Time))
	h.log.Printf("%sError: %v, Query: `%s`, Args: `%q`, Took: %s%s",
		red, err, query, args, duration, reset)
	return err
}
