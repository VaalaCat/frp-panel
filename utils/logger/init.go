package logger

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	frplog "github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/golib/log"
	"github.com/sirupsen/logrus"
)

func initFrpLogger(frpLogLevel log.Level) {
	frplog.Logger = log.New(
		log.WithCaller(true),
		log.AddCallerSkip(1),
		log.WithLevel(frpLogLevel),
		log.WithOutput(LoggerWriter("frp", logrus.InfoLevel)))
}

func InitLogger() {
	// projectRoot, projectPkg, _ := findProjectRootAndModule()

	Instance().SetReportCaller(true)
	Instance().SetFormatter(NewCustomFormatter(false, true))
	Instance().AddHook(NewStackTraceHook())

	logrus.SetReportCaller(true)
	logrus.SetFormatter(NewCustomFormatter(false, true))
}

func UpdateLoggerOpt(frpLogLevel string, logrusLevel string) {
	ctx := context.Background()

	frpLogLevel = strings.ToLower(frpLogLevel)
	logrusLevel = strings.ToLower(logrusLevel)

	if frpLogLevel == "" {
		frpLogLevel = "info"
	}
	if logrusLevel == "" {
		logrusLevel = "info"
	}

	frpLv, err := log.ParseLevel(frpLogLevel)
	if err != nil {
		Logger(ctx).WithError(err).Errorf("invalid frp log level: %s, use info", frpLogLevel)
		frpLv = log.InfoLevel
	}
	logrusLv, err := logrus.ParseLevel(logrusLevel)
	if err != nil {
		Logger(ctx).WithError(err).Errorf("invalid logrus log level: %s, use info", logrusLevel)
		logrusLv = logrus.InfoLevel
	}

	Instance().SetLevel(logrusLv)
	logrus.SetLevel(logrusLv)

	initFrpLogger(frpLv)
}

func NewCallerPrettyfier(projectRoot, projectPkg string) func(frame *runtime.Frame) (function string, file string) {
	return func(frame *runtime.Frame) (function string, file string) {
		file = frame.File
		if relPath, err := filepath.Rel(projectRoot, frame.File); err == nil {
			file = relPath
		}
		file = " " + file + ":" + strconv.Itoa(frame.Line)

		function = frame.Function
		if strings.HasPrefix(function, projectPkg) {
			function = function[len(projectPkg):]
			function = strings.TrimPrefix(function, "/")
			function = strings.TrimPrefix(function, ".")
		}

		return function, file
	}
}

func FindProjectRootAndModule() (projectRoot string, projectModule string, err error) {
	_, filename, _, ok := runtime.Caller(2)
	if !ok {
		err = fmt.Errorf("cannot get caller info")
		return
	}

	dir := filepath.Dir(filename)
	for {
		if dir == "/" || dir == "." {
			err = fmt.Errorf("go.mod not found")
			return
		}
		modFile := filepath.Join(dir, "go.mod")
		if _, statErr := os.Stat(modFile); statErr == nil {
			projectRoot = dir
			projectModule, err = readModuleName(modFile)
			return
		}
		dir = filepath.Dir(dir)
	}
}

func readModuleName(modFile string) (string, error) {
	f, err := os.Open(modFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module ")), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("module name not found in go.mod")
}
