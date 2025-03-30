package riff_test

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/localhots/riff"
	"github.com/localhots/riff/ctx/log"
)

func TestLogger(t *testing.T) {
	cfg := riff.DefaultConfig()
	cfg.Level = riff.LevelTrace
	log.Setup(cfg)
	ctx := context.Background()
	err := errors.New("task already exists")

	log.Trace(ctx, "Parsing message")
	log.Debug(ctx, "Callback received",
		log.Str("device_unique_id", "G4000E-1000-F"),
		log.Int("task_id", 123456),
		log.Str("status", "success"),
		log.Str("template_name", "index.tpl"),
	)
	log.Info(ctx, "Starting task",
		log.Str("device_unique_id", "G4000E-1000-F"),
		log.Int("task_id", 123456),
	)
	log.Info(ctx, "Extremely long message, sorry about that; it is meant to prove that buffer can grow"+strings.Repeat(" @@@", 300),
		log.Str("device_unique_id", "G4000E-1000-F"),
		log.Int("task_id", 123456),
	)
	log.Warn(ctx, "Duplicate task, but also this message exceeds 40 characters",
		log.Int("task_id", 123456),
	)
	log.Warn(ctx, "Duplicate task but exactly 40 characters",
		log.Int("task_id", 123456),
	)
	log.Warn(riff.WithContext(ctx, log.Str("foo", "bar")), "Duplicate task is exactly 39 characters",
		log.Int("task_id", 123456),
	)
	log.Error(ctx, "Failed to process task", log.Cause(err),
		log.Int("task_id", 123456),
	)
	log.Fatal(ctx, "Failed to start service", log.Cause(err),
		log.Str("service", "api"),
	)
}

func TestBare(t *testing.T) {
	log.Setup(riff.Config{
		Level:           riff.LevelDebug,
		Output:          os.Stderr,
		StackTraceLevel: riff.LevelError,
	})
	ctx := context.Background()

	log.Info(ctx, "Starting task",
		log.Str("device_unique_id", "G4000E-1000-F"),
		log.Int("task_id", 123456),
		log.Str("status", "success"),
		log.Str("template_name", "index.tpl"),
	)
}

func TestOptimized(t *testing.T) {
	log.Setup(riff.Config{
		Level:           riff.LevelDebug,
		Output:          os.Stderr,
		Time:            true,
		TimeFormat:      riff.TimeFormat,
		TimePrecision:   1 * time.Millisecond,
		Color:           false,
		MinMessageWidth: 0,
		SortFields:      false,
		StackTraceLevel: riff.LevelError,
	})
	ctx := context.Background()

	log.Info(ctx, "Starting task",
		log.Str("device_unique_id", "G4000E-1000-F"),
		log.Int("task_id", 123456),
		log.Str("status", "success"),
		log.Str("template_name", "index.tpl"),
	)
}

func TestPretty(t *testing.T) {
	log.Setup(riff.Config{
		Level:           riff.LevelDebug,
		Output:          os.Stderr,
		Time:            true,
		TimeFormat:      riff.TimeFormat,
		TimePrecision:   1 * time.Millisecond,
		Color:           true,
		MinMessageWidth: 40,
		SortFields:      false,
		StackTraceLevel: riff.LevelError,
	})
	ctx := context.Background()

	log.Info(ctx, "Starting task",
		log.Str("device_unique_id", "G4000E-1000-F"),
		log.Int("task_id", 123456),
		log.Str("status", "success"),
		log.Str("template_name", "index.tpl"),
	)
}

func TestPrettySorted(t *testing.T) {
	log.Setup(riff.Config{
		Level:           riff.LevelDebug,
		Output:          os.Stderr,
		Time:            true,
		TimeFormat:      riff.TimeFormat,
		TimePrecision:   1 * time.Millisecond,
		Color:           true,
		MinMessageWidth: 40,
		SortFields:      true,
		StackTraceLevel: riff.LevelError,
	})
	ctx := context.Background()

	log.Info(ctx, "Starting task",
		log.Str("device_unique_id", "G4000E-1000-F"),
		log.Int("task_id", 123456),
		log.Str("status", "success"),
		log.Str("template_name", "index.tpl"),
	)
}
func TestPrettySortedContext(t *testing.T) {
	log.Setup(riff.Config{
		Level:           riff.LevelDebug,
		Output:          os.Stderr,
		Time:            true,
		TimeFormat:      riff.TimeFormat,
		TimePrecision:   1 * time.Millisecond,
		Color:           true,
		MinMessageWidth: 40,
		SortFields:      true,
		StackTraceLevel: riff.LevelError,
	})
	ctx := context.Background()
	ctx = log.WithContext(ctx, log.Str("foo", "bar"))
	ctx = log.WithContext(ctx, log.Str("zone", "two"))

	log.Info(ctx, "Starting task",
		log.Str("device_unique_id", "G4000E-1000-F"),
		log.Int("task_id", 123456),
		log.Str("status", "success"),
		log.Str("template_name", "index.tpl"),
	)
}

