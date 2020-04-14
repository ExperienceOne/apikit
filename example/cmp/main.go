package main

import (
	"os"
	"os/signal"
	"syscall"

	basket "github.com/ExperienceOne/apikit/example"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var basketService *basket.Service

func main() {

	app := cli.NewApp()
	app.Name = "Basket Service"

	app.Action = func(ctx *cli.Context) error {

		log.SetFormatter(&log.JSONFormatter{})

		service := basket.NewService()
		basketService = service

		return service.Start(9001)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigint

		if basketService != nil {
			if err := basketService.Stop(); err != nil {
				log.WithError(err).Error("failed to stop Basket Service")
				os.Exit(-1)
			}
		}

		os.Exit(0)
	}()

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Error("failed to start Basket Service")
	}
}
