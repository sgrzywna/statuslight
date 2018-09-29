package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type statusType int

const (
	statusOK statusType = iota
	statusUnstable
	statusError

	maxStatuses = 16
)

var (
	// errTooMuchStatuses is returned when statuses queue is full.
	errTooMuchStatuses = errors.New("too much statuses")
)

// context represents status context, it stores all details necessary to calculate current status.
type context struct {
	host       string
	port       int
	stats      map[string]bool
	client     *http.Client
	colors     map[statusType]string
	sequences  map[statusType]string
	brightness int
}

// status stores status details
type status struct {
	State bool   `json:"state"`
	ID    string `json:"statusId"`
}

// processStatus process status received by http server.
func (c *context) processStatus(s status) error {
	if len(c.stats) == maxStatuses {
		return errTooMuchStatuses
	}
	c.stats[s.ID] = s.State
	sts := c.getStatus()
	sequence, _ := c.sequences[sts]
	if sequence != "" {
		return c.setSequence(sequence)
	}
	color, _ := c.colors[sts]
	return c.setLight(color)
}

// getStatus returns single status for all received statuses.
func (c *context) getStatus() statusType {
	var t, f int
	for _, s := range c.stats {
		switch s {
		case true:
			t++
		case false:
			f++
		}
	}
	res := statusOK
	if t > 0 {
		if f > 0 {
			res = statusUnstable
		}
	} else {
		if f > 0 {
			res = statusError
		}
	}
	return res
}

// setLight sets light through milightd.
func (c *context) setLight(color string) error {
	url := fmt.Sprintf("http://%s:%d/api/v1/light", c.host, c.port)

	var cmd = struct {
		Color      string `json:"color"`
		Brightness int    `json:"brightness"`
		Switch     string `json:"switch"`
	}{
		Color:      color,
		Brightness: c.brightness,
		Switch:     "on",
	}

	d, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(d))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("milightd client: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// setSequence sets sequence of lights through milightd.
func (c *context) setSequence(sequence string) error {
	url := fmt.Sprintf("http://%s:%d/api/v1/seqctrl", c.host, c.port)

	var cmd = struct {
		Name  string `json:"name"`
		State string `json:"state"`
	}{
		Name:  sequence,
		State: "running",
	}

	d, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(d))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("milightd client: unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	var mihost = flag.String("mihost", "127.0.0.1", "milightd network address")
	var miport = flag.Int("miport", 8080, "milightd network port")
	var port = flag.Int("port", 8888, "listening port")
	var okColor = flag.String("ok-color", "green", "color for the OK status")
	var unstableColor = flag.String("unstable-color", "yellow", "color for the unstable status")
	var errorColor = flag.String("error-color", "red", "color for the error status")
	var okSeq = flag.String("ok-seq", "", "sequence for the OK status")
	var unstableSeq = flag.String("unstable-seq", "", "sequence for the unstable status")
	var errorSeq = flag.String("error-seq", "", "sequence for the error status")
	var brightness = flag.Int("brightness", 32, "brightness level")

	flag.Parse()

	ctx := context{
		host:  *mihost,
		port:  *miport,
		stats: make(map[string]bool),
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		colors: map[statusType]string{
			statusOK:       *okColor,
			statusUnstable: *unstableColor,
			statusError:    *errorColor,
		},
		sequences: map[statusType]string{
			statusOK:       *okSeq,
			statusUnstable: *unstableSeq,
			statusError:    *errorSeq,
		},
		brightness: *brightness,
	}

	r := mux.NewRouter()
	v1 := r.PathPrefix("/api/v1/").Subrouter()

	v1.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		statusHandler(w, r, &ctx)
	}).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("statuslight listening @ :%d\n", *port)
	log.Fatal(srv.ListenAndServe())
}

func statusHandler(w http.ResponseWriter, r *http.Request, ctx *context) {
	var s status
	if r.Body == nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = ctx.processStatus(s)
	if err != nil {
		log.Printf("processStatus error: %s\n", err)
		http.Error(w, "statuslight error", http.StatusInternalServerError)
		return
	}
}