func BenchmarkBare(b *testing.B) {
	log.Setup(riff.Config{
		Level:           riff.LevelDebug,
		Output:          io.Discard,
		StackTraceLevel: riff.LevelError,
	})
	ctx := context.Background()

	b.ResetTimer()
	for range b.N {
		log.Info(ctx, "Starting task",
			log.Str("device_unique_id", "G4000E-1000-F"),
			log.Int("task_id", 123456),
			log.Str("status", "success"),
			log.Str("template_name", "index.tpl"),
		)
	}
}

func BenchmarkOptimized(b *testing.B) {
	log.Setup(riff.Config{
		Level:           riff.LevelDebug,
		Output:          io.Discard,
		Time:            true,
		TimeFormat:      riff.TimeFormat,
		TimePrecision:   1 * time.Millisecond,
		Color:           false,
		MinMessageWidth: 0,
		SortFields:      false,
		StackTraceLevel: riff.LevelError,
	})
	ctx := context.Background()

	b.ResetTimer()
	for range b.N {
		log.Info(ctx, "Starting task",
			log.Str("device_unique_id", "G4000E-1000-F"),
			log.Int("task_id", 123456),
			log.Str("status", "success"),
			log.Str("template_name", "index.tpl"),
		)
	}
}

func BenchmarkPretty(b *testing.B) {
	log.Setup(riff.Config{
		Level:           riff.LevelDebug,
		Output:          io.Discard,
		Time:            true,
		TimeFormat:      riff.TimeFormat,
		TimePrecision:   1 * time.Millisecond,
		Color:           true,
		MinMessageWidth: 40,
		SortFields:      false,
		StackTraceLevel: riff.LevelError,
	})
	ctx := context.Background()

	b.ResetTimer()
	for range b.N {
		log.Info(ctx, "Starting task",
			log.Any("device_unique_id", "G4000E-1000-F"),
			log.Any("task_id", 123456),
			log.Any("status", "success"),
			log.Any("template_name", "index.tpl"),
		)
	}
}

func BenchmarkPrettySorted(b *testing.B) {
	log.Setup(riff.Config{
		Level:           riff.LevelDebug,
		Output:          io.Discard,
		Time:            true,
		TimeFormat:      riff.TimeFormat,
		TimePrecision:   1 * time.Millisecond,
		Color:           true,
		MinMessageWidth: 40,
		SortFields:      true,
		StackTraceLevel: riff.LevelError,
	})
	ctx := context.Background()

	b.ResetTimer()
	for range b.N {
		log.Info(ctx, "Starting task",
			log.Any("device_unique_id", "G4000E-1000-F"),
			log.Any("task_id", 123456),
			log.Any("status", "success"),
			log.Any("template_name", "index.tpl"),
		)
	}
}

func BenchmarkPrettySortedContext(b *testing.B) {
	log.Setup(riff.Config{
		Level:           riff.LevelDebug,
		Output:          io.Discard,
		Time:            true,
		TimeFormat:      riff.TimeFormat,
		TimePrecision:   1 * time.Millisecond,
		Color:           true,
		MinMessageWidth: 40,
		SortFields:      true,
		StackTraceLevel: riff.LevelError,
	})
	ctx := context.Background()
	ctx = log.WithContext(ctx, log.Any("foo", "bar"))
	ctx = log.WithContext(ctx, log.Any("one", "two"))

	b.ResetTimer()
	for range b.N {
		log.Info(ctx, "Starting task",
			log.Any("device_unique_id", "G4000E-1000-F"),
			log.Any("task_id", 123456),
			log.Any("status", "success"),
			log.Any("template_name", "index.tpl"),
		)
	}
}
