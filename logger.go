package riff

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type Logger struct {
	cfg       Config
	timeCache func(time.Time) string
	lock      sync.Mutex
}

type Config struct {
	Level           Level
	Output          io.Writer
	Time            bool
	TimeFormat      string
	TimePrecision   time.Duration
	Color           bool
	MinMessageWidth int
	SortFields      bool
	StackTraceLevel Level
	StackTraceSkip  int
}

type Level int

const (
	LevelTrace Level = iota
	LevelDebug       = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

const (
	colorRed      = "\033[31m"
	colorGreen    = "\033[32m"
	colorYellow   = "\033[33m"
	colorBlue     = "\033[34m"
	colorPurple   = "\033[35m"
	colorCyan     = "\033[36m"
	colorOffWhite = "\033[37m"
	colorRedBg    = "\033[48;5;88m"
	colorWhite    = "\033[38;5;255m"
	colorReset    = "\033[0m"
)

var (
	defaultMessageWidth = 40 // characters
	defaultTimeFormat   = "2006-01-02 15:04:05.000"

	DurationPrecision = time.Millisecond
	TimeFormat        = time.RFC3339
)

func New(cfg Config) *Logger {
	l := &Logger{cfg: cfg}
	if l.cfg.TimePrecision > 0 {
		l.timeCache = timeCache(l.cfg.TimeFormat, l.cfg.TimePrecision)
	}
	return l
}

func DefaultConfig() Config {
	return Config{
		Level:           LevelInfo,
		Output:          os.Stderr,
		Time:            true,
		TimeFormat:      defaultTimeFormat,
		TimePrecision:   0, // Disable time cache
		Color:           true,
		MinMessageWidth: defaultMessageWidth,
		SortFields:      true,
		StackTraceLevel: LevelError,
		StackTraceSkip:  4,
	}
}

func (l *Logger) Trace(ctx context.Context, msg string, fields ...Field) {
	if l.cfg.Level == LevelTrace {
		l.print(ctx, LevelTrace, msg, fields)
	}
}

func (l *Logger) Debug(ctx context.Context, msg string, fields ...Field) {
	if l.cfg.Level <= LevelDebug {
		l.print(ctx, LevelDebug, msg, fields)
	}
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...Field) {
	if l.cfg.Level <= LevelInfo {
		l.print(ctx, LevelInfo, msg, fields)
	}
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...Field) {
	if l.cfg.Level <= LevelWarn {
		l.print(ctx, LevelWarn, msg, fields)
	}
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...Field) {
	if l.cfg.Level <= LevelError {
		l.print(ctx, LevelError, msg, fields)
	}
}

func (l *Logger) Panic(ctx context.Context, msg string, fields ...Field) {
	if l.cfg.Level <= LevelPanic {
		l.print(ctx, LevelPanic, msg, fields)
	}
}

func (l *Logger) Fatal(ctx context.Context, msg string, fields ...Field) {
	l.print(ctx, LevelFatal, msg, fields)
	os.Exit(1)
}

//
// Printing
//

func (l *Logger) print(ctx context.Context, lev Level, msg string, fields []Field) {
	if lev < l.cfg.Level {
		return
	}

	buf := getBuffer()
	defer putBuffer(buf)

	l.printTime(buf)
	l.printLevel(buf, lev)
	l.printMessage(buf, msg, len(fields) > 0)
	l.printFields(ctx, buf, lev, fields)
	*buf = append(*buf, '\n')
	l.printStackTrace(buf, lev)

	l.lock.Lock()
	l.cfg.Output.Write(*buf)
	l.lock.Unlock()
}

func (l *Logger) printTime(buf *[]byte) {
	if !l.cfg.Time {
		return
	}

	t := time.Now()
	if l.timeCache != nil {
		*buf = append(*buf, l.timeCache(t)...)
	} else {
		*buf = t.AppendFormat(*buf, l.cfg.TimeFormat)
	}
	*buf = append(*buf, ' ')
}

func (l *Logger) printLevel(buf *[]byte, lev Level) {
	l.writeColorized(buf, lev, l.levelName(lev))
	*buf = append(*buf, ' ')
}

func (l *Logger) printMessage(buf *[]byte, msg string, needsPad bool) {
	*buf = append(*buf, msg...)
	if l.cfg.MinMessageWidth > 0 {
		// Pad the message to the configured width +2 spaces to separate it from
		// the fields.
		for range l.cfg.MinMessageWidth + 2 - len(msg) {
			*buf = append(*buf, ' ')
		}
		// If the message is long enough not to be padded, add an extra space to
		// separate it from the fields
		if len(msg) > l.cfg.MinMessageWidth {
			// Separate message from fields with 2 spaces
			*buf = append(*buf, ' ', ' ')
		}
	} else if needsPad {
		// Separate message from fields with 2 spaces
		*buf = append(*buf, ' ', ' ')
	}
}

func (l *Logger) printFields(ctx context.Context, buf *[]byte, lev Level, fields []Field) {
	if l.cfg.SortFields {
		l.printFieldsSorted(ctx, buf, lev, fields)
	} else {
		l.printFieldsUnsorted(ctx, buf, lev, fields)
	}
}

func (l *Logger) printFieldsUnsorted(ctx context.Context, buf *[]byte, lev Level, fields []Field) {
	for i, f := range fields {
		l.printField(buf, lev, f, i > 0)
	}
	for i, f := range FromContext(ctx) {
		l.printField(buf, lev, f, i+len(fields) > 0)
	}
}

func (l *Logger) printFieldsSorted(ctx context.Context, buf *[]byte, lev Level, fields []Field) {
	// Alias field groups for brevity
	a := FromContext(ctx)
	b := fields

	// Pre-sort both slices
	sortFields(a)
	sortFields(b)

	// Iterate over both slices and print them in sorted order
	var i, j int
	for i < len(a) && j < len(b) {
		if a[i].Key < b[j].Key {
			l.printField(buf, lev, a[i], i+j > 0)
			i++
		} else {
			l.printField(buf, lev, b[j], i+j > 0)
			j++
		}
	}

	// Print remaining fields
	for i < len(a) {
		l.printField(buf, lev, a[i], i+j > 0)
		i++
	}
	for j < len(b) {
		l.printField(buf, lev, b[j], i+j > 0)
		j++
	}
}

func (l *Logger) printField(buf *[]byte, lev Level, f Field, pad bool) {
	if pad {
		*buf = append(*buf, ' ')
	}
	l.writeColorized(buf, lev, f.Key)
	*buf = append(*buf, '=')
	*buf = f.ValueFn(*buf)
}

func (l *Logger) printStackTrace(buf *[]byte, lev Level) {
	if lev >= l.cfg.StackTraceLevel {
		// Print stack trace but skip the first 4 frames which are part of the
		// logger itself.
		*buf = append(*buf, stackTrace(l.cfg.StackTraceSkip)...)
		*buf = append(*buf, '\n')
	}
}

func sortFields(f []Field) {
	if len(f) > 1 {
		insertionSort(f)
	}
}

// insertionSort is great for small slices. Using this custom function instead
// of sort.Slice() reduces the number of allocations to zero.
func insertionSort(f []Field) {
	for i := 1; i < len(f); i++ {
		for j := i; j > 0 && f[j].Key < f[j-1].Key; j-- {
			f[j], f[j-1] = f[j-1], f[j]
		}
	}
}

//
// Helpers
//

func (l *Logger) writeColorized(buf *[]byte, lev Level, str string) {
	if !l.cfg.Color {
		*buf = append(*buf, str...)
		return
	}

	switch lev {
	case LevelTrace, LevelDebug:
		*buf = append(*buf, colorOffWhite...)
	case LevelInfo:
		*buf = append(*buf, colorCyan...)
	case LevelWarn:
		*buf = append(*buf, colorYellow...)
	case LevelError:
		*buf = append(*buf, colorRed...)
	case LevelPanic, LevelFatal:
		*buf = append(*buf, colorRedBg...)
		*buf = append(*buf, colorWhite...)
	}
	*buf = append(*buf, str...)
	*buf = append(*buf, colorReset...)
}

// levelName returns level label that is consistently 4 characters long.
func (l *Logger) levelName(lev Level) string {
	switch lev {
	case LevelTrace:
		return "TRAC"
	case LevelDebug:
		return "DEBU"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERRO"
	case LevelPanic:
		return "PANI"
	case LevelFatal:
		return "FATA"
	default:
		panic("unreachable")
	}
}

func stackTrace(skip int) string {
	// Get up to 100 stack frames
	pc := make([]uintptr, 100)
	// +2 frames to skip for runtime.Callers and stackTrace itself
	n := runtime.Callers(skip+2, pc)
	frames := runtime.CallersFrames(pc[:n])

	var buf bytes.Buffer
	for {
		f, more := frames.Next()
		buf.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", f.Function, f.File, f.Line))
		if !more {
			break
		}
	}
	return buf.String()
}

func timeCache(format string, precision time.Duration) func(time.Time) string {
	var lastTime time.Time
	var lastTimeStr string

	return func(t time.Time) string {
		if !lastTime.IsZero() && t.Sub(lastTime) < precision {
			return lastTimeStr
		}

		lastTime = t
		lastTimeStr = t.Format(format)
		return lastTimeStr
	}
}
