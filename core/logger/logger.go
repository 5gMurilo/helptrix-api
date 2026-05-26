package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	once     sync.Once
	instance *slog.Logger
)

func Init() {
	once.Do(func() {
		handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
		instance = slog.New(handler)
		slog.SetDefault(instance)
	})
}

func Get() *slog.Logger {
	if instance == nil {
		Init()
	}
	return instance
}
