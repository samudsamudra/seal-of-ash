package events

type Event struct {
	Type      string
	Entity    string
	EntityID uint
	ActorID  uint
	IP       string
	UA       string
}

var EventBus = make(chan Event, 100)
