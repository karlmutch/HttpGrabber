package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mgutz/logxi/v1"
)

var (
	logLevel = flag.String("loglevel", "debug", "Set the desired log level")

	url          = flag.String("url", "http://127.0.0.1:12345", "URL to query for content")
	scenarioPath = flag.String("path", "./", "Path into which the time based directories will be written.")
)

var (
	// create Logger interface
	logW = log.NewLogger(log.NewConcurrentWriter(os.Stdout), "HttpRoller")

	startedAt = time.Now()
)

func main() {

	flag.Parse()

	switch strings.ToLower(*logLevel) {
	case "debug":
		logW.SetLevel(log.LevelDebug)
	case "info":
		logW.SetLevel(log.LevelInfo)
	}

	_, err := filepath.Abs(*scenarioPath)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(-1)
	}

}

func retrieve() {
	interval := time.Duration(500 * time.Millisecond)
	alarm := time.NewTicker(interval)
	defer alarm.Stop()

	lastPayload := ""

	for {
		select {
		case <-alarm.C:
			// Retrieve the URL contento
			payload := func() (payload string) {
				resp, err := http.Get(*url)
				if err != nil {
					logW.Error(fmt.Sprintf("unable to load %s due to %s", *url, err.Error()), "error", err, "url", *url)
					return payload
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				return string(body)
			}()
			if payload != lastPayload {
				outputChange(payload)
				lastPayload = payload
			}
		}
	}
}

func outputChange(payload string) {
	logW.Info(payload)

}
