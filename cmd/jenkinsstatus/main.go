package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/sgrzywna/statuslight/internal/app/jenkinsstatus"
	"github.com/sgrzywna/statuslight/internal/app/statuslightclient"
)

// config stores jenkinsstatus configuration.
type config struct {
	StatusLight statuslight `toml:"statuslight"`
	Jenkins     jenkins     `toml:"jenkins"`
	Jobs        []job       `toml:"job"`
}

// statuslight stores statuslight daemon configuration.
type statuslight struct {
	URL string `toml:"url"`
}

// jenkins stores Jenkins configuration.
type jenkins struct {
	URL         string `toml:"url"`
	Username    string `toml:"username"`
	Password    string `toml:"password"`
	CheckPeriod int    `toml:"check_period"`
}

// job stores Jenkins job configuration.
type job struct {
	Description string `toml:"description"`
	Path        string `toml:"path"`
}

// jenkinsStatusReceiver implements jenkinsstatus.Receiver interface.
type jenkinsStatusReceiver struct {
	client *statuslightclient.Client
}

func (r *jenkinsStatusReceiver) OnStatus(job []string, status string) {
	var sts bool

	switch status {
	case "ABORTED", "FAILURE", "NOT_BUILT", "UNSTABLE":
	case "SUCCESS":
		sts = true
	default:
		log.Printf("jenkinsStatusReceiver.OnStatus unsupported status: %s", status)
		return
	}

	err := r.client.SetStatus(strings.Join(job, "/"), sts)
	if err != nil {
		log.Printf("jenkinsStatusReceiver.OnStatus error: %s", err)
	}
}

func main() {
	var cfgPath = flag.String("config", "config.toml", "full path to the configuration file")

	flag.Parse()

	var cfg config
	if _, err := toml.DecodeFile(*cfgPath, &cfg); err != nil {
		log.Fatalf("configuration error: %s", err)
	}

	jenkins, err := jenkinsstatus.NewJenkinsClient(cfg.Jenkins.URL, cfg.Jenkins.Username, cfg.Jenkins.Password)
	if err != nil {
		log.Fatalf("jenkins error: %s", err)
	}

	statusLightClient := statuslightclient.NewClient(cfg.StatusLight.URL)

	rcv := jenkinsStatusReceiver{
		client: statusLightClient,
	}

	var jobs [][]string

	for _, job := range cfg.Jobs {
		log.Printf("Loading '%s' (%v)", job.Description, job.Path)
		path := strings.Split(job.Path, "/")
		if len(path) == 0 {
			continue
		}
		// last <-> first
		path[0], path[len(path)-1] = path[len(path)-1], path[0]
		// reverse(second, last)
		for left, right := 1, len(path)-1; left < right; left, right = left+1, right-1 {
			path[left], path[right] = path[right], path[left]
		}
		jobs = append(jobs, path)
	}

	jenkinsStatus := jenkinsstatus.NewJenkinsStatus(jenkins, jobs, time.Duration(cfg.Jenkins.CheckPeriod)*time.Second, &rcv)
	defer jenkinsStatus.Close()

	select {}
}
