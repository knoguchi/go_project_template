package services

import (
	"context"
	"fmt"
	"reflect"
	"time"
)

type IServiceConfig interface {
	//GetName() string
}

type ServiceConfig struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Verbose bool   `json:"verbose"`
}

//func (sc *ServiceConfig) GetName() string {
//	return sc.Name
//}

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
	Config          *IServiceConfig
	ConfigChange    chan IServiceConfig
	ConfigTimestamp time.Time
}

func (s *Service) SetRegistry(registry *ServiceRegistry) {
	s.Registry = registry
}

func (s *Service) GetServiceConfig() IServiceConfig {
	if s.Config != nil {
		return *s.Config
	}
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

// ServiceRegistry provides a useful pattern for managing services.
// It allows for ease of dependency management and ensures services
// dependent on others use the same references in memory.
type ServiceRegistry struct {
	services      map[reflect.Type]IService // map of types to services.
	servicesByKey map[string]*IService      // map of keys to services
	serviceTypes  []reflect.Type            // keep an ordered slice of registered services types.
}

// NewServiceRegistry starts a registry instance for convenience
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services:      make(map[reflect.Type]IService),
		servicesByKey: make(map[string]*IService),
	}
}

// StartAll initialized each services in order of registration.
func (s *ServiceRegistry) ConfigureAll() {
	log.Infof("Configuring %d services: %v", len(s.serviceTypes), s.serviceTypes)
	for _, kind := range s.serviceTypes {
		log.Debugf("Configure services type %v", kind)
		go s.services[kind].Configure()
	}
}

// StartAll initialized each services in order of registration.
func (s *ServiceRegistry) StartAll(ctx context.Context) {
	log.Infof("Starting %d services: %v", len(s.serviceTypes), s.serviceTypes)
	for _, kind := range s.serviceTypes {
		log.Infof("Starting services type %v", kind)
		go s.services[kind].Start(ctx)
	}
}

// StopAll ends every services in reverse order of registration, logging a
// panic if any of them fail to stop.
func (s *ServiceRegistry) StopAll() {
	for i := len(s.serviceTypes) - 1; i >= 0; i-- {
		kind := s.serviceTypes[i]
		service := s.services[kind]
		if err := service.Stop(); err != nil {
			log.Panicf("Could not stop the following services: %v, %v", kind, err)
		}
	}
}

// Statuses returns a map of IService type -> error. The map will be populated
// with the results of each services.Status() method call.
func (s *ServiceRegistry) Statuses() map[reflect.Type]error {
	m := make(map[reflect.Type]error)
	for _, kind := range s.serviceTypes {
		m[kind] = s.services[kind].Status()
	}
	return m
}

// RegisterService appends a services constructor function to the services
// registry.
func (s *ServiceRegistry) RegisterService(service IService) error {
	// attach itself
	service.SetRegistry(s)

	kind := reflect.TypeOf(service)
	if _, exists := s.services[kind]; exists {
		return fmt.Errorf("services already exists: %v", kind)
	}
	s.services[kind] = service
	s.servicesByKey[service.GetKey()] = &service
	s.serviceTypes = append(s.serviceTypes, kind)
	log.Infof("Registered %s", kind.String())
	return nil
}

// FetchService takes in a struct pointer and sets the value of that pointer
// to a services currently stored in the services registry. This ensures the input argument is
// set to the right pointer that refers to the originally registered services.
func (s *ServiceRegistry) FetchService(service interface{}) error {
	if reflect.TypeOf(service).Kind() != reflect.Ptr {
		return fmt.Errorf("input must be of pointer type, received value type instead: %T", service)
	}
	element := reflect.ValueOf(service).Elem()
	if running, ok := s.services[element.Type()]; ok {
		element.Set(reflect.ValueOf(running))
		return nil
	}
	return fmt.Errorf("unknown services: %T", service)
}

func (s *ServiceRegistry) NotifyConfigChange(key string, cfg IServiceConfig) error {
	if running, ok := s.servicesByKey[key]; ok {
		(*running).ChangeConfig(cfg)
		return nil
	}
	return fmt.Errorf("unknown service key: %s", key)
}
