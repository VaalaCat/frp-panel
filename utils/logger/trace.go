package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

type StackTraceHook struct{}

func NewStackTraceHook() *StackTraceHook {
	return &StackTraceHook{}
}

func (hook *StackTraceHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.FatalLevel,
	}
}

func (hook *StackTraceHook) Fire(entry *logrus.Entry) error {
	stack := getConciseStackTrace()
	if stack != "" {
		entry.Data["stack"] = stack
	}
	return nil
}

func getConciseStackTrace() string {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(4, pcs)
	if n == 0 {
		return ""
	}

	var conciseStack []string
	frames := runtime.CallersFrames(pcs[:n])

	for {
		frame, more := frames.Next()

		if strings.HasPrefix(frame.Function, "runtime.") ||
			strings.HasPrefix(frame.Function, "reflect") ||
			strings.HasPrefix(frame.Function, "github.com/sirupsen/logrus") ||
			strings.HasPrefix(frame.Function, "go.uber.org/fx") ||
			strings.HasPrefix(frame.Function, "go.uber.org/dig") {
			if more {
				continue
			} else {
				break
			}
		}

		fileName := filepath.Base(frame.File)
		frameStr := fmt.Sprintf("%s(%s:%d)", simpleFuncName(frame.Function), fileName, frame.Line)
		conciseStack = append(conciseStack, frameStr)

		if !more {
			break
		}
	}

	return strings.Join(conciseStack, " <- ")
}

func simpleFuncName(fullName string) string {
	lastSlash := strings.LastIndex(fullName, "/")
	if lastSlash >= 0 {
		fullName = fullName[lastSlash+1:]
	}
	lastDot := strings.LastIndex(fullName, ".")
	if lastDot >= 0 {
		if _, err := strconv.Atoi(fullName[lastDot+1:]); err == nil {
			return fullName
		}
		return fullName[lastDot+1:]
	}
	return fullName
}
