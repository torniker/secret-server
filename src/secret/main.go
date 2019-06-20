package main

import (
	"secret/app"
	"secret/handler"
	"secret/log"
)

func main() {
	go setupLogChannel()
	log.Info("service is starting...")

	a := app.New()
	a.DefaultHandler = handler.Default
	err := a.StartHTTP(":8080")
	if err != nil {
		log.WithError(err).Error("could not start server")
	}
}

func setupLogChannel() {
	ch := log.NewChan()
	for {
		_ = <-ch
		// TODO: handle log store somewhere
	}
}
