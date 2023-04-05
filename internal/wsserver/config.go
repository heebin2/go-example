package wsserver

import (
	"fmt"
	"strconv"
	"time"
)

type Identity = map[string]any
type OpenHandler = func(Identity) error
type MessageHandler = func(Identity, []byte) []byte
type CloseHandler = func(Identity)
type SendFilter = func(Identity) bool

type Config struct {
	Port           int
	Singletone     bool
	OpenHandler    OpenHandler
	MessageHandler MessageHandler
	CloseHandler   CloseHandler
	SendFilter     SendFilter
	ReadDeadline   time.Duration
}

func (c Config) Invalid() error {
	if c.Port < 1000 || c.Port > 65535 {
		return fmt.Errorf("invalid port")
	}

	return nil
}

func (c Config) Address() string {
	return ":" + strconv.Itoa(c.Port)
}
