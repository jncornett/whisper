package main

import (
	"flag"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/jncornett/whisper"
)

type Engine struct {
	MaxPeerAge            time.Duration
	PeerFlushFreq         time.Duration
	PeerTableTransmitFreq time.Duration
	peers                 whisper.PeerSet
	transponder           whisper.EventTransponder
}

func (e Engine) ServeTransmitPeerTables() {
	log.Print("starting TransmitPeerTables service")
	for {
		time.Sleep(e.PeerTableTransmitFreq)
		peers := e.peers.GetAll()
		if len(peers) == 0 {
			log.Print("not transmitting empty peer table")
			continue
		} else {
			log.Print("transmitting peer table")
		}
		recip := peers[rand.Int()%len(peers)]
		client, err := rpc.DialHTTP("tcp", recip.Address)
		if err != nil {
			log.Print(err)
			continue
		}
		err = client.Call("Service.UpdatePeers", &peers, nil)
		if err != nil {
			log.Print(err) // FIXME need better logging
			continue
		}
		// Success! Now we should update the peer table to reflect that
		e.peers.Add(recip.Address)
	}
}

func (e Engine) ServeExpirePeers() {
	log.Print("starting ExpirePeers service")
	for {
		time.Sleep(e.PeerFlushFreq)
		log.Print("expiring old peers")
		e.peers.Expire(e.MaxPeerAge)
		if e.peers.Empty() {
			log.Print("warning: no peers in peer table")
		}
	}
}

func (e Engine) ServeTransponder() {
	log.Print("starting Transponder service")
	e.transponder.Serve(func(m whisper.Message, peers []whisper.Peer) {
		if m.TTL <= 0 {
			return
		}
		m.TTL--
		log.Printf("transmitting message to peers: %+v", m)
		for _, recip := range peers {
			go func() {
				client, err := rpc.DialHTTP("tcp", recip.Address)
				if err != nil {
					log.Print(err)
					return
				}
				err = client.Call("Service.Push", &m, nil)
				if err != nil {
					log.Print(err)
				}
				// Success! now we should update the peer table to reflect that
				e.peers.Add(recip.Address)
			}()
		}
	})
}

type Service struct {
	wrapped *Engine
}

func (s Service) Push(m *whisper.Message, reply *bool) error {
	log.Printf("service: responding to Push(%v)", m)
	s.wrapped.transponder.Push(*m)
	// FIXME update the peer table with the sender of this push
	// which means the sender will have to send a 'from' field
	*reply = true
	return nil
}

func (s Service) UpdatePeers(peers *[]whisper.Peer, reply *bool) error {
	log.Print("service: responding to UpdatePeers(%v)", peers)
	for _, peer := range *peers {
		s.wrapped.peers.Add(peer)
	}
	*reply = true
	return nil
}

const (
	// FIXME there is a good debugging value, and a good production value...
	// A good debugging value for this might be 5 min,
	// while a good production value for this might be 1 week!
	defaultMaxPeerAge            = time.Minute
	defaultPeerFlushFreq         = 10 * time.Second
	defaultBufferSize            = 10
	defaultPeerTableTransmitFreq = 5 * time.Second
	defaultListenAddress         = ":8081"
)

func main() {
	maxPeerAge := flag.Duration("maxpeerage", defaultMaxPeerAge, "max age of peers in seconds")
	peerFlushFreq := flag.Duration("peerflushfreq", defaultPeerFlushFreq, "frequency of cache flushing in seconds")
	listenAddress := flag.String("listen", defaultListenAddress, ":port to listen on")
	flag.Parse()
	peers := whisper.SyncedPeerSet{}
	engine := Engine{
		MaxPeerAge:            *maxPeerAge,
		PeerFlushFreq:         *peerFlushFreq,
		PeerTableTransmitFreq: defaultPeerTableTransmitFreq,
		peers:       &peers,
		transponder: whisper.NewEventTransponder(&peers, defaultBufferSize),
	}
	svc := Service{&engine}
	rpc.HandleHTTP()
	rpc.Register(&svc)
	ln, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		log.Fatal(err)
	}

	go engine.ServeExpirePeers()
	go engine.ServeTransmitPeerTables()
	go engine.ServeTransponder()

	http.Serve(ln, nil)
}
