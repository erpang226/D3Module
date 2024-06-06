package msghandler

import (
	"fmt"
	"main/pkg/actor/core/model"
	northClients "main/pkg/north/adapter"
	"sync"
)

// Handler handler different messages
type Handler interface {
	Filter(message *model.Message) bool
	Process(message *model.Message, northAdapter northClients.Adapter) error
}

var Handlers []Handler
var lock sync.RWMutex

func init() {
	Handlers = make([]Handler, 0)
}

// RegisterHandler register msg handler.
func RegisterHandler(handler Handler) {
	lock.Lock()
	defer lock.Unlock()
	Handlers = append(Handlers, handler)
}

// ProcessHandler return true if handler filtered
func ProcessHandler(message model.Message, client northClients.Adapter) error {
	lock.RLock()
	defer lock.RUnlock()
	for _, handle := range Handlers {
		if handle.Filter(&message) {
			err := handle.Process(&message, client)
			if err != nil {
				return fmt.Errorf("failed to handle message, message group: %s, error: %+v", message.GetGroup(), err)
			}
			return nil
		}
	}
	return fmt.Errorf("failed to handle message, no handler found for the message, message group: %s", message.GetGroup())
}
