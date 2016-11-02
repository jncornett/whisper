package whisper

type Handler func(Message, []Peer)

type EventQueue interface {
	Push(Message)
	Serve(Handler)
}

type eventQueue struct {
	peers    PeerSet
	messages chan Message
}

// Push submits a message to the broadcast queue. This method is blocking
// unless the queue is buffere
func (q *eventQueue) Push(m Message) {
	q.messages <- m
}

func (q *eventQueue) Serve(cb Handler) {
	for {
		m, more := <-q.messages
		if !more {
			break
		}
		cb(m, q.peers.GetAll())
	}
}

func NewEventQueue(bufsize int) EventQueue {
	return new(eventQueue)
}
