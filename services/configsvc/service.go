package configsvc

import (
	"context"
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"github.com/google/go-cmp/cmp"
	"github.com/knoguchi/go_project_template/services"
	"io/ioutil"
	"os"
)

func New() *ConfigSvc {
	cfg := &ConfigSvc{}
	cfg.Config = &ConfigSvcConfig{}
	cfg.ConfigChange = make(chan services.IServiceConfig)
	cfg.Key = "config"
	cfg.JConfig = &JsonConfig{}
	cfg.JConfig.Services = map[string]services.IServiceConfig{}
	return cfg
}

func (c *ConfigSvc) Configure() {
	c.LoadConfig(c.GetConfigFilePath())
}

func (c *ConfigSvc) Start(ctx context.Context) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	c.FsWatcher = watcher

	err = watcher.Add(c.GetConfigFilePath())
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				//log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					//log.Println("modified file:", event.Name)
					newCfg, err := c.ReadConfigFromFile(c.GetConfigFilePath())
					if err != nil {
						log.Error("can't read config file.  ignoring")
						// c.Status = "WARNING"
					} else {
						// config json is good
						for key := range newCfg.Services {
							currentCfg := c.Registry.GetCurrentConfig(key)
							//log.Infof("Checking config for [%s]", key)
							if diff := cmp.Diff(currentCfg, newCfg.Services[key]); diff != "" {
								log.Infof("Config for [%s] has changed", key)
								log.Infof("Current config: %#v", currentCfg)
								log.Infof("New     config: %#v", newCfg.Services[key])
								if err := c.Registry.NotifyConfigChange(key, newCfg.Services[key]); err != nil {
									log.Error(err)
								}
							}
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			case <-ctx.Done():
				log.Info("context done")
				return
			}
		}
	}()

	log.Info("config service started")
	return nil
}

func (c *ConfigSvc) Status() error {
	return nil
}
func (c *ConfigSvc) Stop() error {
	if err := c.FsWatcher.Close(); err != nil {
		return err
	}
	log.Info("config service stopped")
	return nil
}

// AddService /* The purpose of this function is to build a tree of config structs
func (c *ConfigSvc) AddService(svc services.IService) {
	svcCfg := svc.GetServiceConfig()
	c.JConfig.Services[svc.GetKey()] = svcCfg
}

// ReadConfigFromFile reads the configuration from the given file
// if target file is encrypted, prompts for encryption key
// Also - if not in dryrun mode - it checks if the configuration needs to be encrypted
func (c *ConfigSvc) ReadConfigFromFile(configPath string) (cfg *JsonConfig, err error) {
	//defaultPath, _, err := GetFilePath(configPath)
	//if err != nil {
	//	return err
	//}
	confFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer confFile.Close()
	byteValue, _ := ioutil.ReadAll(confFile)

	// fill fields except services
	result := &JsonConfig{}
	if err := json.Unmarshal(byteValue, result); err != nil {
		return nil, err
	}
	// unmarshal services
	_result := &_JsonConfig{}
	if err := json.Unmarshal(byteValue, _result); err != nil {
		return nil, err
	}
	result.Services = _result.Services
	return result, nil
}

// LoadConfig loads your configuration file into your configuration object
func (c *ConfigSvc) LoadConfig(configPath string) error {
	cfg, err := c.ReadConfigFromFile(configPath)
	if err != nil {
		log.Errorf(ErrFailureOpeningConfig, configPath, err)
		return err
	}
	c.JConfig = cfg
	return nil
}

func (c *ConfigSvc) GetConfigFilePath() string {
	return "./config.json"
}
