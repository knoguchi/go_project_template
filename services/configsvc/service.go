package configsvc

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/fsnotify/fsnotify"
	"github.com/google/go-cmp/cmp"
	"github.com/knoguchi/go_project_template/services"
	"io"
	"io/ioutil"
	"os"
)

func New() *ConfigSvc {
	cfg := &ConfigSvc{}
	cfg.ConfigChange = make(chan services.IServiceConfig)
	cfg.Key = "config"
	cfg.MainConfig = &MainConfig{}
	cfg.MainConfig.Services = map[string]services.IServiceConfig{}
	return cfg
}

func (c *ConfigSvc) Configure() {
	log.Info("implement me")
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
					newCfg, err := c.ReadConfigFromFile(c.GetConfigFilePath(), false)
					if err != nil {
						log.Error("can't read config file.  ignoring")
						// c.Status = "WARNING"
					} else {
						// config json is good
						for key := range newCfg.Services {
							if diff := cmp.Diff(c.MainConfig.Services[key], newCfg.Services[key]); diff != "" {
								log.Infof("Config for [%s] has changed: %s", key, diff)
								c.MainConfig = newCfg
								c.Registry.NotifyConfigChange(key, newCfg.Services[key])
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
		log.Error("fs check exit")
	}()

	log.Info("config service started")
	return nil
}

func (c *ConfigSvc) Status() error {
	return nil
}
func (c *ConfigSvc) Stop() error {
	c.FsWatcher.Close()
	log.Info("config service stopped")
	return nil
}

/* The purpose of this function is to build a tree of config structs
Then let json.Marshal to populate
*/
func (c *ConfigSvc) AddService(svc services.IService) {
	svcCfg := svc.GetServiceConfig()
	c.MainConfig.Services[svc.GetKey()] = svcCfg
}

// ReadConfig verifies and checks for encryption and loads the config from a JSON object.
// Prompts for decryption key, if target data is encrypted.
// Returns the loaded configuration and whether it was encrypted.
func ReadConfig(configReader io.Reader) (*MainConfig, error) {
	reader := bufio.NewReader(configReader)

	// Read unencrypted configuration
	decoder := json.NewDecoder(reader)
	c := &MainConfig{}
	err := decoder.Decode(c)
	return c, err
}

// ReadConfigFromFile reads the configuration from the given file
// if target file is encrypted, prompts for encryption key
// Also - if not in dryrun mode - it checks if the configuration needs to be encrypted
func (c *ConfigSvc) ReadConfigFromFile(configPath string, dryrun bool) (cfg *MainConfig, err error) {
	//defaultPath, _, err := GetFilePath(configPath)
	//if err != nil {
	//	return err
	//}
	confFile, err := os.Open(c.GetConfigFilePath())
	if err != nil {
		return nil, err
	}
	defer confFile.Close()
	byteValue, _ := ioutil.ReadAll(confFile)

	// fill fields except services
	result := &MainConfig{}
	if err := json.Unmarshal(byteValue, result); err != nil {
		return nil, err
	}
	// unmarshal services
	_result := &_MainConfig{}
	if err := json.Unmarshal(byteValue, _result); err != nil {
		return nil, err
	}
	result.Services = _result.Services
	return result, nil
}

// LoadConfig loads your configuration file into your configuration object
func (c *ConfigSvc) LoadConfig(configPath string, dryrun bool) error {
	cfg, err := c.ReadConfigFromFile(configPath, dryrun)
	if err != nil {
		log.Errorf(ErrFailureOpeningConfig, configPath, err)
		return err
	}
	c.MainConfig = cfg
	return nil
}

func (c *ConfigSvc) GetConfigFilePath() string {
	return "./config.json"
}
