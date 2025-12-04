package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// ANSI color codes for log levels.
const (
	colorReset  = "\033[0m"
	colorBlue   = "\033[34m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
)

func logWithLevel(level string, msg string, args ...any) {
	// Build the JSON payload.
	payload := map[string]any{
		"time":  time.Now().Format(time.RFC3339Nano),
		"level": level,
		"msg":   msg,
	}

	// Interpret args as key/value pairs: "key1", value1, "key2", value2, ...
	for i := 0; i+1 < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		payload[key] = args[i+1]
	}

	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		// Fallback to a simple print if JSON marshalling fails.
		fmt.Fprintf(os.Stderr, "logger marshal error: %v\n", err)
		return
	}

	// Human prefix: time - colored level - JSON payload
	nowStr := time.Now().Format(time.RFC3339Nano)
	levelColor := levelColor(level)

	fmt.Fprintf(os.Stdout, "%s - %s%s%s - %s\n",
		nowStr,
		levelColor, level, colorReset,
		string(jsonBytes),
	)
}

func levelColor(level string) string {
	switch level {
	case "DEBUG":
		return colorBlue
	case "INFO":
		return colorGreen
	case "WARN":
		return colorYellow
	case "ERROR":
		return colorRed
	default:
		return colorReset
	}
}

// Debug logs a debug-level message.
func Debug(msg string, args ...any) {
	logWithLevel("DEBUG", msg, args...)
}

// Info logs an info-level message.
func Info(msg string, args ...any) {
	logWithLevel("INFO", msg, args...)
}

// Warn logs a warning-level message.
func Warn(msg string, args ...any) {
	logWithLevel("WARN", msg, args...)
}

// Error logs an error-level message.
func Error(msg string, args ...any) {
	logWithLevel("ERROR", msg, args...)
}
