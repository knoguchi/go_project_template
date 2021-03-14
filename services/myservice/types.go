package myservice

import (
	"github.com/knoguchi/go_project_template/services"
	"sync"
)

// Instance holds all information for a services instance
type Instance struct {
	//DB       *services
	MyServiceConfig *MyServiceConfig
	Connected       bool
	Mu              sync.RWMutex
}

// MyServiceConfig holds all services configurable options including enable/disabled & settings
type MyServiceConfig struct {
	services.ServiceConfig
	Interval int
}
