package configsvc

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/knoguchi/go_project_template/services"
	"io"
	"io/ioutil"
	"os"
)

func New() *ConfigSvc {
	cfg := &ConfigSvc{}
	cfg.Config.Services = map[string]services.IServiceConfig{}
	return cfg
}

func (c *ConfigSvc) Configure() {
	log.Info("implement me")
}

func (c *ConfigSvc) Start(ctx context.Context) error {
	return nil
}

func (c *ConfigSvc) Status() error {
	return nil
}
func (c *ConfigSvc) Stop() error {
	return nil
}

/* The purpose of this function is to build a tree of config structs
Then let json.Marshal to populate
 */
func (c *ConfigSvc) AddServiceConfig(cfg services.IServiceConfig) {
	key := cfg.GetName()
	c.Config.Services[key] = cfg
}

// ReadConfig verifies and checks for encryption and loads the config from a JSON object.
// Prompts for decryption key, if target data is encrypted.
// Returns the loaded configuration and whether it was encrypted.
func ReadConfig(configReader io.Reader) (*Config, error) {
	reader := bufio.NewReader(configReader)

	// Read unencrypted configuration
	decoder := json.NewDecoder(reader)
	c := &Config{}
	err := decoder.Decode(c)
	return c, err
}

// ReadConfigFromFile reads the configuration from the given file
// if target file is encrypted, prompts for encryption key
// Also - if not in dryrun mode - it checks if the configuration needs to be encrypted
func (c *ConfigSvc) ReadConfigFromFile(configPath string, dryrun bool) (err error) {
	//defaultPath, _, err := GetFilePath(configPath)
	//if err != nil {
	//	return err
	//}
	confFile, err := os.Open("./config.json")
	if err != nil {
		return err
	}
	defer confFile.Close()
	byteValue, _ := ioutil.ReadAll(confFile)
	result := Config{}
	result.Services = map[string]services.IServiceConfig{}
	log.Info("%s", string(byteValue))
	log.Infof("%v", c.Config)
	//result.Services = []services.IServiceConfig{}
	err = json.Unmarshal(byteValue, &c.Config)
	if err != nil {
		return err
	}
	log.Infof("%v", result)
	//// Override values in the current Config
	//*c.Config = *result
	return nil
}

// CheckConfig checks all config settings
func (c *ConfigSvc) CheckConfig() error {
	return nil
}

// LoadConfig loads your configuration file into your configuration object
func (c *ConfigSvc) LoadConfig(configPath string, dryrun bool) error {
	err := c.ReadConfigFromFile(configPath, dryrun)
	if err != nil {
		log.Errorf(ErrFailureOpeningConfig, configPath, err)
		return err
	}

	return c.CheckConfig()
}
