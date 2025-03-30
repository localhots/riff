package riff

import (
	"fmt"
	"strconv"
	"time"
)

// Field represents a key-value pair of a log entry.
type Field struct {
	Key     string
	ValueFn ValueFn
}

// ValueFn writes to the given byte slice and returns the result.
type ValueFn func([]byte) []byte

// Cause returns a field that wraps the given error in a standardized way.
func Cause(err error) Field {
	return Str("error", err.Error())
}

// Str returns a field with the given key and a string value.
func Str(key, value string) Field {
	return field(key, func(b []byte) []byte {
		return append(b, value...)
	})
}

// Int returns a field with the given key and an int value.
func Int(key string, value int) Field {
	return Int64(key, int64(value))
}

// Int64 returns a field with the given key and an int64 value.
func Int64(key string, value int64) Field {
	return field(key, func(b []byte) []byte {
		return strconv.AppendInt(b, value, 10)
	})
}

// Uint returns a field with the given key and a uint value.
func Uint(key string, value uint) Field {
	return Uint64(key, uint64(value))
}

// Uint64 returns a field with the given key and a uint64 value.
func Uint64(key string, value uint64) Field {
	return field(key, func(b []byte) []byte {
		return strconv.AppendUint(b, value, 10)
	})
}

// Bool returns a field with the given key and a boolean value.
func Bool(key string, value bool) Field {
	return field(key, func(b []byte) []byte {
		return strconv.AppendBool(b, value)
	})
}

// Float64 returns a field with the given key and a float64 value.
func Float64(key string, value float64) Field {
	return field(key, func(b []byte) []byte {
		return strconv.AppendFloat(b, value, 'f', -1, 64)
	})
}

// Float32 returns a field with the given key and a float32 value.
func Float32(key string, value float32) Field {
	return field(key, func(b []byte) []byte {
		return strconv.AppendFloat(b, float64(value), 'f', -1, 32)
	})
}

// Duration returns a field with the given key and a time.Duration value.
// The duration is truncated to the configured precision (default is
// milliseconds).
func Duration(key string, value time.Duration) Field {
	return field(key, func(b []byte) []byte {
		return append(b, value.Truncate(DurationPrecision).String()...)
	})
}

// Time returns a field with the given key and a time.Time value. Time is
// formatted using the configured time format (default is RFC3339).
func Time(key string, value time.Time) Field {
	return field(key, func(b []byte) []byte {
		return value.AppendFormat(b, TimeFormat)
	})
}

// Any returns a field with the given key and an any value. For most built-in
// types the value is written to the buffer using most effcient method, for
// other types the value is converted to a string using fmt.Sprint.
func Any(key string, value any) Field {
	switch v := value.(type) {
	case string:
		return Str(key, v)
	case []byte:
		return Str(key, string(v))
	case int:
		return Int(key, v)
	case int8:
		return Int64(key, int64(v))
	case int16:
		return Int64(key, int64(v))
	case int32:
		return Int64(key, int64(v))
	case int64:
		return Int64(key, v)
	case uint:
		return Uint(key, v)
	case uint8:
		return Uint64(key, uint64(v))
	case uint16:
		return Uint64(key, uint64(v))
	case uint32:
		return Uint64(key, uint64(v))
	case uint64:
		return Uint64(key, v)
	case float32:
		return Float32(key, v)
	case float64:
		return Float64(key, v)
	case bool:
		return Bool(key, v)
	case time.Duration:
		return Duration(key, v)
	case time.Time:
		return Time(key, v)
	default:
		// TODO: Add support for custom encoders
		return field(key, func(b []byte) []byte {
			return append(b, fmt.Sprint(value)...)
		})
	}
}

func field(key string, fn ValueFn) Field {
	return Field{
		Key:     key,
		ValueFn: fn,
	}
}
