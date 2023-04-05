package wsserver

import (
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
)

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	send       chan []byte
	sendKey    chan SendMessage
	register   chan *Client
	unregister chan *Client

	sendFilter SendFilter

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func Newhub() *hub {
	ctx, cancel := context.WithCancel(context.TODO())
	return &hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		send:       make(chan []byte),
		sendKey:    make(chan SendMessage),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		ctx:        ctx,
		cancel:     cancel,
	}
}

func (h *hub) open(sendFilter SendFilter) error {

	h.sendFilter = sendFilter

	h.wg.Add(1)
	go h.run()

	return nil
}

func (h *hub) close() {
	close(h.broadcast)
	close(h.send)
	close(h.sendKey)
	close(h.register)
	close(h.unregister)

	h.cancel()
	h.wg.Wait()
}

func (h *hub) size() int {
	return len(h.clients)
}

func (h *hub) run() {
	defer h.wg.Done()
	for {
		select {
		// register
		case client := <-h.register:
			h.clients[client] = true
			log.Debug("Websocket Server Registed Client : ", client.conn.RemoteAddr(), " client len=", len(h.clients))

		// unregister
		case client := <-h.unregister:
			length := len(h.clients)
			adr := client.conn.RemoteAddr()
			delete(h.clients, client)
			close(client.send)
			log.Debug("Websocket Server Unregisted Client : ", adr, ", client len=", length)

		// broadcast message
		case message := <-h.broadcast:
			for client := range h.clients {
				client.send <- message
			}

		// send message
		case message := <-h.send:
			for client := range h.clients {
				if h.sendFilter == nil {
					log.Warn("sendfilter is null")
					client.send <- message
				} else {
					if h.sendFilter(client.id) {
						client.send <- message
					}
				}
			}

		// send message
		case message := <-h.sendKey:
			for client := range h.clients {
				if v, exist := client.id[message.Key]; exist && v == message.Value {
					client.send <- message.Message
				}
			}

		// goroutine close
		case <-h.ctx.Done():
			for client := range h.clients {
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}
