package main

import (
	"bufio"
	"bytes"
	"github.com/tylerb/graceful"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var stream = NewStream()

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", stream)

	server := &graceful.Server{
		Timeout: 10 * time.Second,
		BeforeShutdown: func() bool {
			log.Println("Stopping http server...")
			return true
		},
		Server: &http.Server{
			Addr:    ":8081",
			Handler: mux,
		},
	}

	go runFF()

	log.Println("Starting application 127.0.0.1:8081")

	if err := server.ListenAndServe(); err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			log.Fatalln(err)
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
		log.Fatalln(err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(scanImages)

	go func() {
		for scanner.Scan() {
			buf := scanner.Bytes()
			log.Printf("Received jpeg with length %d\n", len(buf))
			stream.UpdateJPEG(buf)
		}
	}()

	if err := cmd.Start(); err != nil {
		log.Fatalln(err)
	}

	err = cmd.Wait()
	if err != nil {
		os.Exit(1)
	}
}

func scanImages(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{255, 217}); i >= 0 {
		return i + 2, data[0 : i+2], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}
