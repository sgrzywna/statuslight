package statuslight

import (
	"errors"
	"log"
	"time"

	"github.com/sgrzywna/milightd/pkg/milightdclient"
	"github.com/sgrzywna/milightd/pkg/models"
)

// statusType represents type of status.
type statusType int

// StatusMap stores mapping between status and related action.
type StatusMap map[statusType]string

const (
	// StatusOK represents OK status.
	StatusOK statusType = iota
	// StatusUnstable represents unstable status.
	StatusUnstable
	// StatusError represents error status.
	StatusError
	// maxStatuses defines maximal number of different statuses that can be processed by statuslight daemon.
	maxStatuses = 16
	// setStatusPeriod defines how often statuslight daemon will connect to milightd daemon to set the light.
	setStatusPeriod = 30 * time.Second
)

var (
	// errTooMuchStatuses is returned when statuses queue is full.
	errTooMuchStatuses = errors.New("too much statuses")
)

// Status stores status details
type Status struct {
	State bool   `json:"state"`
	ID    string `json:"statusId"`
}

// StatusLight represents status context, it stores all details necessary to calculate current status.
type StatusLight struct {
	stats      map[string]bool
	client     *milightdclient.Client
	colors     StatusMap
	sequences  StatusMap
	brightness int
	quit       chan struct{}
}

// NewStatusLight returns initialized StatusLight object.
func NewStatusLight(miURL string, colors, sequences StatusMap, brightness int) *StatusLight {
	statusLight := StatusLight{
		stats:      make(map[string]bool),
		client:     milightdclient.NewClient(miURL),
		colors:     colors,
		sequences:  sequences,
		brightness: brightness,
		quit:       make(chan struct{}),
	}
	go statusLight.statusLoop()
	return &statusLight
}

// Close terminates status loop.
func (c *StatusLight) Close() {
	// Blocks until received.
	c.quit <- struct{}{}
}

// processStatus process status received by http server.
func (c *StatusLight) processStatus(s Status) error {
	if len(c.stats) == maxStatuses {
		return errTooMuchStatuses
	}
	c.stats[s.ID] = s.State
	return nil
}

// statusLoop is the main processing loop.
func (c *StatusLight) statusLoop() {
	// set status immediately
	lastStatus := c.getStatus()
	err := c.setStatus(lastStatus)
	if err != nil {
		log.Printf("statuslight.setStatus error: %s", err)
	}

	for {
		select {
		case <-c.quit:
			return
		case <-time.After(setStatusPeriod):
			sts := c.getStatus()
			if lastStatus != sts {
				err = c.setStatus(sts)
				if err != nil {
					log.Printf("statuslight.setStatus error: %s", err)
				} else {
					lastStatus = sts
				}
			}
		}
	}
}

// getStatus returns single status for all received statuses.
func (c *StatusLight) getStatus() statusType {
	var t, f int
	for _, s := range c.stats {
		switch s {
		case true:
			t++
		case false:
			f++
		}
	}
	res := StatusOK
	if t > 0 {
		if f > 0 {
			res = StatusUnstable
		}
	} else {
		if f > 0 {
			res = StatusError
		}
	}
	return res
}

// setStatus send command to milightd daemon to set light according to provided status.
func (c *StatusLight) setStatus(sts statusType) error {
	sequence, _ := c.sequences[sts]
	if sequence != "" {
		return c.setSequence(sequence)
	}
	color, _ := c.colors[sts]
	return c.setLight(color)
}

// setLight sets light through milightd.
func (c *StatusLight) setLight(color string) error {
	var light models.Light

	light.SetColor(color)
	light.SetBrightness(c.brightness)
	light.SetSwitch(true)

	return c.client.SetLight(light)
}

// setSequence sets sequence of lights through milightd.
func (c *StatusLight) setSequence(sequence string) error {
	state := models.SequenceState{
		Name:  sequence,
		State: models.SeqRunning,
	}
	return c.client.SetSequenceState(state)
}
