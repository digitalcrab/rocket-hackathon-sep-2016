package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"github.com/tarm/serial"
	"github.com/tylerb/graceful"
	"net"
	"net/http"
	"time"
)

var (
	logLevel     = flag.String("log", "debug", "Logs level")
	listenAddr   = flag.String("listen", ":8080", "Listen on address")
	staticFolder = flag.String("staticPath", "static", "Path to static folder")
	redDevice    = flag.String("redDevice", "/dev/cu.usbmodem1411", "Red Car USB device")
	greenDevice  = flag.String("greenDevice", "/dev/cu.usbmodem1411", "Green Car USB device")

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	hub = pool{
		broadcast:   make(chan WsCommand),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool),
		cars:        make(map[string]*car),
	}
)

func serveWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			logrus.WithError(err).Errorln("Error on upgrade connection to ws")
		}
		return
	}

	logrus.WithField("ip", ws.RemoteAddr().String()).Debugln("Client connected")

	c := &connection{
		hub:  &hub,
		send: make(chan []byte, 256),
		ws:   ws,
		ip:   ws.RemoteAddr().String(),
	}

	hub.register <- c
	go c.writer()
	c.reader()
}

func main() {
	flag.Parse()

	lvl, err := logrus.ParseLevel(*logLevel)
	if err != nil {
		lvl = logrus.DebugLevel
	}

	logrus.SetLevel(lvl)
	logrus.WithFields(logrus.Fields{
		"listen": *listenAddr,
		"static": *staticFolder,
	}).Debugln("Starting application...")

	redConfig := &serial.Config{
		Name:        *redDevice,
		Baud:        9600,
		ReadTimeout: time.Second * 5,
	}

	redPort, redErr := serial.OpenPort(redConfig)
	if redErr != nil {
		logrus.WithError(redErr).Fatalln("Unable to open a USB for red car")
	}
	defer redPort.Close()

	hub.cars["red"] = &car{
		port: redPort,
	}

	go hub.run()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(*staticFolder)))
	mux.HandleFunc("/game", serveWs)

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

	if err := server.ListenAndServe(); err != nil {
		if opErr, ok := err.(*net.OpError); !ok || (ok && opErr.Op != "accept") {
			logrus.WithError(err).Fatalln("Error on listen port")
		}
	}
}
