package main

import (
	"context"
	"fmt"
	"maintenance-tasks/api"
	"maintenance-tasks/manager"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type handler struct {
	apiHandler http.Handler
}

func main() {
	manager := manager.Create()
	manager.Start()

	apiHandler, err := api.CreateHandler(manager)
	if err != nil {
		panic(err)
	}

	h := handler{
		apiHandler: apiHandler,
	}

	srv := &http.Server{Addr: ":3000", Handler: h}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go startServer(&wg, srv)

	stopServer(srv)
	wg.Wait()
}

func startServer(wg *sync.WaitGroup, srv *http.Server) {
	fmt.Println("Starting server")
	defer wg.Done()
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

func stopServer(srv *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	fmt.Println("Stopping http server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}

}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.apiHandler.ServeHTTP(w, r)
}
