package myservice

import (
	"github.com/knoguchi/go_project_template/services"
	"sync"
)

// Instance holds all information for a services instance
type Instance struct {
	//DB       *services
	Config    *Config
	Connected bool
	Mu        sync.RWMutex
}

// Config holds all services configurable options including enable/disabled & settings
type Config struct {
	services.ServiceConfig
	Interval int
}
