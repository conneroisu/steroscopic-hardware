package logger

import (
	"context"
	"log/slog"
	"sync"
)

// MultiHandler combines multiple slog.Handler instances
type MultiHandler struct {
	handlers []slog.Handler
	mu       sync.Mutex
}

// NewMultiHandler creates a new MultiHandler
func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

// Enabled implements slog.Handler.Enabled
func (h *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle implements slog.Handler.Handle
func (h *MultiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, handler := range h.handlers {
		if err := handler.Handle(ctx, r); err != nil {
			return err
		}
	}
	return nil
}

// WithAttrs implements slog.Handler.WithAttrs
func (h *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return NewMultiHandler(handlers...)
}

// WithGroup implements slog.Handler.WithGroup
func (h *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return NewMultiHandler(handlers...)
}

// SetLevel sets the minimum level for all handlers that support it
func (h *MultiHandler) SetLevel(level slog.Level) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, handler := range h.handlers {
		if ch, ok := handler.(*ChannelHandler); ok {
			ch.SetLevel(level)
		} else if mh, ok := handler.(*MultiHandler); ok {
			mh.SetLevel(level)
		} else if lh, ok := handler.(interface{ SetLevel(slog.Level) }); ok {
			// Try to set level on any handler that has a SetLevel method
			lh.SetLevel(level)
		}
	}
}
