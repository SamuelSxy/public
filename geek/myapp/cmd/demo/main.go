package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSTOP)

	serv := http.Server{Addr: ":8080"}
	http.HandleFunc("/", helloWord)

	g.Go(func() error {
		go func() {
			<-ctx.Done()
			serv.Shutdown(context.TODO())
		}()
		return serv.ListenAndServe()
	})

	g.Go(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case s := <-sigChan:
			return errors.Errorf("signal: %v", s)
		}
	})

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}

}

func helloWord(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello golang http!")
}
