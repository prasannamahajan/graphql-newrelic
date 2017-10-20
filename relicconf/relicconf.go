package relicconf

import (
	newrelic "github.com/newrelic/go-agent"
)

var app newrelic.Application

func GetRelicApp() newrelic.Application {
	return app
}

func InitNewRelic() {
	var err error
	const licKey = "newrelickey"
	config := newrelic.NewConfig("tut4", licKey)
	app, err = newrelic.NewApplication(config)
	if err != nil {
		panic(err)
	}
}
