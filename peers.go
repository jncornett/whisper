package whisper

import (
	"sync"
	"time"
)

type Peer struct {
	Address     string
	LastUpdated time.Time
}

type PeerSet interface {
	Add(Peer)
	Remove(string)
	Lookup(string) Peer
	GetAll() []Peer
	Expire(time.Duration)
}

type SyncedPeerSet struct {
	peers map[string]Peer
	mux   sync.RWMutex
}

func (ps *SyncedPeerSet) Add(peer Peer) {
	if ps == nil {
		return
	}
	ps.mux.Lock()
	defer ps.mux.Unlock()
	if ps.peers == nil {
		ps.peers = make(map[string]Peer)
	}
	ps.peers[peer.Address] = peer
}

func (ps *SyncedPeerSet) Remove(address string) {
	if ps == nil {
		return
	}
	ps.mux.Lock()
	defer ps.mux.Unlock()
	delete(ps.peers, address)
}

func (ps *SyncedPeerSet) Lookup(address string) Peer {
	if ps == nil {
		return Peer{}
	}
	ps.mux.RLock()
	defer ps.mux.RUnlock()
	peer, _ := ps.peers[address]
	return peer
}

func (ps *SyncedPeerSet) GetAll() []Peer {
	if ps == nil {
		return []Peer{}
	}
	// get the write lock so that the map doesn't change size
	ps.mux.Lock()
	defer ps.mux.Unlock()
	list := make([]Peer, len(ps.peers))
	for _, peer := range ps.peers {
		list = append(list, peer)
	}
	return list
}

func (ps *SyncedPeerSet) Expire(maxAge time.Duration) {
	if ps == nil {
		return
	}
	for _, peer := range ps.GetAll() {
		if time.Since(peer.LastUpdated) > maxAge {
			ps.Remove(peer.Address)
		}
	}
}
