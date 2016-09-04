package main

import (
	"bufio"
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/tylerb/graceful"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
	"bytes"
)

var (
	logLevel   = flag.String("log", "debug", "Logs level")
	listenAddr = flag.String("listen", ":8082", "Listen on address")
	stream     = NewStream()
)

func main() {
	flag.Parse()

	lvl, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		lvl = logrus.DebugLevel
	}

	logrus.SetLevel(lvl)
	logrus.WithFields(logrus.Fields{
		"listen": *listenAddr,
	}).Debugln("Starting application...")

	mux := http.NewServeMux()
	mux.Handle("/", stream)

	server := &graceful.Server{
		Timeout: 10 * time.Second,
		BeforeShutdown: func() bool {
			logrus.Debugln("Stopping http server...")
			return true
		},
		Server: &http.Server{
			Addr:    *listenAddr,
			Handler: mux,
		},
	}

	go runFF()

	if err := server.ListenAndServe(); err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			logrus.WithError(err).Fatalln("Error on listen port")
		}
	}
}

func runFF() {
	cmd := exec.Command(
		"ffmpeg",
		"-f",
		"avfoundation",
		"-video_size",
		"640x480",
		"-framerate",
		"30",
		"-i",
		"1",
		"-an",
		"-f",
		"mjpeg",
		"-",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Fatalln(err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(scanImages)

	go func() {
		for scanner.Scan() {
			buf := scanner.Bytes()
			logrus.WithField("length", len(buf)).Infoln("Received jpeg")
			stream.UpdateJPEG(buf)
		}
	}()

	if err := cmd.Start(); err != nil {
		logrus.Fatalln(err)
	}

	err = cmd.Wait()
	if err != nil {
		os.Exit(1)
	}
}

func scanImages(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{255, 217}); i >= 0 {
		return i + 2, data[0:i + 2], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}