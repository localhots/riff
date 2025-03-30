package log

import (
	"context"
	"time"

	"github.com/localhots/riff"
)

var logger *riff.Logger

// Setup initializes the logger.
func Setup(cfg riff.Config) {
	logger = riff.New(cfg)
}

// Trace logs a message at the Trace level, which is the most verbose level.
func Trace(ctx context.Context, msg string, fields ...riff.Field) {
	logger.Trace(ctx, msg, fields...)
}

// Debug logs a message at the Debug level, which is less verbose than Trace,
// but is still excessively detailed.
func Debug(ctx context.Context, msg string, fields ...riff.Field) {
	logger.Debug(ctx, msg, fields...)
}

// Info logs a message at the Info level, which is great for general kind of
// records.
func Info(ctx context.Context, msg string, fields ...riff.Field) {
	logger.Info(ctx, msg, fields...)
}

// Warn logs a message at the Warn level, which indicates a potential problem
// that should be looked at.
func Warn(ctx context.Context, msg string, fields ...riff.Field) {
	logger.Warn(ctx, msg, fields...)
}

// Error logs a message at the Error level.
func Error(ctx context.Context, msg string, fields ...riff.Field) {
	logger.Error(ctx, msg, fields...)
}

// Panic logs a message at the Panic level, which indicates a very serious
// problem. It doesn't actually panic, but it should be treated as an emergency.
func Panic(ctx context.Context, msg string, fields ...riff.Field) {
	logger.Panic(ctx, msg, fields...)
}

// Fatal logs a message at the Fatal level, which indicates an unrecoverable
// error. It will terminate the program after logging the message.
func Fatal(ctx context.Context, msg string, fields ...riff.Field) {
	logger.Fatal(ctx, msg, fields...)
}

// WithContext adds logging fields to the context.
func WithContext(ctx context.Context, fields ...riff.Field) context.Context {
	return riff.WithContext(ctx, fields...)
}

// FromContext returns the logging fields from the context.
func FromContext(ctx context.Context) []riff.Field {
	return riff.FromContext(ctx)
}

//
// Types
//

// Cause returns a Field that wraps the given error in a standardized way.
func Cause(err error) riff.Field {
	return riff.Cause(err)
}

// Str returns a Field with the given key and string value.
func Str(key, value string) riff.Field {
	return riff.Str(key, value)
}

// Int returns a Field with the given key and int value.
func Int(key string, value int) riff.Field {
	return riff.Int(key, value)
}

// Int64 returns a Field with the given key and int64 value.
func Int64(key string, value int64) riff.Field {
	return riff.Int64(key, value)
}

// Int32 returns a Field with the given key and int32 value converted to int64.
func Int32(key string, value int32) riff.Field {
	return riff.Int64(key, int64(value))
}

// Int16 returns a Field with the given key and int16 value converted to int64.
func Int16(key string, value int16) riff.Field {
	return riff.Int64(key, int64(value))
}

// Int8 returns a Field with the given key and int8 value converted to int64.
func Int8(key string, value int8) riff.Field {
	return riff.Int64(key, int64(value))
}

// Uint returns a Field with the given key and uint value.
func Uint(key string, value uint) riff.Field {
	return riff.Uint(key, value)
}

// Uint64 returns a Field with the given key and uint64 value.
func Uint64(key string, value uint64) riff.Field {
	return riff.Uint64(key, value)
}

// Uint32 returns a Field with the given key and uint32 value converted to
// uint64.
func Uint32(key string, value uint32) riff.Field {
	return riff.Uint64(key, uint64(value))
}

// Uint16 returns a Field with the given key and uint16 value converted to
// uint64.
func Uint16(key string, value uint16) riff.Field {
	return riff.Uint64(key, uint64(value))
}

// Uint8 returns a Field with the given key and uint8 value converted to
// uint64.
func Uint8(key string, value uint8) riff.Field {
	return riff.Uint64(key, uint64(value))
}

// Float64 returns a Field with the given key and float64 value.
func Float64(key string, value float64) riff.Field {
	return riff.Float64(key, value)
}

// Float32 returns a Field with the given key and float32 value.
func Float32(key string, value float32) riff.Field {
	return riff.Float32(key, value)
}

// Bool returns a Field with the given key and boolean value.
func Bool(key string, value bool) riff.Field {
	return riff.Bool(key, value)
}

// Duration returns a Field with the given key and time.Duration value.
// The value is truncated to the configured precision (default is 1ms).
func Duration(key string, value time.Duration) riff.Field {
	return riff.Duration(key, value)
}

// Time returns a Field with the given key and time.Time value.
// The value is formatted using the configured time format.
func Time(key string, value time.Time) riff.Field {
	return riff.Time(key, value)
}

// Any returns a Field with the given key and any value.
func Any(key string, value any) riff.Field {
	return riff.Any(key, value)
}
