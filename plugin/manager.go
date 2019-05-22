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
	Plugins []*Plugin
	Timeout time.Duration
}

func (manager *Manager) Load() error {
	matches, err := filepath.Glob(filepath.Join(config.PLUGIN_DIR, config.PLUGIN_NAME_PATTERN))
	if err != nil {
		return err
	}
	for _, path := range matches {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			continue
		}
		plugin := new(Plugin)
		err = json.Unmarshal(bytes, plugin)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"json":  path,
				"error": err,
			}).Error("error load plugin")
		} else {
			plugin.Timeout = manager.Timeout
			logrus.WithFields(logrus.Fields{
				"json":   path,
				"plugin": *plugin,
			}).Debug("success load plugin")
			manager.Plugins = append(manager.Plugins, plugin)
		}
	}
	return nil
}

func (manager Manager) Search(query string) []*Plugin {
	plugins := make([]*Plugin, 0)
	qs := strings.Split(query, ",")
	for _, plugin := range manager.Plugins {
		for _, q := range qs {
			if q == "*" || q == "all" || strings.Contains(q, strings.ToLower(plugin.Name)) || strings.Contains(q, strings.ToLower(plugin.Desc)) {
				plugins = append(plugins, plugin)
			}
		}
	}
	return plugins
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
		jobs <- Job{&target, plugin}
	case IPPORT:
		for _, protocol := range plugin.Request.Protocols {
			for _, suffix := range plugin.Request.Suffixes {
				manager.ParseJobs(plugin, Target{URL, fmt.Sprintf("%s://%s/%s", protocol, target.Raw, strings.TrimLeft(suffix, "/"))}, jobs)
			}
		}
	case IP:
		for _, port := range plugin.Request.Ports {
			manager.ParseJobs(plugin, Target{IPPORT, fmt.Sprintf("%s:%d", target.Raw, port)}, jobs)
		}
	}

}

func (manager Manager) Execute(targets []Target, plugins []*Plugin, worker int) <-chan interface{} {
	ppool := pool.New(worker)
	ppool.Start()
	go func() {
		jobs := make(chan Job)
		go func() {
			for _, target := range targets {
				for _, plugin := range plugins {
					manager.ParseJobs(plugin, target, jobs)
				}
			}
			close(jobs)
		}()
		for job := range jobs {
			logrus.WithFields(logrus.Fields{
				"target": job.Target,
				"plugin": job.Plugin.Name,
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
