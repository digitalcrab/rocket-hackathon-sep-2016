package main

import "github.com/Sirupsen/logrus"

type (
	pool struct {
		connections map[*connection]bool
		commands    chan cmd
		register    chan *connection
		unregister  chan *connection
		cars        map[string]*car
	}
)

func (p *pool) run() {
	for {
		select {
		case c := <-p.register:
			p.connections[c] = true
			logrus.WithField("ip", c.ip).Debugln("Client registered")

		case c := <-p.unregister:
			if _, ok := p.connections[c]; ok {
				delete(p.connections, c)
				close(c.send)
				logrus.WithField("ip", c.ip).Debugln("Client removed")
			}

		case m := <-p.commands:
			logrus.WithFields(logrus.Fields{
				"car":        m.Car,
				"speed":      m.Speed,
				"directions": m.Directions,
			}).Debugln("Command received")

			if car, ok := p.cars[m.Car]; ok {
				car.send(m)
			}
		}
	}
}
