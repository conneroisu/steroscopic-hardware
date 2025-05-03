package logger

import (
	"context"
	"log/slog"
	"sync"
)

// ChannelHandler implements slog.Handler and sends logs to a channel
type ChannelHandler struct {
	logChan chan LogEntry
	attrs   []slog.Attr
	group   string
	level   slog.Level
	mu      sync.Mutex
}

// NewChannelHandler creates a new ChannelHandler
func NewChannelHandler(logChan chan LogEntry, level slog.Level) *ChannelHandler {
	return &ChannelHandler{
		logChan: logChan,
		level:   level,
	}
}

// Enabled implements slog.Handler.Enabled
func (h *ChannelHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

// Handle implements slog.Handler.Handle
func (h *ChannelHandler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	// Create a copy of the attributes
	h.mu.Lock()
	attrs := make([]slog.Attr, 0)
	copy(attrs, h.attrs)
	h.mu.Unlock()

	// Add record attributes
	r.Attrs(func(attr slog.Attr) bool {
		if h.group != "" {
			attr = slog.Group(h.group, attr)
		}
		attrs = append(attrs, attr)
		return true
	})

	// Send to channel
	h.logChan <- LogEntry{
		Level:   r.Level,
		Time:    r.Time,
		Message: r.Message,
		Attrs:   attrs,
	}

	return nil
}

// WithAttrs implements slog.Handler.WithAttrs
func (h *ChannelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := &ChannelHandler{
		logChan: h.logChan,
		level:   h.level,
		group:   h.group,
	}

	h2.mu.Lock()
	defer h2.mu.Unlock()

	h2.attrs = make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(h2.attrs, h.attrs)
	copy(h2.attrs[len(h.attrs):], attrs)

	return h2
}

// WithGroup implements slog.Handler.WithGroup
func (h *ChannelHandler) WithGroup(name string) slog.Handler {
	h2 := &ChannelHandler{
		logChan: h.logChan,
		level:   h.level,
		attrs:   h.attrs,
	}

	if h.group != "" {
		h2.group = h.group + "." + name
	} else {
		h2.group = name
	}

	return h2
}

// SetLevel sets the minimum level for logging
func (h *ChannelHandler) SetLevel(level slog.Level) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.level = level
}
