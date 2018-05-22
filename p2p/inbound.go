package p2p

import (
	"context"

	manet "gx/ipfs/QmRK2LxanhK2gZq6k6R7vk5ZoYZk8ULSSTB7FzDsMUX6CB/go-multiaddr-net"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	net "gx/ipfs/QmXfkENeeBvh3zYA51MaSdGUdBjhQ99cP5WQe8zgr6wchG/go-libp2p-net"
	protocol "gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
)

// inboundListener accepts libp2p streams and proxies them to a manet host
type inboundListener struct {
	p2p *P2P

	// Application proto identifier.
	proto string

	// Address to proxy the incoming connections to
	addr ma.Multiaddr
}

// NewListener creates new p2p listener
func (p2p *P2P) NewListener(ctx context.Context, proto string, addr ma.Multiaddr) (Listener, error) {
	listenerInfo := &inboundListener{
		p2p: p2p,

		proto: proto,
		addr:  addr,
	}

	p2p.peerHost.SetStreamHandler(protocol.ID(proto), func(remote net.Stream) {
		local, err := manet.Dial(addr)
		if err != nil {
			remote.Reset()
			return
		}

		stream := Stream{
			Protocol: proto,

			LocalPeer: p2p.identity,
			LocalAddr: addr,

			RemotePeer: remote.Conn().RemotePeer(),
			RemoteAddr: remote.Conn().RemoteMultiaddr(),

			Local:  local,
			Remote: remote,

			Registry: &p2p.Streams,
		}

		p2p.Streams.Register(&stream)
		stream.startStreaming()
	})

	p2p.Listeners.Register(listenerInfo)

	return listenerInfo, nil
}

func (l *inboundListener) Protocol() string {
	return l.proto
}

func (l *inboundListener) Address() string {
	return l.addr.String()
}

func (l *inboundListener) Close() error {
	l.p2p.peerHost.RemoveStreamHandler(protocol.ID(l.proto))
	l.p2p.Listeners.Deregister(l.proto)
	return nil
}
