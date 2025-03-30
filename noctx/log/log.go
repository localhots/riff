package log

import (
	"context"
	"time"

	"github.com/localhots/riff"
)

var logger *riff.Logger

func Setup(cfg riff.Config) {
	logger = riff.New(cfg)
}

func Trace(msg string, fields ...riff.Field) {
	logger.Trace(context.Background(), msg, fields...)
}

func Debug(msg string, fields ...riff.Field) {
	logger.Debug(context.Background(), msg, fields...)
}

func Info(msg string, fields ...riff.Field) {
	logger.Info(context.Background(), msg, fields...)
}

func Warn(msg string, fields ...riff.Field) {
	logger.Warn(context.Background(), msg, fields...)
}

func Error(msg string, fields ...riff.Field) {
	logger.Error(context.Background(), msg, fields...)
}

func Panic(msg string, fields ...riff.Field) {
	logger.Panic(context.Background(), msg, fields...)
}

func Fatal(msg string, fields ...riff.Field) {
	logger.Fatal(context.Background(), msg, fields...)
}

//
// Types
//

func Cause(err error) riff.Field {
	return riff.Cause(err)
}

func Str(key, value string) riff.Field {
	return riff.Str(key, value)
}

func Int(key string, value int) riff.Field {
	return riff.Int(key, value)
}

func Int64(key string, value int64) riff.Field {
	return riff.Int64(key, value)
}

func Int32(key string, value int32) riff.Field {
	return riff.Int64(key, int64(value))
}

func Int16(key string, value int16) riff.Field {
	return riff.Int64(key, int64(value))
}

func Int8(key string, value int8) riff.Field {
	return riff.Int64(key, int64(value))
}

func Uint(key string, value uint) riff.Field {
	return riff.Uint(key, value)
}

func Uint64(key string, value uint64) riff.Field {
	return riff.Uint64(key, value)
}

func Uint32(key string, value uint32) riff.Field {
	return riff.Uint64(key, uint64(value))
}

func Uint16(key string, value uint16) riff.Field {
	return riff.Uint64(key, uint64(value))
}

func Uint8(key string, value uint8) riff.Field {
	return riff.Uint64(key, uint64(value))
}

func Bool(key string, value bool) riff.Field {
	return riff.Bool(key, value)
}

func Float64(key string, value float64) riff.Field {
	return riff.Float64(key, value)
}

func Float32(key string, value float32) riff.Field {
	return riff.Float32(key, value)
}

func Duration(key string, value time.Duration) riff.Field {
	return riff.Duration(key, value)
}

func Time(key string, value time.Time) riff.Field {
	return riff.Time(key, value)
}

func Any(key string, value any) riff.Field {
	return riff.Any(key, value)
}
