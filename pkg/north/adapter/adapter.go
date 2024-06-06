package adapter

// Adapter is a north client interface
type Adapter interface {
	Init() error
	UnInit()
	SendToNorth(payload interface{}, params ...string) error
	ReceiveFromNorth(params ...interface{}) error
}
