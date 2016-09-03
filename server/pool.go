package main

import "github.com/Sirupsen/logrus"

type (
	pool struct {
		connections map[*connection]bool
		broadcast   chan WsCommand
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

		case m := <-p.broadcast:
			logrus.WithFields(logrus.Fields{
				"car":       m.Car,
				"speed":     m.Speed,
				"direction": m.Direction,
			}).Debugln("Broadcast message received")

			if car, ok := p.cars[m.Car]; ok {
				car.send(m)
			}

			/*
				for c := range p.connections {
					select {
					case c.send <- m:
					default:
						logrus.WithFields(logrus.Fields{
							"ip":  c.ip,
							"msg": string(m),
						}).Errorln("Error on write to web socket")

						close(c.send)
						delete(p.connections, c)
					}
				}
			*/
		}
	}
}
