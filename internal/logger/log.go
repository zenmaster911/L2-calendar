package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

type Logger struct {
	*slog.Logger
	file *os.File
}

func NewSLog(logPath string, level slog.Level) (*Logger, error) {
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to stat logger %w", err)
	}

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to create log file %w", err)
	}

	multiwriter := io.MultiWriter(os.Stdout, file)
	options := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}
	handler := slog.NewTextHandler(multiwriter, options)

	return &Logger{
		Logger: slog.New(handler),
		file:   file,
	}, nil
}

func (l *Logger) Close() error {
	return l.file.Close()
}

func (l *Logger) InfoCtx(ctx context.Context, msg string, args ...any) {
	l.Logger.Log(ctx, slog.LevelInfo, msg, args...)
}

func (l *Logger) ErrorCtx(ctx context.Context, msg string, err error, args ...any) {
	args = append(args, slog.Any("error", err))
	l.Logger.Log(ctx, slog.LevelError, msg, args...)
}
