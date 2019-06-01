package plugin

import (
	"biu/config"
	"biu/pool"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"path/filepath"
	"strings"
	"time"

	gocidr "github.com/apparentlymart/go-cidr/cidr"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	POCs    []*POC
	Timeout time.Duration
}

func (manager *Manager) Load() error {
	logrus.WithFields(logrus.Fields{
		"target": config.POC_DIR,
	}).Debug("load poc")
	matches, err := filepath.Glob(filepath.Join(config.POC_DIR, config.POC_NAME_PATTERN))
	if err != nil {
		return err
	}
	for _, path := range matches {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}
		poc := new(POC)
		err = json.Unmarshal(bytes, poc)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"json":  path,
				"error": err,
			}).Error("error load poc")
		} else {
			poc.Timeout = manager.Timeout
			logrus.WithFields(logrus.Fields{
				"json": path,
				"poc":  *poc,
			}).Debug("success load poc")
			manager.POCs = append(manager.POCs, poc)
		}
	}
	return nil
}

func (manager Manager) Search(query string) []*POC {
	pocs := make([]*POC, 0)
	qs := strings.Split(query, ",")
	for _, poc := range manager.POCs {
		for _, q := range qs {
			if q == "*" || q == "all" || strings.Contains(q, strings.ToLower(poc.Name)) || strings.Contains(q, strings.ToLower(poc.Desc)) {
				pocs = append(pocs, poc)
			}
		}
	}
	return pocs
}

func (manager Manager) ParseTargets(target string, cidr string) []Target {
	targets := make([]Target, 0)
	ts := strings.Split(target, ",")
	for _, t := range ts {
		if strings.HasPrefix(t, "http://") || strings.HasPrefix(t, "https://") {
			targets = append(targets, Target{URL, t})
		} else {
			elements := strings.Split(t, ":")
			if len(elements) > 1 {
				targets = append(targets, Target{IPPORT, t})
			} else {
				targets = append(targets, Target{IP, t})
			}
		}
	}
	cs := strings.Split(cidr, ",")
	for _, c := range cs {
		_, ipnet, err := net.ParseCIDR(c)
		if err == nil {
			count := gocidr.AddressCount(ipnet)
			ip, _ := gocidr.AddressRange(ipnet)
			for count > 1 {
				targets = append(targets, Target{IP, ip.String()})
				ip = gocidr.Inc(ip)
				count--
			}
		}
	}
	return targets
}

type Job struct {
	Target *Target
	Plugin *Plugin
}

func (manager *Manager) ParseJobs(plugin *Plugin, target Target, jobs chan<- Job) {
	switch target.Type {
	case URL:
		for _, suffix := range plugin.POC.Request.Suffixes {
			jobs <- Job{&Target{URL, fmt.Sprintf("%s/%s", strings.TrimRight(target.Raw, "/"), strings.TrimLeft(suffix, "/"))}, plugin}
		}
	case IPPORT:
		for _, protocol := range plugin.POC.Request.Protocols {
			manager.ParseJobs(plugin, Target{URL, fmt.Sprintf("%s://%s", protocol, target.Raw)}, jobs)
		}
	case IP:
		for _, port := range plugin.POC.Request.Ports {
			manager.ParseJobs(plugin, Target{IPPORT, fmt.Sprintf("%s:%d", target.Raw, port)}, jobs)
		}
	}

}

func (manager Manager) Execute(targets []Target, pocs []*POC, worker int) <-chan interface{} {
	ppool := pool.New(worker)
	ppool.Start()
	jobs := make(chan Job)
	go func() {
		for _, poc := range pocs {
			plugins := poc.BuildPlugins()
			for _, target := range targets {
				for _, plugin := range plugins {
					manager.ParseJobs(plugin, target, jobs)
				}
			}
		}
		close(jobs)
	}()
	go func() {
		for job := range jobs {
			logrus.WithFields(logrus.Fields{
				"target": job.Target,
				"plugin": job.Plugin.POC.Name,
			}).Debug("add task to pool")
			ppool.Add(func(job Job) func() interface{} {
				return func() interface{} {
					return job.Plugin.Execute(*job.Target)
				}
			}(job))
		}
		ppool.CloseAndWait()
	}()

	return ppool.Results
}
