package handler

import (
	"context"
	"log/slog"
	"runtime"
)

type (
	HandlerOptions struct {
		StripFilePath func(string) string
		AddSource     bool
	}

	SourceHandler struct {
		wrapped slog.Handler
		opts    HandlerOptions
	}
)

func NewSourceHandler(baseHandler slog.Handler, opts *HandlerOptions) slog.Handler {
	if opts == nil {
		opts = &HandlerOptions{}
	}
	return &SourceHandler{
		wrapped: baseHandler,
		opts:    *opts,
	}
}

func (h *SourceHandler) Source(r slog.Record) *slog.Source {
	fc := runtime.FuncForPC(r.PC)
	if fc == nil {
		return nil
	}
	file, line := fc.FileLine(r.PC)
	return &slog.Source{
		Function: fc.Name(),
		File:     h.opts.StripFilePath(file),
		Line:     line,
	}
}

func (h *SourceHandler) Handle(ctx context.Context, record slog.Record) error {
	if !h.opts.AddSource && record.Level < slog.LevelError {
		return h.wrapped.Handle(ctx, record)
	}

	source := h.Source(record)
	record.AddAttrs(slog.Any(slog.SourceKey, source))
	return h.wrapped.Handle(ctx, record)
}

func (h *SourceHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.wrapped.Enabled(ctx, level)
}

func (h *SourceHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewSourceHandler(h.wrapped.WithAttrs(attrs), &h.opts)
}

func (h *SourceHandler) WithGroup(name string) slog.Handler {
	return NewSourceHandler(h.wrapped.WithGroup(name), &h.opts)
}
