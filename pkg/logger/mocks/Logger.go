package mocks

import (
	"testing"
	"time"
)

const (
	LevelError = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	// colorGray   = "\033[90m"
)

const TestTimeStampFormat = "15:04:05.000"

type Logger struct {
	// a flexible int to indicate level of logger. Default to log all.
	// the larger this number is, the more verbose logs will be printed
	// -1, no log
	// 0, error
	// 1. warn
	// 2, info
	// 3. debug
	level int

	t *testing.T
}

func NewLogger(t *testing.T, level ...int) *Logger {
	l := &Logger{t: t, level: LevelDebug}
	if len(level) > 0 {
		l.level = level[0]
	}
	return l
}

// Debugf provides a mock function with given fields: format, a
func (_m *Logger) Debugf(format string, a ...interface{}) {
	if _m.level >= LevelDebug {
		b := []interface{}{colorYellow, time.Now().Format(TestTimeStampFormat), colorReset}
		a = append(b, a...)
		_m.t.Logf("%s [%s] DEBUG:%s "+format, a...)
	}
}

// Infof provides a mock function with given fields: format, a
func (_m *Logger) Infof(format string, a ...interface{}) {
	if _m.level >= LevelInfo {
		b := []interface{}{colorBlue, time.Now().Format(TestTimeStampFormat), colorReset}
		a = append(b, a...)
		_m.t.Logf("%s [%s] INFO:%s "+format, a...)
	}
}

// Warnf provides a mock function with given fields: format, a
func (_m *Logger) Warnf(format string, a ...interface{}) {
	if _m.level >= LevelWarn {
		b := []interface{}{colorGreen, time.Now().Format(TestTimeStampFormat), colorReset}
		a = append(b, a...)
		_m.t.Logf("%s [%s] WARN:%s "+format, a...)
	}
}

// Errorf provides a mock function with given fields: format, a
func (_m *Logger) Errorf(format string, a ...interface{}) {
	if _m.level >= LevelError {
		b := []interface{}{colorRed, time.Now().Format(TestTimeStampFormat), colorReset}
		a = append(b, a...)
		_m.t.Logf("%s [%s] ERROR:%s "+format, a...)
	}
}
