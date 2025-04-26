package logger

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	// projectRoot, projectPkg, _ := findProjectRootAndModule()

	Instance().SetReportCaller(true)
	Instance().AddHook(NewStackTraceHook())

	logrus.SetReportCaller(true)
	logrus.SetReportCaller(true)
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
