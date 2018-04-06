package p2p

import (
	peer "gx/ipfs/QmcJukH2sAFjY3HdBKq35WDzWoL3UUu2gt9wdfqZTUyM74/go-libp2p-peer"
	pstore "gx/ipfs/QmdeiKhUy1TVGBaKxt7y1QmBDLBdisSrLJ1x58Eoj4PXUh/go-libp2p-peerstore"
	p2phost "gx/ipfs/QmfZTdmunzKzAGJrSvXXQbQ5kLLUiEMX5vdwux7iXkdk7D/go-libp2p-host"
)

// P2P structure holds information on currently running streams/listeners
type P2P struct {
	Listeners ListenerRegistry
	Streams   StreamRegistry

	identity  peer.ID
	peerHost  p2phost.Host
	peerstore pstore.Peerstore
}

// CheckProtoExists checks whether a proto handler is registered to
// mux handler
func (p2p *P2P) CheckProtoExists(proto string) bool {
	protos := p2p.peerHost.Mux().Protocols()

	for _, p := range protos {
		if p != proto {
			continue
		}
		return true
	}
	return false
}
