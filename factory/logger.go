// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/luisguillenc/yalogi"
	"github.com/sirupsen/logrus"

	"github.com/luids-io/common/config"
)

// Logger is a factory for a logger
func Logger(cfg *config.LoggerCfg, debug bool) (yalogi.Logger, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid logger config: %v", err)
	}
	level, _ := logrus.ParseLevel(cfg.Level)
	logger := logrus.New()
	if debug {
		level = logrus.DebugLevel
		logger.SetReportCaller(true)
		logger.SetFormatter(&logrus.TextFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				gopath := os.Getenv("GOPATH")
				if gopath == "" {
					gopath = fmt.Sprintf("%s/go", os.Getenv("HOME"))
				}
				repopath := fmt.Sprintf("%s/src/github.com/luids-io", gopath)
				filename := strings.Replace(f.File, repopath, "~", -1)
				function := f.Function
				fpath := strings.Split(f.Function, "/")
				if len(fpath) > 0 {
					function = fpath[len(fpath)-1]
				}
				return fmt.Sprintf("%s()", function), fmt.Sprintf("%s:%d", filename, f.Line)
			}})
	}
	logger.SetLevel(level)
	switch strings.ToLower(cfg.Format) {
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{})
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	return logger, nil
}
