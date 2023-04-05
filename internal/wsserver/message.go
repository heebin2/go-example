package wsserver

import "fmt"

type SendMessage struct {
	Message []byte
	Key     string
	Value   any
}

func (sm SendMessage) Invalid() error {
	if len(sm.Message) <= 0 {
		return fmt.Errorf("message empty")
	}

	if len(sm.Key) <= 0 {
		return fmt.Errorf("key empty")
	}

	if sm.Value == nil {
		return fmt.Errorf("value empty")
	}

	return nil
}
