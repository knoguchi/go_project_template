package services

import (
	"context"
	"time"
)

type IServiceConfig interface {
}

type ServiceConfig struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Verbose bool   `json:"verbose"`
}

// IService is a struct that can be registered into a ServiceRegistry for
// easy dependency management.
type IService interface {
	// Start spawns any goroutines required by the services.
	Start(context.Context) error
	// Stop terminates all goroutines belonging to the services,
	// blocking until they are all terminated.
	Stop() error
	// Status Returns error if the services is not considered healthy.
	Status() error
	SetRegistry(registry *ServiceRegistry)
	GetServiceConfig() IServiceConfig
	Configure()
	MarkConfigTimestamp()
	GetKey() string
	ChangeConfig(config IServiceConfig)
}

type Service struct {
	Key             string
	Registry        *ServiceRegistry
	Verbose         bool
	Config          IServiceConfig
	ConfigChange    chan IServiceConfig
	ConfigTimestamp time.Time
}

func (s *Service) SetRegistry(registry *ServiceRegistry) {
	s.Registry = registry
}

func (s *Service) GetServiceConfig() IServiceConfig {
	if s.Config != nil {
		return s.Config
	}
	log.Error("Config should never be nil")
	return nil
}

func (s *Service) GetKey() string {
	return s.Key
}

func (s *Service) ChangeConfig(config IServiceConfig) {
	log.Infof("%p %s: pushing new config to channel: %#v", s.ConfigChange, s.Key, config)
	s.ConfigChange <- config
	log.Infof("pushed")
}

func (s *Service) MarkConfigTimestamp() {
	s.ConfigTimestamp = time.Now()
}
