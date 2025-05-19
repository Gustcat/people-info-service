package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
)

func SetupLogger(levelLog slog.Level) *slog.Logger {
	var log *slog.Logger

	lumberjackLogger := &lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    10, // МБ
		MaxBackups: 1,  // Кол-во старых файлов
		MaxAge:     2,  // Дней хранить
		Compress:   true,
	}

	log = slog.New(
		slog.NewJSONHandler(lumberjackLogger, &slog.HandlerOptions{Level: levelLog}),
	)

	return log
}
