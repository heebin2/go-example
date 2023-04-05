package wsserver

import (
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	// If the write is not processed within the writeWait
	// time after the write attempt, the connection is disconnected.
	// Time allowed to write a message to the peer.
	writeWait = 3 * time.Second

	// If there is no pong for the pongWait time, the connection will be disconnected.
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 100 * 1024 * 1024
)

var (
	newline = []byte{'\n'}
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Information by client
	id Identity

	// receive handler
	messageHandler MessageHandler
	closeHandler   CloseHandler
	openHandler    OpenHandler
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(readDeadline time.Duration) {
	defer func() {
		if c.closeHandler != nil {
			c.closeHandler(c.id)
		}
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize) //1GB
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	if c.openHandler != nil {
		if err := c.openHandler(c.id); err != nil {
			return
		}
	}

	for {
		if readDeadline > 0 {
			c.conn.SetReadDeadline(time.Now().Add(readDeadline))
		}
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Warn("client error, closed connection ", c.conn.RemoteAddr(), " : ", err)
			break
		}

		if response := c.messageHandler(c.id, message); response != nil {
			c.send <- response
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}