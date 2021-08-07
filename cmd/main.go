package main

import (
	"fmt"
	"github.com/caarlos0/env"
	"github.com/inview-team/raptor.stream-server/internal/app/connector"
	"github.com/inview-team/raptor.stream-server/internal/config"
	"github.com/inview-team/raptor.stream-server/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var conf config.Settings
	if err := env.Parse(&conf); err != nil {
		log.Fatal(err)
	}
	broadcaster := connector.New(&conf)
	srv := server.New(conf.StreamServerAddress, broadcaster)

	done := make(chan os.Signal, 1)
	errs := make(chan error, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		if err := srv.Stop(); err != nil {
			log.Fatal(fmt.Errorf("server stopped with error: %w", err))
			return
		}
	}()

	go func() {
		log.Printf("server started at %s", conf.StreamServerAddress)
		errs <- srv.Start()
	}()

	select {
	case <-done:
		signal.Stop(done)
		return
	case err := <-errs:
		if err != nil {
			log.Fatal("server exited with error: %w", err)
		}
		return
	}
}
