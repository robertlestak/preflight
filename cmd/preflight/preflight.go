package main

import (
	"encoding/json"
	"flag"
	"net/http"
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

func server(addr string, remoteToken string) error {
	l := log.WithFields(log.Fields{
		"app": "preflight",
	})
	l.Debug("starting preflight server")
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			l.WithField("method", r.Method).Warn("invalid method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if remoteToken != "" {
			if r.Header.Get("Authorization") != "token "+remoteToken {
				l.Warn("invalid token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		pf := &preflight.Preflight{}
		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(pf); err != nil {
			l.WithError(err).Error("unable to parse request body")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		pf.Remote = ""
		if err := pf.Run(); err != nil {
			l.WithError(err).Error("preflight checks failed")
			w.WriteHeader(http.StatusInternalServerError)
			// send error message
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	l.WithField("addr", addr).Info("listening")
	return http.ListenAndServe(addr, nil)
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
	serverMode := preflightFlags.Bool("server", false, "run in server mode")
	serverAddr := preflightFlags.String("server-addr", ":8090", "server listen address")
	serverToken := preflightFlags.String("server-token", "", "token to use when running in server mode")
	remote := preflightFlags.String("remote", "", "remote preflight server to run checks against")
	remoteToken := preflightFlags.String("remote-token", "", "token to use when running remote checks")
	equiv := preflightFlags.Bool("equiv", false, "print equivalent command")
	preflightFlags.Parse(os.Args[1:])
	ll, err := log.ParseLevel(*logLevel)
	if err != nil {
		ll = log.InfoLevel
	}
	log.SetLevel(ll)
	preflight.SetLogger(l.Logger)
	l.WithField("log-level", ll).Debug("log level set")
	l.WithField("config", *configFile).Debug("loading config file")
	if *serverMode {
		if err := server(*serverAddr, *serverToken); err != nil {
			l.WithError(err).Fatal("server failed")
		}
		return
	}
	p, err := preflight.LoadConfig(*configFile)
	if err != nil {
		l.WithError(err).Fatal("unable to load config file")
	}
	if p.Concurrency == 0 {
		p.Concurrency = *concurrency
	}
	if p.Remote == "" && *remote != "" {
		p.Remote = *remote
	}
	if p.RemoteToken == "" && *remoteToken != "" {
		p.RemoteToken = *remoteToken
	}
	if *equiv {
		p.Equivalent = true
	}
	l.Debug("running preflight checks")
	if err := p.Run(); err != nil {
		l.WithError(err).Fatal("preflight checks failed")
	}
	if !*equiv {
		l.Info("preflight checks passed")
	}
}
