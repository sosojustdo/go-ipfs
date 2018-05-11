package p2p

import (
	"context"
	"time"

	manet "gx/ipfs/QmRK2LxanhK2gZq6k6R7vk5ZoYZk8ULSSTB7FzDsMUX6CB/go-multiaddr-net"
	ma "gx/ipfs/QmWWQ2Txc2c6tqjsBpzg5Ar652cHPGNsQQp2SejkNmkUMb/go-multiaddr"
	pstore "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	net "gx/ipfs/QmXfkENeeBvh3zYA51MaSdGUdBjhQ99cP5WQe8zgr6wchG/go-libp2p-net"
	protocol "gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	peer "gx/ipfs/QmZoWKhxUmZ2seW4BzX6fJkNR8hh9PsGModr7q171yq2SS/go-libp2p-peer"
)

// inboundListener accepts libp2p streams and proxies them to a manet host
type outboundListener struct {
	ctx    context.Context
	cancel context.CancelFunc

	p2p *P2P
	id  peer.ID

	proto string
	peer  peer.ID

	listener manet.Listener
}

// Dial creates new P2P stream to a remote listener
func (p2p *P2P) Dial(ctx context.Context, peer peer.ID, proto string, bindAddr ma.Multiaddr) (Listener, error) {
	maListener, err := manet.Listen(bindAddr)
	if err != nil {
		return nil, err
	}

	listener := &outboundListener{
		p2p: p2p,
		id:  p2p.identity,

		proto: proto,
		peer:  peer,

		listener: maListener,
	}

	go listener.acceptConns()

	return listener, nil
}

func (l *outboundListener) dial() (net.Stream, error) {
	ctx, cancel := context.WithTimeout(l.ctx, time.Second*30) //TODO: configurable?
	defer cancel()

	err := l.p2p.peerHost.Connect(ctx, pstore.PeerInfo{ID: l.peer})
	if err != nil {
		return nil, err
	}

	return l.p2p.peerHost.NewStream(l.ctx, l.peer, protocol.ID(l.proto))
}

func (l *outboundListener) acceptConns() {
	for {
		local, err := l.listener.Accept()
		if err != nil {
			return
		}

		remote, err := l.dial()
		if err != nil {
			local.Close()
			return
		}

		stream := Stream{
			Protocol: l.proto,

			LocalPeer: l.id,
			LocalAddr: l.listener.Multiaddr(),

			RemotePeer: remote.Conn().RemotePeer(),
			RemoteAddr: remote.Conn().RemoteMultiaddr(),

			Local:  local,
			Remote: remote,

			Registry: &l.p2p.Streams,
		}

		l.p2p.Streams.Register(&stream)
		stream.startStreaming()
	}
}

func (l *outboundListener) Close() error {
	l.listener.Close()
	l.p2p.Listeners.Deregister(l.proto)
	return nil
}

func (l *outboundListener) Protocol() string {
	return l.proto
}

func (l *outboundListener) Address() string {
	return "/ipfs/" + l.peer.String()
}
