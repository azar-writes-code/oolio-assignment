package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Init configures and sets the global slog default logger.
// Optionally accepts extra io.Writers for fan-out (e.g. Loki) alongside stdout.
func Init(logLevel, format string, extra ...io.Writer) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(strings.ToUpper(logLevel))); err != nil {
		level = slog.LevelInfo
	}

	writers := append([]io.Writer{os.Stdout}, extra...)
	out := io.MultiWriter(writers...)

	var handler slog.Handler
	if strings.ToLower(format) == "json" {
		opts := &slog.HandlerOptions{Level: level, AddSource: true}
		handler = slog.NewJSONHandler(out, opts)
	} else {
		handler = &CustomTextHandler{
			w:     out,
			level: level,
		}
	}

	slog.SetDefault(slog.New(handler))
}


// CustomTextHandler implements slog.Handler with pattern: "LOG_LEVEL: (date :: file) : message <json_attrs>"
type CustomTextHandler struct {
	w     io.Writer
	level slog.Level
	attrs []slog.Attr
}

func (h *CustomTextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *CustomTextHandler) Handle(ctx context.Context, r slog.Record) error {
	levelStr := r.Level.String()
	var coloredLevel string
	switch r.Level {
	case slog.LevelDebug:
		coloredLevel = "\033[36m" + levelStr + "\033[0m" // Cyan
	case slog.LevelInfo:
		coloredLevel = "\033[32m" + levelStr + "\033[0m" // Green
	case slog.LevelWarn:
		coloredLevel = "\033[33m" + levelStr + "\033[0m" // Yellow
	case slog.LevelError:
		coloredLevel = "\033[31m" + levelStr + "\033[0m" // Red
	default:
		coloredLevel = levelStr
	}

	timeStr := r.Time.Format("02-01-2006 - 15:04:05") // DD-MM-YYYY - time

	file := "unknown"
	// Fetch source file explicitly when AddSource isn't automatic via standard handler framework
	// skip [Callers, slog.Record.PC...] -> depth varies, but typically we can just capture pc
	// r.PC is captured automatically when we use slog.Info etc! 
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			file = fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
		}
	}

	// Prepare attributes seamlessly as JSON if present
	attrMap := make(map[string]any)

	for _, a := range h.attrs {
		attrMap[a.Key] = a.Value.Any()
	}

	r.Attrs(func(a slog.Attr) bool {
		attrMap[a.Key] = a.Value.Any()
		return true
	})

	var attrStr string
	if len(attrMap) > 0 {
		b, err := json.Marshal(attrMap)
		if err == nil {
			attrStr = " " + string(b)
		}
	}

	// User requested layout: LOG_LEVEL: [date :: file] : message
	msg := fmt.Sprintf("%s: \033[1;36m[%s :: %s]\033[0m : %s%s\n", coloredLevel, timeStr, file, r.Message, attrStr)
	_, err := h.w.Write([]byte(msg))
	return err
}

func (h *CustomTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CustomTextHandler{
		w:     h.w,
		level: h.level,
		attrs: append(h.attrs[:len(h.attrs):len(h.attrs)], attrs...),
	}
}

func (h *CustomTextHandler) WithGroup(name string) slog.Handler {
	// TODO: implement group-prefixing for structured log keys (e.g. "group.key")
	// Currently a no-op: group semantics are intentionally not supported in the custom text format.
	return h
}
