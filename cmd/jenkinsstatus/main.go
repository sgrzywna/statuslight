package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/sgrzywna/statuslight/internal/app/jenkinsstatus"
	"github.com/sgrzywna/statuslight/internal/app/statuslightclient"
)

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
	var jenkinsURL = flag.String("jenkinsurl", "http://admin:admin@127.0.0.1:8080", "Jenkins URL")
	var statusURL = flag.String("statusurl", "http://127.0.0.1:8888", "status light daemon URL")
	var checkPeriod = flag.Int("checkperiod", 180, "check period in seconds")

	flag.Parse()

	url, err := url.Parse(*jenkinsURL)
	if err != nil {
		log.Fatal(err)
	}

	strippedURL := fmt.Sprintf("%s://%s", url.Scheme, url.Host)
	username := url.User.Username()
	password, _ := url.User.Password()

	jenkins, err := jenkinsstatus.NewJenkinsClient(strippedURL, username, password)
	if err != nil {
		log.Fatal(err)
	}

	statusLightClient := statuslightclient.NewClient(*statusURL)

	rcv := jenkinsStatusReceiver{
		client: statusLightClient,
	}

	jobs := [][]string{
		[]string{"webapp-develop-config-tests", "webapp-develop"},
	}

	jenkinsStatus := jenkinsstatus.NewJenkinsStatus(jenkins, jobs, time.Duration(*checkPeriod)*time.Second, &rcv)
	defer jenkinsStatus.Close()

	select {}
}
