package webservice

import (
	"context"
	"errors"
	"github.com/knoguchi/go_project_template/services"
	"net/http"
)

type WebService struct {
	services.Service
}

func New() *WebService {
	s := &WebService{}
	s.ConfigChange = make(chan services.IServiceConfig, 100)
	s.Key = "webservice"
	return s
}

func (ws *WebService) Configure() {

}

func (ws *WebService) Start(ctx context.Context) error {
	log.Infof("%p start", ws.ConfigChange)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router01(ws.Registry),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
			//ws.ErrCh <- err
		}
	}()

	for {
		select {
		case newCfg := <-ws.ConfigChange:
			log.Infof("got config change notification: %v", newCfg)
			ws.MarkConfigTimestamp()
		case <-ctx.Done():
			log.Info("webservice shutting down")
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatal("Server forced to shutdown:", err)
				//ws.ErrCh <- err
			}
			break
		}
	}
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
