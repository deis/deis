package event

type Event struct {
	Type    string
	Payload interface{}
	Context interface{}
}
