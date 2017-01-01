package main

import (
	"github.com/Sirupsen/logrus"
	"net/http"
	"net"
	"time"
	"github.com/tylerb/graceful"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugln("Starting application 127.0.0.1:8888")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("static")))

	server := &graceful.Server{
		Timeout: 10 * time.Second,
		BeforeShutdown: func() bool {
			logrus.Debugln("Stopping http server...")
			return true
		},
		Server: &http.Server{
			Addr:    ":8888",
			Handler: mux,
		},
	}

	if err := server.ListenAndServe(); err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			logrus.WithError(err).Fatalln("Error on listen port")
		}
	}
}