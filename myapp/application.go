package myapp

import (
	"context"
	"github.com/knoguchi/go_project_template/services"
	"github.com/knoguchi/go_project_template/services/configsvc"
	"github.com/knoguchi/go_project_template/services/database"
	"github.com/knoguchi/go_project_template/services/kafka"
	"github.com/knoguchi/go_project_template/services/myservice"
	"github.com/knoguchi/go_project_template/services/webservice"
	"github.com/knoguchi/go_project_template/version"
	"github.com/sirupsen/logrus"
	"sync"
)

type Application struct {
	lock     sync.RWMutex
	registry *services.ServiceRegistry
	db       database.Database
}

// New starts a new myapp
func New(configPath string) (*Application, error) {
	registry := services.NewServiceRegistry()
	app := &Application{

		registry: registry,
	}

	if err := app.registerServices(); err != nil {
		return nil, err
	}

	app.registry.ConfigureAll()

	return app, nil
}

func (app *Application) Start(parentCtx context.Context) error {

	// use config service, and setup tracing.  that way cliCtx is not necessary
	//if err := tracing.Setup(
	//	"myapp", // FIXME: change the name
	//	cliCtx.String(cmd.TracingProcessNameFlag.Name),
	//	cliCtx.String(cmd.TracingEndpointFlag.Name),
	//	cliCtx.Float64(cmd.TraceSampleFractionFlag.Name),
	//	cliCtx.Bool(cmd.EnableTracingFlag.Name),
	//); err != nil {
	//	return err
	//}

	ctx, cancel := context.WithCancel(parentCtx)
	// ---- main stuff
	app.lock.Lock()
	err := app.registry.StartAll(ctx)
	app.lock.Unlock()
	if err != nil {
		log.Errorf("Startup failed: %v", err)
		cancel()
		return err
	}
	log.WithFields(logrus.Fields{"version": version.Version}).Info("Started")
	return nil
}

//// Close handles graceful shutdown of the system.
//func (app *Application) Close() {
//	app.lock.Lock()
//	defer app.lock.Unlock()
//
//	log.Info("Stopping myapp")
//	app.registry.StopAll()
//	//if err := app.database.Close(); err != nil {
//	//	log.Errorf("Failed to close database: %v", err)
//	//}
//	log.Info("closed")
//}

func (app *Application) registerServices() (err error) {
	cfgsvc := configsvc.New()
	err = app.registry.RegisterService(cfgsvc)
	if err != nil {
		panic(err)
		return
	}
	cfgsvc.AddService(cfgsvc)

	websvc := webservice.New()
	err = app.registry.RegisterService(websvc)
	if err != nil {
		panic(err)
		return
	}
	cfgsvc.AddService(websvc)

	mysvc := myservice.New()
	err = app.registry.RegisterService(mysvc)
	if err != nil {
		return
	}
	cfgsvc.AddService(mysvc)

	kafkasvc := kafka.New()
	err = app.registry.RegisterService(kafkasvc)
	if err != nil {
		return err
	}
	cfgsvc.AddService(kafkasvc)

	return nil
}
