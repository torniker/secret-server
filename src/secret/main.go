package main

import (
	"secret/app"
	"secret/handler"
	"secret/log"

	"github.com/prometheus/client_golang/prometheus"
)

func main() {
	go setupLogChannel()
	log.Info("service is starting...")

	a := app.New()
	a.DefaultHandler = handler.Default
	setupPrometheus(a)
	err := a.StartHTTP(":8080")
	if err != nil {
		log.WithError(err).Error("could not start server")
	}
}

func setupPrometheus(a *app.App) {
	// TODO: need to think of a better way of tracking routes and
	// registering appropriate prometheus object
	a.Summery["addSecret"] = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "secret",
			Name:      "addSecret",
		},
		[]string{"secret"},
	)
	prometheus.Register(a.Summery["addSecret"])
	a.Counter["addSecret"] = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "secret",
		Name:      "addSecret",
	})

	a.Summery["getSecretByHash"] = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "secret",
			Name:      "getSecretByHash",
		},
		[]string{"secret"},
	)
	prometheus.Register(a.Summery["getSecretByHash"])
	a.Counter["getSecretByHash"] = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: "secret",
		Name:      "getSecretByHash",
	})
}

func setupLogChannel() {
	ch := log.NewChan()
	for {
		_ = <-ch
		// TODO: handle log store somewhere
	}
}
