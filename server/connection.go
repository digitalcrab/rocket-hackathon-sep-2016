package main

import (
	"bytes"
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	"time"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type (
	connection struct {
		hub  *pool
		ws   *websocket.Conn
		send chan []byte
		ip   string
	}
)

func (c *connection) reader() {
	defer func() {
		c.hub.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(512)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logrus.WithError(err).Errorln("Error on read")
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))
		logrus.WithField("msg", string(message)).Infoln("Message reseived")

		var wsMsg WsCommand
		err = json.Unmarshal(message, &wsMsg)
		if err != nil {
			logrus.WithError(err).Errorln("Unable to Unmarshal message")
		} else {
			c.hub.broadcast <- wsMsg
		}
	}
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return c.ws.WriteMessage(mt, payload)
}

func (c *connection) writer() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	logrus.WithField("ip", c.ip).Debugln("Client writer started")

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				logrus.WithField("ip", c.ip).Errorln("Error on read from output channel")
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				logrus.WithFields(logrus.Fields{
					"err": err,
					"ip":  c.ip,
					"msg": string(message),
				}).Errorln("Error on write to web socket")
				return
			}

		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				logrus.WithFields(logrus.Fields{
					"ip":  c.ip,
					"err": err,
				}).Errorln("Error on ping client")
				return
			}
		}
	}
}
