package jenkinsstatus

import (
	"log"
	"time"
)

// Receiver represents receiver of Jenkins job statuses.
type Receiver interface {
	OnStatus(job []string, status string)
}

// JenkinsStatus represents Jenkins high level client.
type JenkinsStatus struct {
	jenkins     *JenkinsClient
	jobs        [][]string
	checkPeriod time.Duration
	rcv         Receiver
	close       chan bool
}

// NewJenkinsStatus returns initialized JenkinsStatus object.
func NewJenkinsStatus(jenkins *JenkinsClient, jobs [][]string, checkPeriod time.Duration, rcv Receiver) *JenkinsStatus {
	jenkinsStatus := JenkinsStatus{
		jenkins:     jenkins,
		jobs:        jobs,
		checkPeriod: checkPeriod,
		rcv:         rcv,
		close:       make(chan bool),
	}
	go jenkinsStatus.probeLoop()
	return &jenkinsStatus
}

// Close stops probing loop.
func (s *JenkinsStatus) Close() {
	s.close <- true
}

// probeLoop is the main processing loop.
func (s *JenkinsStatus) probeLoop() {
	for {
		select {
		case <-s.close:
			return
		case <-time.After(s.checkPeriod):
			s.checkStatus()
		}
	}
}

// checkStatus probes jenkins for job status
func (s *JenkinsStatus) checkStatus() {
	for _, job := range s.jobs {
		sts, err := s.jenkins.GetStatus(job[0], job[1:]...)
		if err != nil {
			log.Printf("jenkins.GetStatus error: %s", err)
		} else {
			s.rcv.OnStatus(job, sts)
		}
	}
}
