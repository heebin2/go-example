package main

import (
	"go-helper/internal/helper"
	"go-helper/internal/wsserver"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

// type Identity = map[string]any
// type OpenHandler = func(Identity) error
// type MessageHandler = func(Identity, []byte) []byte
// type CloseHandler = func(Identity)
// type SendFilter = func(Identity) bool

var pool int

func Open(id wsserver.Identity) error {
	id["id"] = strconv.Itoa(pool)
	pool++
	logrus.Info("connect ", pool)
	return nil
}

func Message(id wsserver.Identity, message []byte) []byte {
	return []byte("id : " + id["id"].(string))
}

func Close(id wsserver.Identity) {
	logrus.Info("disconnect ", id["id"])
}

func SendFilter(id wsserver.Identity) bool {
	i, err := strconv.Atoi(id["id"].(string))
	if err != nil {
		logrus.Error(err)
		return false
	}

	logrus.Info("filter : ", i, " ", id["id"])
	return i%2 == 0
}

func main() {

	helper.InitLogger("wsserver", "trace", true)

	ws := wsserver.NewWebsocketServer()
	if err := ws.Open(wsserver.Config{
		Port:           12020,
		OpenHandler:    Open,
		MessageHandler: Message,
		CloseHandler:   Close,
		SendFilter:     SendFilter,
		// ReadDeadline:   3 * time.Second,
	}); err != nil {
		logrus.Error(err)
		return
	}

	for {
		// ws.SendWithKey(wsserver.SendMessage{
		// 	Message: []byte("sendmessage"),
		// 	Key:     "id",
		// 	Value:   "3",
		// })
		ws.Broadcast([]byte(time.Now().String()))
		time.Sleep(1 * time.Second)
	}
}
