package main

import (
	"flag"
	"log"

	"github.com/sgrzywna/statuslight/internal/app/statuslight"
)

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

	statusLight := statuslight.NewStatusLight(
		*mihost,
		*miport,
		statuslight.StatusMap{
			statuslight.StatusOK:       *okColor,
			statuslight.StatusUnstable: *unstableColor,
			statuslight.StatusError:    *errorColor,
		},
		statuslight.StatusMap{
			statuslight.StatusOK:       *okSeq,
			statuslight.StatusUnstable: *unstableSeq,
			statuslight.StatusError:    *errorSeq,
		},
		*brightness,
	)

	srv := statuslight.NewHTTPServer(*port, statusLight)

	log.Printf("statuslight listening @ :%d\n", *port)
	log.Fatal(srv.ListenAndServe())
}
