package main

import (
	"os"
	"os/signal"
	"syscall"

	todo "github.com/ExperienceOne/apikit/example"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var todoService *todo.Service

func main() {

	app := cli.NewApp()
	app.Name = "Todo Service"

	app.Action = func(ctx *cli.Context) error {

		log.SetFormatter(&log.JSONFormatter{})

		service := todo.NewService()
		todoService = service

		return service.Start(9001)
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigint

		if todoService != nil {
			if err := todoService.Stop(); err != nil {
				log.WithError(err).Error("failed to stop todo service")
				os.Exit(-1)
			}
		}

		os.Exit(0)
	}()

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Error("failed to start todo service")
	}
}
