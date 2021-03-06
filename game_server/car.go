package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/tarm/serial"
)

const (
	DriveForward  byte = 1
	DriveBackward byte = 2
	DriveLeft     byte = 4
	DriveRight    byte = 8
)

type (
	car struct {
		port *serial.Port
	}

	cmd struct {
		Car        string   `json:"car"`
		Speed      byte     `json:"speed"`
		Directions []string `json:"directions"`
	}
)

func (c *car) send(cmd cmd) {
	var direction byte

	for i := range cmd.Directions {
		switch cmd.Directions[i] {
		case "stop":
			direction = 0
		case "up":
			direction += DriveForward
		case "down":
			direction += DriveBackward
		case "left":
			direction += DriveLeft
		case "right":
			direction += DriveRight
		}
	}

	msg := []byte{direction, cmd.Speed}

	logrus.WithField("msg", msg).Infoln("Sending to arduino")

	if _, err := c.port.Write(msg); err != nil {
		logrus.WithError(err).Errorln("Error on send message")
	}
}
