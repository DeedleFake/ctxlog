package ctxlog

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type contextKey int

const (
	depthKey contextKey = iota
)

func depthFromContext(ctx context.Context) (depth int, ok bool) {
	depth, ok = ctx.Value(depthKey).(int)
	return
}

func Logf(ctx context.Context, format string, args ...interface{}) context.Context {
	depth, _ := depthFromContext(ctx)

	format = strings.TrimSpace(format)
	fmt.Fprintf(os.Stderr, strings.Repeat("\t", depth)+format+"\n", args...)

	return context.WithValue(ctx, depthKey, depth+1)
}
