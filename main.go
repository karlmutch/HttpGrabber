package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mgutz/logxi/v1"
)

var (
	logLevel = flag.String("loglevel", "debug", "Set the desired log level")

	url          = flag.String("url", "http://127.0.0.1:12345/module/status/json", "URL to query for content")
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

	files, err := ioutil.ReadDir(*scenarioPath)
	if len(files) > 0 {
		log.Fatal(fmt.Sprintf("%s already contains files aborting", *scenarioPath))
		os.Exit(-1)
	}
	if _, err := os.Stat(*scenarioPath); !os.IsNotExist(err) {
		if os.MkdirAll(*scenarioPath, 0777); err != nil {
			log.Fatal(fmt.Sprintf("%s could not be created due to %s", *scenarioPath, err.Error()))
			os.Exit(-1)
		}
	}

	retrieve()
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

func writeFile(fileName string, payload string) (err error) {
	f, err := os.Create(fileName)
	if err != nil {
		return nil
	}
	defer f.Close()

	_, err = io.Copy(f, strings.NewReader(payload))
	return err
}

func outputChange(payload string) {
	seconds := int(math.Floor(time.Since(startedAt).Seconds()))
	since := time.Duration(seconds) * time.Second

	logW.Info(fmt.Sprintf("%s %s", since.String(), payload))

	fp := filepath.Join(*scenarioPath, strconv.Itoa(seconds), "module", "status")
	if err := os.MkdirAll(fp, 0777); err != nil {
		logW.Error(fmt.Sprintf("%s could not be created due to %s", fp, err.Error()))
		return
	}

	fp = filepath.Join(fp, "json")
	if err := writeFile(fp, payload); err != nil {
		logW.Error(fmt.Sprintf("%s could not be created due to %s", fp, err.Error()))
		return
	}
}
