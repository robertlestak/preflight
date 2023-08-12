package main

import (
	"flag"
	"os"

	"github.com/robertlestak/preflight/pkg/preflight"
	log "github.com/sirupsen/logrus"
)

func init() {
	ll, err := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
}

func main() {
	l := log.WithFields(log.Fields{
		"app": "preflight",
	})
	l.Debug("starting preflight")
	preflightFlags := flag.NewFlagSet("preflight", flag.ExitOnError)
	configFile := preflightFlags.String("config", "preflight.yaml", "path to config file")
	logLevel := preflightFlags.String("log-level", log.GetLevel().String(), "log level")
	concurrency := preflightFlags.Int("concurrency", 1, "number of concurrent checks to run")
	preflightFlags.Parse(os.Args[1:])
	ll, err := log.ParseLevel(*logLevel)
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
	preflight.SetLogger(l.Logger)
	l.WithField("log-level", ll).Debug("log level set")
	l.WithField("config", *configFile).Debug("loading config file")
	p, err := preflight.LoadConfig(*configFile)
	if err != nil {
		l.WithError(err).Fatal("unable to load config file")
	}
	if p.Concurrency == 0 {
		p.Concurrency = *concurrency
	}
	l.Debug("running preflight checks")
	if err := p.Run(); err != nil {
		l.WithError(err).Fatal("preflight checks failed")
	}
	l.Info("preflight checks passed")
}
