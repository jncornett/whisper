package whisper

import (
	"sync"
	"time"
)

type Peer struct {
	Address     string
	LastUpdated time.Time
}

type PeerSet struct {
	peers map[string]Peer
	mux   sync.RWMutex
}

func (ps *PeerSet) Add(peer Peer) {
	ps.mux.Lock()
	defer ps.mux.Unlock()
	if ps.peers == nil {
		ps.peers = make(map[string]Peer)
	}
	ps.peers[peer.Address] = peer
}

func (ps *PeerSet) Remove(address string) {
	ps.mux.Lock()
	defer ps.mux.Unlock()
	delete(ps.peers, address)
}

func (ps *PeerSet) Lookup(address string) Peer {
	ps.mux.RLock()
	defer ps.mux.RUnlock()
	peer, _ := ps.peers[address]
	return peer
}

func (ps *PeerSet) GetAll() []Peer {
	// get the write lock so that the map doesn't change size
	ps.mux.Lock()
	defer ps.mux.Unlock()
	list := make([]Peer, len(ps.peers))
	for _, peer := range ps.peers {
		list = append(list, peer)
	}
	return list
}
