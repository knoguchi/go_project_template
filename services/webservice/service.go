package webservice

import (
	"context"
	"errors"
	"github.com/knoguchi/go_project_template/services"
	"net/http"
	"time"
)

type WebService struct {
	services.Service
	srv *http.Server
}

func New() *WebService {
	s := &WebService{}
	return s
}

func (ws *WebService) Configure() {

}

func (ws *WebService) Start(ctx context.Context) error {
	ws.srv = &http.Server{
		Addr:    ":8080",
		Handler: router01(ws.Registry),
	}

	go func() {
		if err := ws.srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
			//ws.ErrCh <- err
		}

		// allow 5sec for http server to shutdown
		srvCtx, srvCancel := context.WithTimeout(ctx, 5*time.Second)
		defer srvCancel()
		for {
			select {
			case <-ctx.Done():
				log.Info("webservice shutting down")
				if err := ws.srv.Shutdown(srvCtx); err != nil {
					log.Fatal("Server forced to shutdown:", err)
					//ws.ErrCh <- err
				}
				break
			}
		}
		log.Info("webservice stopped")
	}()

	log.Info("webservice started")
	return nil
}

func (ws *WebService) GetName() string {
	return "REST API"
}

func (ws *WebService) Stop() error {
	return nil
}

func (ws *WebService) Status() error {
	return nil
}
