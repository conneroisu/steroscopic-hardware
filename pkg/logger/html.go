package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	slogformatter "github.com/samber/slog-formatter"
)

// HTMLRecord represents a log record in a format suitable for HTML display
type HTMLRecord struct {
	Time    time.Time
	Level   slog.Level
	Message string
	Attrs   map[string]any
}

// HTMLHandler is a slog.Handler that serves logs as HTML
type HTMLHandler struct {
	mu          sync.RWMutex
	records     []HTMLRecord
	maxRecords  int
	formatter   slog.Handler
	baseHandler slog.Handler
}

// HandlerOptions contains options for HTMLHandler
type HandlerOptions struct {
	// MaxRecords is the maximum number of records to keep in memory
	MaxRecords int
	// BaseHandler is an optional underlying handler (for multi-output logging)
	BaseHandler slog.Handler
	// Formatters are the slog-formatter formatters to apply
	Formatters []slogformatter.Formatter
}

// NewHTMLHandler creates a new HTMLHandler
func NewHTMLHandler(opts HandlerOptions) *HTMLHandler {
	if opts.MaxRecords <= 0 {
		opts.MaxRecords = 1000 // Default to 1000 records
	}

	var baseHandler slog.Handler
	if opts.BaseHandler != nil {
		baseHandler = opts.BaseHandler
	} else {
		baseHandler = slog.NewJSONHandler(nil, nil)
	}

	h := &HTMLHandler{
		records:     make([]HTMLRecord, 0, opts.MaxRecords),
		maxRecords:  opts.MaxRecords,
		baseHandler: baseHandler,
	}

	// Apply formatters if provided
	if len(opts.Formatters) > 0 {
		h.formatter = slogformatter.NewFormatterHandler(opts.Formatters...)(baseHandler)
	} else {
		h.formatter = baseHandler
	}

	return h
}

// Enabled implements slog.Handler
func (h *HTMLHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.formatter.Enabled(ctx, level)
}

// Handle implements slog.Handler
func (h *HTMLHandler) Handle(ctx context.Context, record slog.Record) error {
	// Convert the record to our HTMLRecord format
	htmlRecord := HTMLRecord{
		Time:    record.Time,
		Level:   record.Level,
		Message: record.Message,
		Attrs:   make(map[string]interface{}),
	}

	// Extract attributes
	record.Attrs(func(attr slog.Attr) bool {
		// Store the attribute in our map
		h.extractAttr(htmlRecord.Attrs, "", attr)
		return true
	})

	// Add the record to our in-memory store
	h.mu.Lock()
	if len(h.records) >= h.maxRecords {
		// Remove the oldest record
		h.records = h.records[1:]
	}
	h.records = append(h.records, htmlRecord)
	h.mu.Unlock()

	// Forward to the base handler if it exists
	return h.formatter.Handle(ctx, record)
}

// WithAttrs implements slog.Handler
func (h *HTMLHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Create a new handler with the attributes
	newHandler := &HTMLHandler{
		records:     h.records,
		maxRecords:  h.maxRecords,
		formatter:   h.formatter.WithAttrs(attrs),
		baseHandler: h.baseHandler.WithAttrs(attrs),
	}
	return newHandler
}

// WithGroup implements slog.Handler
func (h *HTMLHandler) WithGroup(name string) slog.Handler {
	// Create a new handler with the group
	newHandler := &HTMLHandler{
		records:     h.records,
		maxRecords:  h.maxRecords,
		formatter:   h.formatter.WithGroup(name),
		baseHandler: h.baseHandler.WithGroup(name),
	}
	return newHandler
}

