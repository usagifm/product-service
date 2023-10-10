package logger

import (
	"context"

	"github.com/sirupsen/logrus"
)

type ctxKey string

const loggerCtxKey ctxKey = "logger"

func Init(_ context.Context) {
	logrus.SetReportCaller(true)
	logFormat := &logrus.TextFormatter{
		ForceColors:     false,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
	logrus.SetFormatter(logFormat)
	logrus.SetLevel(logrus.DebugLevel)
}

func GetLogger(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(loggerCtxKey)
	if logger == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}
	return logger.(*logrus.Entry)
}

func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, logger)
}
