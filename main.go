package main

import (
	"fmt"
	"maintenance-tasks/api"
	"net/http"
	"sync"
)

type handler struct {
	apiHandler http.Handler
}

func main() {
	apiHandler, err := api.CreateHandler()
	if err != nil {
		panic(err)
	}

	h := handler{
		apiHandler: apiHandler,
	}

	srv := &http.Server{Addr: ":3000", Handler: h}
	wg := sync.WaitGroup{}
	wg.Add(1)
	fmt.Println("Starting server")
	go startServer(&wg, srv)

	wg.Wait()
}

func startServer(wg *sync.WaitGroup, srv *http.Server) {
	defer wg.Done()
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		panic(err)
	}
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.apiHandler.ServeHTTP(w, r)
}
