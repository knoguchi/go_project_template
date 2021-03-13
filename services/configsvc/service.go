package configsvc

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/knoguchi/go_project_template/services"
	"io"
	"os"
)

func New() *ConfigSvc {
	cfg := &ConfigSvc{}
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

func (c *ConfigSvc) AddServiceConfig(cfg services.IServiceConfig) {
	c.config.Services = append(c.config.Services, cfg)
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

	reader := bufio.NewReader(confFile)
	decoder := json.NewDecoder(reader)
	result := &Config{}
	err = decoder.Decode(c)
	if err != nil {
		return fmt.Errorf("error reading config %w", err)
	}

	// Override values in the current config
	*c.Config = *result
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
