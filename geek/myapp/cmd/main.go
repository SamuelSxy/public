package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type App struct {
	ctx    context.Context
	cancel func()
	http   http.Server
}

func New(addr string, handler http.Handler) *App {

	http := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	return &App{
		ctx:    ctx,
		cancel: cancel,
		http:   http,
	}

}

func (a *App) Run() {
	eg, ctx := errgroup.WithContext(a.ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSTOP)

	eg.Go(func() error {
		<-ctx.Done()
		return a.http.Shutdown(ctx)
	})

	eg.Go(func() error {
		return a.http.ListenAndServe()
	})

	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-sigChan:
			return a.Stop()
		}

	})

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
	}
}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil

}

func NewApp() *App {

	http := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service, _ := initApp()
		fmt.Fprintln(w, "Http START")
		fmt.Fprintln(w, service.Name)
		fmt.Fprintln(w, "Http END")
	})

	return New("0.0.0.0:8080", http)

}

func main() {
	app := NewApp()
	app.Run()
}
