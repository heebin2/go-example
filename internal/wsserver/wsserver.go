package wsserver

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type WebsocketServer struct {
	Address string
	Port    int
	hub     *hub
	server  *http.Server
	wg      sync.WaitGroup
	config  Config
}

func NewWebsocketServer() *WebsocketServer {
	return &WebsocketServer{}
}

func (ws *WebsocketServer) Open(config Config) error {
	if err := config.Invalid(); err != nil {
		return fmt.Errorf("websocket server open fail : %s", err)
	}

	// close previous connections before connecting.
	ws.Close()

	// config
	ws.Address = config.Address()
	ws.Port = config.Port
	ws.config = config

	// hub
	ws.hub = Newhub()
	if err := ws.hub.open(config.SendFilter); err != nil {
		return fmt.Errorf("websocket server hub open fail : %s", err)
	}

	// server
	mux := http.NewServeMux()
	mux.HandleFunc("/", ws.serve)
	ws.server = &http.Server{
		Addr:    ws.Address,
		Handler: mux,
	}

	ws.wg.Add(1)
	go ws.runServer()

	log.Info("listen websocket server : ws://localhost:", ws.Port)
	return nil
}

func (ws *WebsocketServer) Close() {

	// close server
	if ws.server != nil {
		ws.server.Close()
		ws.wg.Wait()
		ws.server = nil
	}

	// clouse hub
	if ws.hub != nil {
		ws.hub.close()
		ws.hub = nil
	}
}

// You can use it after registering the filter.
// Compares the filter function with the id of the client and sends them if they match.
func (ws *WebsocketServer) SendWithFilter(message []byte) error {
	if length := len(message); length <= 0 || length > maxMessageSize {
		return fmt.Errorf("websocket send with filter buffer range over %d", length)
	}

	ws.hub.send <- message

	return nil
}

// It can be sent to a client with a specific key.
func (ws *WebsocketServer) SendWithKey(message SendMessage) error {
	if err := message.Invalid(); err != nil {
		return fmt.Errorf("websocket send with key err %s", err)
	}

	ws.hub.sendKey <- message

	return nil
}

// Send to all clients.
func (ws *WebsocketServer) Broadcast(message []byte) error {
	if length := len(message); length <= 0 || length > maxMessageSize {
		return fmt.Errorf("websocket broadcast send message buffer range over %d", length)
	}

	ws.hub.broadcast <- message

	return nil
}

// Get the number of clients.
func (ws *WebsocketServer) Size() int {
	if ws.hub == nil {
		return 0
	}

	return ws.hub.size()
}

func (ws *WebsocketServer) runServer() error {
	defer ws.wg.Done()

	return ws.server.ListenAndServe()
}

// serveWs handles websocket requests from the peer.
func (ws *WebsocketServer) serve(w http.ResponseWriter, r *http.Request) {
	log.Trace("Websocket Upgrade request received : ", r.RemoteAddr)
	upgrader := websocket.Upgrader{
		ReadBufferSize:  0,
		WriteBufferSize: 0,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error(err)
		return
	}

	client := &Client{
		hub:            ws.hub,
		conn:           conn,
		send:           make(chan []byte),
		id:             make(Identity),
		openHandler:    ws.config.OpenHandler,
		messageHandler: ws.config.MessageHandler,
		closeHandler:   ws.config.CloseHandler,
	}

	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump(ws.config.ReadDeadline)
}
