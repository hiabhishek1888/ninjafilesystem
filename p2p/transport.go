package p2p

import "net"

// peer represent a remote connected node.
type Peer interface {
	Send([]byte) error
	// RemoteAddress() net.Addr
	net.Conn
	// Close() error
	WgDone()
}

// transport will handle the communication between nodes in the p2p network
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan Payload
	CloseConnection() error
}
