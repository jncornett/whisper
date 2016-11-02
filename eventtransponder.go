package whisper

type Handler func(Message, []Peer)

type EventTransponder interface {
	Push(Message)
	Serve(Handler)
}

type eventTransponder struct {
	peers    PeerSet
	messages chan Message
}

// Push submits a message to the broadcast queue. This method is blocking
// unless the queue is buffere
func (q *eventTransponder) Push(m Message) {
	q.messages <- m
}

func (q *eventTransponder) Serve(cb Handler) {
	for {
		m, more := <-q.messages
		if !more {
			break
		}
		cb(m, q.peers.GetAll())
	}
}

func NewEventTransponder(peers PeerSet, bufsize int) EventTransponder {
	return &eventTransponder{peers: peers, messages: make(chan Message, bufsize)}
}
