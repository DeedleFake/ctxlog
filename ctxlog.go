package ctxlog

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type contextKey int

const (
	depthKey contextKey = iota
)

// Depth returns the logging depth of a context.
func Depth(ctx context.Context) int {
	depth, _ := ctx.Value(depthKey).(int)
	return depth
}

type Logger struct {
	once sync.Once

	output io.Writer
	indent string
}

func (log *Logger) init() {
	log.once.Do(func() {
		log.output = os.Stderr
		log.indent = "\t"
	})
}

// SetOutput sets the output that the logger will write to. It
// defaults to os.Stderr.
func (log *Logger) SetOutput(output io.Writer) {
	log.init()

	log.output = output
}

// SetIndent sets the indent that will be prepended to log entries. It
// will be multiplied by the depth of a log call's context. It
// defaults to "\t".
func (log *Logger) SetIndent(indent string) {
	log.init()

	log.indent = indent
}

// Log queues str to be printed when ctx ends. It returns a new
// context with ctx as its parent that is deeper than ctx.
func (log *Logger) Log(ctx context.Context, str string) context.Context {
	log.init()

	str = strings.TrimSpace(str)

	depth := Depth(ctx)
	fmt.Fprintln(output, strings.Repeat(indent, depth)+str)

	return context.WithValue(ctx, depthKey, depth+1)
}

// Logf is like Log, but does Printf-style formatting.
func (log *Logger) Logf(ctx context.Context, format string, args ...interface{}) context.Context {
	log.init()

	format = strings.TrimSpace(format) + "\n"

	depth := Depth(ctx)
	fmt.Fprintf(output, strings.Repeat(indent, depth)+format, args...)

	return context.WithValue(ctx, depthKey, depth+1)
}
