package relicconf

import (
	newrelic "github.com/newrelic/go-agent"
)

var app newrelic.Application

func GetRelicApp() newrelic.Application {
	return app
}

func InitNewRelic(appname string, lickey string) {
	var err error
	config := newrelic.NewConfig(appname, lickey)
	app, err = newrelic.NewApplication(config)
	if err != nil {
		panic(err)
	}
}
