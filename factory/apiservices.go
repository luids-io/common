// Copyright 2019 Luis Guillén Civera <luisguillenc@gmail.com>. View LICENSE.

package factory

import (
	"fmt"

	"github.com/luisguillenc/yalogi"

	"github.com/luids-io/common/config"
	"github.com/luids-io/common/util"
	"github.com/luids-io/core/apiservice"
)

// APIServices is a factory
func APIServices(cfg *config.APIServicesCfg, logger yalogi.Logger) (*apiservice.Registry, error) {
	if cfg.Empty() {
		return apiservice.NewRegistry(), nil
	}
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	defs, err := getServiceDefs(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("loading servicedefs: %v", err)
	}
	services := apiservice.NewRegistry()
	for _, def := range defs {
		if def.Disabled {
			logger.Debugf("'%s' is disabled", def.ID)
			continue
		}
		svc, err := apiservice.Build(def, logger)
		if err != nil {
			return nil, fmt.Errorf("building '%s': %v", def.ID, err)
		}
		services.Register(def.ID, svc)
	}
	return services, nil
}

// APIAutoloader is a factory
func APIAutoloader(cfg *config.APIServicesCfg, logger yalogi.Logger) (*apiservice.Autoloader, error) {
	if cfg.Empty() {
		return apiservice.NewAutoloader(
			[]apiservice.Definition{},
			apiservice.SetLogger(logger),
		), nil
	}
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("bad config: %v", err)
	}
	defs, err := getServiceDefs(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("loading servicedefs: %v", err)
	}
	return apiservice.NewAutoloader(defs, apiservice.SetLogger(logger)), nil
}

func getServiceDefs(cfg *config.APIServicesCfg, logger yalogi.Logger) ([]apiservice.Definition, error) {
	dbFiles, err := util.GetFilesDB("json", cfg.ConfigFiles, cfg.ConfigDirs)
	if err != nil {
		return nil, err
	}
	loadedDB := make([]apiservice.Definition, 0)
	for _, file := range dbFiles {
		logger.Debugf("loading file '%s'", file)
		entries, err := apiservice.DefsFromFile(file)
		if err != nil {
			return nil, fmt.Errorf("loading file '%s': %v", file, err)
		}
		loadedDB = append(loadedDB, entries...)
	}
	return loadedDB, nil
}
