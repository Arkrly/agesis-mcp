package observability

import (
	"encoding/json"
	"io"
	"log"
	"sync"
	"time"
)

// Logger emits structured JSON logs.
type Logger struct {
	base *log.Logger
	mu   sync.Mutex
}

func NewLogger(w io.Writer) *Logger {
	return &Logger{
		base: log.New(w, "", 0),
	}
}

func (l *Logger) Info(message string, fields map[string]any) {
	l.write("info", message, fields)
}

func (l *Logger) Error(message string, fields map[string]any) {
	l.write("error", message, fields)
}

func (l *Logger) write(level string, message string, fields map[string]any) {
	payload := map[string]any{
		"ts":      time.Now().UTC().Format(time.RFC3339Nano),
		"level":   level,
		"message": message,
	}
	for key, value := range fields {
		payload[key] = value
	}
	body, err := json.Marshal(payload)
	if err != nil {
		l.base.Printf(`{"level":"error","message":"log marshal failed","error":%q}`, err.Error())
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.base.Print(string(body))
}