// ServeHTTP implements http.Handler for viewing logs
func (h *HTMLHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Copy the records to avoid holding the lock while rendering
	h.mu.RLock()
	records := make([]HTMLRecord, len(h.records))
	copy(records, h.records)
	h.mu.RUnlock()

	// Check request parameters
	query := r.URL.Query()

	// Check if JSON is requested
	if query.Get("format") == "json" {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(records); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Check if this is a partial HTMX update request
	isHtmxRequest := r.Header.Get("HX-Request") == "true"

	// Get last seen index if provided
	var lastIndex int
	if lastIndexStr := query.Get("lastIndex"); lastIndexStr != "" {
		if idx, err := strconv.Atoi(lastIndexStr); err == nil {
			lastIndex = idx
		}
	}

	// Filter records based on lastIndex
	var newRecords []HTMLRecord
	if lastIndex > 0 && lastIndex < len(records) {
		newRecords = records[lastIndex:]
	} else {
		newRecords = records
	}

	// Prepare template data
	data := struct {
		Records         []HTMLRecord
		LastIndex       int
		IsPartialUpdate bool
	}{
		Records:         newRecords,
		LastIndex:       len(records),
		IsPartialUpdate: isHtmxRequest,
	}

	// Set content type
	w.Header().Set("Content-Type", "text/html")

	// Select the appropriate template based on request type
	var templateName string
	if isHtmxRequest {
		templateName = "logs-partial"
		// Set headers for SSE with HTMX
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("HX-Trigger", "newLogsLoaded")
	} else {
		templateName = "logs-full"
	}

	// Parse and execute the template
	tmpl, err := template.New(templateName).Parse(htmlTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RegisterHTTPHandler registers the HTML handler at the given path
func (h *HTMLHandler) RegisterHTTPHandler(path string, mux *http.ServeMux) {
	if mux == nil {
		mux = http.DefaultServeMux
	}
	mux.Handle(path, h)
}

// Helper function to extract attributes including nested ones
func (h *HTMLHandler) extractAttr(m map[string]interface{}, prefix string, attr slog.Attr) {
	if attr.Value.Kind() == slog.KindGroup {
		// Handle group attribute
		group := attr.Value.Group()
		groupPrefix := prefix
		if attr.Key != "" {
			if prefix != "" {
				groupPrefix = prefix + "." + attr.Key
			} else {
				groupPrefix = attr.Key
			}
		}

		for _, groupAttr := range group {
			h.extractAttr(m, groupPrefix, groupAttr)
		}
	} else {
		// Handle non-group attribute
		key := attr.Key
		if prefix != "" {
			key = prefix + "." + key
		}

		// Store the value based on its kind
		switch attr.Value.Kind() {
		case slog.KindString:
			m[key] = attr.Value.String()
		case slog.KindInt64:
			m[key] = attr.Value.Int64()
		case slog.KindUint64:
			m[key] = attr.Value.Uint64()
		case slog.KindFloat64:
			m[key] = attr.Value.Float64()
		case slog.KindBool:
			m[key] = attr.Value.Bool()
		case slog.KindDuration:
			m[key] = attr.Value.Duration().String()
		case slog.KindTime:
			m[key] = attr.Value.Time().Format(time.RFC3339)
		case slog.KindAny:
			m[key] = fmt.Sprintf("%v", attr.Value.Any())
		default:
			m[key] = fmt.Sprintf("%v", attr.Value)
		}
	}
}

// HTML template for rendering logs
const htmlTemplate = `
{{define "logs-full"}}
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Application Logs</title>
    <!-- HTMX for live updates -->
    <script src="https://unpkg.com/htmx.org@1.9.9"></script>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        h1 {
            color: #2c3e50;
            border-bottom: 1px solid #eee;
            padding-bottom: 10px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .log-status {
            font-size: 14px;
            font-weight: normal;
            background-color: #eee;
            padding: 4px 8px;
            border-radius: 4px;
        }
        .log-status.connected {
            background-color: #d4edda;
            color: #155724;
        }
        .log-status.disconnected {
            background-color: #f8d7da;
            color: #721c24;
        }
        .log-container {
            margin-top: 20px;
        }
        .log-entry {
            border: 1px solid #ddd;
            border-radius: 4px;
            margin-bottom: 10px;
            padding: 12px;
            animation: fadeIn 0.3s ease-in;
        }
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(-10px); }
            to { opacity: 1; transform: translateY(0); }
        }
        .log-entry.DEBUG { background-color: #f8f9fa; border-left: 5px solid #6c757d; }
        .log-entry.INFO { background-color: #f1f8ff; border-left: 5px solid #0366d6; }
        .log-entry.WARN { background-color: #fff8f1; border-left: 5px solid #f66a0a; }
        .log-entry.ERROR { background-color: #ffdce0; border-left: 5px solid #d73a49; }
        .log-time {
            color: #666;
            font-size: 0.9em;
        }
        .log-level {
            display: inline-block;
            padding: 2px 6px;
            border-radius: 3px;
            font-size: 0.8em;
            font-weight: bold;
            margin-right: 8px;
        }
        .log-level.DEBUG { background-color: #e9ecef; color: #495057; }
        .log-level.INFO { background-color: #cce5ff; color: #004085; }
        .log-level.WARN { background-color: #fff3cd; color: #856404; }
        .log-level.ERROR { background-color: #f8d7da; color: #721c24; }
        .log-message {
            font-weight: 500;
            margin: 8px 0;
        }
        .log-attrs {
            background-color: rgba(0,0,0,0.03);
            border-radius: 3px;
            padding: 8px 12px;
            margin-top: 8px;
            overflow-x: auto;
        }
        .log-attr {
            margin: 4px 0;
            display: flex;
        }
        .log-attr-key {
            font-weight: 500;
            color: #0366d6;
            margin-right: 8px;
            flex: 0 0 200px;
        }
        .log-attr-value {
            font-family: SFMono-Regular, Consolas, "Liberation Mono", Menlo, monospace;
            word-break: break-word;
        }
        .controls {
            margin-bottom: 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .btn {
            background-color: #0366d6;
            color: white;
            border: none;
            padding: 8px 16px;
            border-radius: 4px;
            cursor: pointer;
            margin-left: 8px;
        }
        .btn:hover {
            background-color: #0258c3;
        }
        .btn.pause {
            background-color: #f66a0a;
        }
        .btn.pause:hover {
            background-color: #e36209;
        }
        .btn.clear {
            background-color: #d73a49;
        }
        .btn.clear:hover {
            background-color: #cb2431;
        }
        .search-container {
            display: flex;
            gap: 8px;
            flex-grow: 1;
            max-width: 500px;
        }
        #search-input {
            padding: 8px;
            border: 1px solid #ddd;
            border-radius: 4px;
            width: 100%;
        }
        .empty-state {
            text-align: center;
            padding: 40px;
            color: #666;
        }
        .auto-scroll {
            display: flex;
            align-items: center;
            margin-left: 16px;
        }
        .auto-scroll input {
            margin-right: 8px;
        }
    </style>
</head>
<body>
    <h1>
        Application Logs
        <span id="connection-status" class="log-status connected">Connected</span>
    </h1>
    
    <div class="controls">
        <div class="search-container">
            <input type="text" id="search-input" placeholder="Search logs...">
        </div>
        <div class="auto-scroll">
            <input type="checkbox" id="auto-scroll" checked>
            <label for="auto-scroll">Auto-scroll</label>
        </div>
        <button id="pause-btn" class="btn pause" onclick="togglePolling()">Pause</button>
        <button class="btn clear" onclick="clearLogs()">Clear</button>
    </div>
    
    <div id="log-container" class="log-container"
         hx-get="?lastIndex={{.LastIndex}}"
         hx-trigger="newLogsLoaded from:body delay:1s, every 3s"
         hx-swap="beforeend">
        {{if len .Records eq 0}}
            <div class="empty-state">
                <p>No logs available</p>
            </div>
        {{else}}
            {{template "log-entries" .}}
        {{end}}
    </div>
</body>
</html>
{{end}}

{{define "logs-partial"}}
    {{template "log-entries" .}}
{{end}}

{{define "log-entries"}}
    {{range $index, $record := .Records}}
        <div class="log-entry {{$record.Level}}" data-log-content="{{$record.Message}}">
            <div>
                <span class="log-time">{{$record.Time.Format "2006-01-02 15:04:05.000"}}</span>
                <span class="log-level {{$record.Level}}">{{$record.Level}}</span>
            </div>
            <div class="log-message">{{$record.Message}}</div>
            {{if $record.Attrs}}
                <div class="log-attrs">
                    {{range $key, $value := $record.Attrs}}
                        <div class="log-attr">
                            <div class="log-attr-key">{{$key}}</div>
                            <div class="log-attr-value">{{$value}}</div>
                        </div>
                    {{end}}
                </div>
            {{end}}
        </div>
    {{end}}
{{end}}
`
