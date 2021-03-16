package webservice

import (
	"context"
	"errors"
	"github.com/knoguchi/go_project_template/services"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type WebService struct {
	services.Service
}

func New() *WebService {
	s := &WebService{}
	s.Config = &WebServiceConfig{}
	s.ConfigChange = make(chan services.IServiceConfig, 100)
	s.Key = "webservice"
	return s
}

func (ws *WebService) Configure() {

}

func (ws *WebService) Start(ctx context.Context) error {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router01(ws.Registry),
	}

	g, gctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
			return err
		}
		return nil
	})

	g.Go(func() error {
		for {
			select {
			case newCfg := <-ws.ConfigChange:
				log.Infof("got config change notification: %v", newCfg)
				ws.Config = &newCfg
				ws.MarkConfigTimestamp()
			case <-gctx.Done():
				log.Info("ctx done")
				if err := srv.Shutdown(ctx); err != nil {
					log.Infof("Server forced to shutdown: %v", err)
					return gctx.Err()
				}
				return nil
			}
		}
	})

	// this is bad
	go func() {
		err := g.Wait()
		if err != nil {
			log.Infof("webservice error: %v", err)
			return
		}
		return
	}()

	log.Info("webservice started")
	return nil
}

//func (ws *WebService) GetName() string {
//	return "webservice"
//}

func (ws *WebService) Stop() error {
	return nil
}

func (ws *WebService) Status() error {
	return nil
}
