package myservice

import (
	"context"
	"errors"
	"github.com/knoguchi/go_project_template/services"
	"sync"
	"sync/atomic"
	"time"
)

type MyService struct {
	services.Service
	started int32
	stopped int32
	Mu      sync.RWMutex
	count   int
}

func New() *MyService {
	s := &MyService{}
	s.Key = "myservice"
	s.Config = &MyServiceConfig{}
	s.ConfigChange = make(chan services.IServiceConfig)
	return s
}

func (a *MyService) Configure() {

}

func (a *MyService) Start(ctx context.Context) error {
	log.Debugln("services starting...")
	a.count = 0
	go a.run(ctx)
	return nil
}

func (a *MyService) Status() (err error) {
	return nil
}

func (a *MyService) _Status() bool {
	return atomic.LoadInt32(&a.started) == 1
}

/*
Start services.  This function wraps run() and make sure only one
*/
func (a *MyService) _Start(ctx context.Context) (err error) {
	if atomic.AddInt32(&a.started, 1) != 1 {
		return errors.New("services already started")
	}

	defer func() {
		if err != nil {
			atomic.CompareAndSwapInt32(&a.started, 1, 0)
		}
	}()

	log.Debugln("services starting...")
	a.count = 0
	go a.run(ctx)

	return err
}

func (a *MyService) Stop() error {
	return nil
	//if atomic.LoadInt32(&a.started) == 0 {
	//	return errors.New("services not started")
	//}
	//if atomic.AddInt32(&a.stopped, 1) != 1 {
	//	return errors.New("services is already stopping")
	//}
	//
	//return nil
}

func (a *MyService) run(ctx context.Context) {
	log.Debugln("my services started.")
	//a.wg.Add(1)

	t := time.NewTicker(time.Second * 2)

	defer func() {
		t.Stop()
		atomic.CompareAndSwapInt32(&a.stopped, 1, 0)
		atomic.CompareAndSwapInt32(&a.started, 1, 0)
		log.Debugln("services shutdown.")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			a.checkConnection()
		}
	}
}

func (a *MyService) checkConnection() {
	a.Mu.Lock()
	defer a.Mu.Unlock()

	//err := x.SQL.Ping()
	//if err != nil {
	//	log.Errorf( "connection error: %v\n", err)
	//	x.Connected = false
	//	return
	//}
	log.Info("ok!")
	//if !x.Connected {
	//	log.Info("connection reestablished")
	//	x.Connected = true
	//}
	a.count += 1
}

func (a *MyService) GetCount() interface{} {
	return a.count
}
