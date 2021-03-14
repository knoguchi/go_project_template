package myapp

import (
	"context"
	"github.com/knoguchi/go_project_template/services"
	"github.com/knoguchi/go_project_template/services/configsvc"
	"github.com/knoguchi/go_project_template/services/database"
	"github.com/knoguchi/go_project_template/services/webservice"
	"github.com/knoguchi/go_project_template/version"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Application struct {
	lock     sync.RWMutex
	registry *services.ServiceRegistry
	db       database.Database
}

// New starts a new myapp
func New() (*Application, error) {
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

func (app *Application) Start(ctx context.Context) error {

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

	g, gctx := errgroup.WithContext(ctx)

	// goroutine to check for signals to gracefully finish all functions
	g.Go(func() error {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

		select {
		case sig := <-signalChannel:
			log.Info("Received signal: %s\n", sig)
			gctx.Done()
		case <-gctx.Done():
			log.Info("closing signal goroutine\n")
			return gctx.Err()
		}

		return nil
	})

	// ---- main stuff
	app.lock.Lock()
	app.registry.StartAll(gctx)
	app.lock.Unlock()

	log.WithFields(logrus.Fields{"version": version.Version}).Info("Starting app")

	return g.Wait()
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
		return
	}

	websvc := webservice.New()
	err = app.registry.RegisterService(websvc)
	if err != nil {
		return
	}
	cfgsvc.AddServiceConfig(websvc)

	//mysvc := myservice.New()
	//err = app.registry.RegisterService(mysvc)
	//if err != nil {
	//	return
	//}
	//cfgsvc.AddServiceConfig(mysvc.GetServiceConfig())

	//kafkasvc := kafka.New()
	//err = app.registry.RegisterService(kafkasvc)
	//if err != nil {
	//	return err
	//}
	//cfgsvc.AddServiceConfig(kafkasvc.GetServiceConfig())

	cfgsvc.LoadConfig("./config.json", false)

	return nil
}
