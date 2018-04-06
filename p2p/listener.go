package p2p

import (
	p2phost "gx/ipfs/QmNmJZL7FQySMtE2BQuLMuZg2EB2CLEunJJUSVSc9YnnbV/go-libp2p-host"
	pstore "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	peer "gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
)

type Listener interface {
	Protocol() string
	Address() string

	// Close closes the listener. Does not affect child streams
	Close() error
}

// NewP2P creates new P2P struct
func NewP2P(identity peer.ID, peerHost p2phost.Host, peerstore pstore.Peerstore) *P2P {
	return &P2P{
		identity:  identity,
		peerHost:  peerHost,
		peerstore: peerstore,
	}
}

// ListenerRegistry is a collection of local application proto listeners.
type ListenerRegistry struct {
	Listeners map[string]Listener
}

// Register registers listenerInfo2 in this registry
func (c *ListenerRegistry) Register(listenerInfo Listener) {
	c.Listeners[listenerInfo.Protocol()] = listenerInfo
}

// Deregister removes p2p listener from this registry
func (c *ListenerRegistry) Deregister(proto string) {
	delete(c.Listeners, proto)
}
