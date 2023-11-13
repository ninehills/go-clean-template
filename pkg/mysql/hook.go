package mysql

import (
	"context"
	"strings"
	"time"

	"github.com/ninehills/go-webapp-template/pkg/logger"
)

type ctxKeySQLStarted struct{}

type Hook struct {
	// FIXME: this logger cannot auto-reload.
	l logger.Logger
}

func NewHook(l logger.Logger) *Hook {
	return &Hook{
		l: l,
	}
}

func (h *Hook) Before(ctx context.Context, _ string, _ ...interface{}) (context.Context, error) {
	return context.WithValue(ctx, ctxKeySQLStarted{}, time.Now()), nil
}

func (h *Hook) After(ctx context.Context, query string, args ...interface{}) (context.Context, error) {
	h.l.Debugf("SQL Query: `%s`, Args: `%q`. took: %s",
		strings.ReplaceAll(query, "\n", " "), args, getTimeUsedFromCtx(ctx),
	)

	return ctx, nil
}

func (h *Hook) OnError(ctx context.Context, err error, query string, args ...interface{}) error {
	h.l.Debugf("SQL Error: %v, Query: `%s`, Args: `%q`, Took: %s",
		err, strings.ReplaceAll(query, "\n", " "), args, getTimeUsedFromCtx(ctx),
	)

	return err
}

func getTimeUsedFromCtx(ctx context.Context) string {
	started, ok := ctx.Value(ctxKeySQLStarted{}).(time.Time)
	if !ok {
		return "-"
	}

	return time.Since(started).String()
}
