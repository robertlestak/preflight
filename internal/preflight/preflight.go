package preflight

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/robertlestak/preflight-dns/pkg/preflightdns"
	"github.com/robertlestak/preflight-env/pkg/preflightenv"
	"github.com/robertlestak/preflight-id/pkg/preflightid"
	"github.com/robertlestak/preflight-netpath/pkg/preflightnetpath"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type PreflightDriverName string

const (
	DriverNameDNS     PreflightDriverName = "dns"
	DriverNameEnv     PreflightDriverName = "env"
	DriverNameID      PreflightDriverName = "id"
	DriverNameNetPath PreflightDriverName = "netpath"
)

type Preflight struct {
	Concurrency int                                 `json:"concurrency" yaml:"concurrency"`
	DNS         []preflightdns.PreflightDNS         `json:"dns" yaml:"dns"`
	Env         map[string]string                   `json:"env" yaml:"env"`
	ID          []preflightid.PreflightID           `json:"id" yaml:"id"`
	Netpath     []preflightnetpath.PreflightNetPath `json:"netpath" yaml:"netpath"`
}

func SetLogger(l *log.Logger) {
	preflightdns.Logger = l
	preflightenv.Logger = l
	preflightid.Logger = l
	preflightnetpath.Logger = l
}

func LoadConfig(filepath string) (*Preflight, error) {
	l := log.WithFields(log.Fields{
		"package": "preflight",
		"fn":      "LoadConfig",
	})
	l.Debug("loading config file")
	p := Preflight{}
	fd, err := os.ReadFile(filepath)
	if err != nil {
		l.WithError(err).Error("unable to read config file")
		return nil, err
	}
	err = yaml.Unmarshal(fd, &p)
	if err != nil {
		// try json
		l.Debug("unable to parse as yaml, trying json")
		err = json.Unmarshal(fd, &p)
		if err != nil {
			l.WithError(err).Error("unable to parse config file")
			return nil, err
		}
	}
	l.Debugf("loaded config file: %+v", p)
	return &p, nil
}

type PreflightJob struct {
	Driver PreflightDriverName `json:"driver" yaml:"driver"`
	Job    any                 `json:"job" yaml:"job"`
}

func preflightDriverWorker(jobs chan PreflightJob, res chan error) {
	l := log.WithFields(log.Fields{
		"package": "preflight",
		"fn":      "preflightDriverWorker",
	})
	l.Debug("starting worker")
	for j := range jobs {
		l.WithField("driver", j.Driver).Debug("running job")
		switch j.Driver {
		case DriverNameDNS:
			i, ok := j.Job.(preflightdns.PreflightDNS)
			if !ok {
				l.WithField("job", j.Job).Error("invalid job")
				res <- errors.New("invalid job")
				continue
			}
			if err := i.Run(); err != nil {
				l.WithError(err).Error("unable to run job")
				res <- err
				continue
			}
			res <- nil
		case DriverNameEnv:
			i, ok := j.Job.(preflightenv.PreflightEnv)
			if !ok {
				l.WithField("job", j.Job).Error("invalid job")
				res <- errors.New("invalid job")
				continue
			}
			if err := i.Run(); err != nil {
				l.WithError(err).Error("unable to run job")
				res <- err
				continue
			}
			res <- nil
		case DriverNameID:
			i, ok := j.Job.(preflightid.PreflightID)
			if !ok {
				l.WithField("job", j.Job).Error("invalid job")
				res <- errors.New("invalid job")
				continue
			}
			if err := i.Run(); err != nil {
				l.WithError(err).Error("unable to run job")
				res <- err
				continue
			}
			res <- nil
		case DriverNameNetPath:
			i, ok := j.Job.(preflightnetpath.PreflightNetPath)
			if !ok {
				l.WithField("job", j.Job).Error("invalid job")
				res <- errors.New("invalid job")
				continue
			}
			if err := i.Run(); err != nil {
				l.WithError(err).Error("unable to run job")
				res <- err
				continue
			}
			res <- nil
		default:
			l.WithField("driver", j.Driver).Error("invalid driver")
			res <- errors.New("invalid driver")
			continue
		}
	}
}

func (p *Preflight) jobCount() int {
	l := log.WithFields(log.Fields{
		"package": "preflight",
		"fn":      "jobCount",
	})
	l.Debug("counting jobs")
	c := 0
	if len(p.Env) > 0 {
		c++
	}
	c += len(p.DNS)
	c += len(p.ID)
	c += len(p.Netpath)
	return c
}

func (p *Preflight) Run() error {
	l := log.WithFields(log.Fields{
		"package": "preflight",
		"fn":      "Run",
	})
	l.Debug("running preflight checks")
	if p.Concurrency == 0 {
		p.Concurrency = 1
	}
	count := p.jobCount()
	l.WithField("count", count).Debug("creating jobs")
	jobs := make(chan PreflightJob, count)
	res := make(chan error, count)
	for i := 0; i < p.Concurrency; i++ {
		go preflightDriverWorker(jobs, res)
	}
	for _, d := range p.DNS {
		jobs <- PreflightJob{
			Driver: DriverNameDNS,
			Job:    d,
		}
	}
	if len(p.Env) > 0 {
		ev := preflightenv.PreflightEnv{
			EnvVars: p.Env,
		}
		jobs <- PreflightJob{
			Driver: DriverNameEnv,
			Job:    ev,
		}
	}
	for _, d := range p.ID {
		jobs <- PreflightJob{
			Driver: DriverNameID,
			Job:    d,
		}
	}
	for _, d := range p.Netpath {
		jobs <- PreflightJob{
			Driver: DriverNameNetPath,
			Job:    d,
		}
	}
	close(jobs)
	for i := 0; i < count; i++ {
		if err := <-res; err != nil {
			l.WithError(err).Error("preflight check failed")
			return err
		}
	}
	return nil
}
