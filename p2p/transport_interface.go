package p2p

import "net"

// Peer interface represent a remote connected node.
//
// It implements Send method and have a waitgroup done().
//
// And since it is peer means some other device established connection with it over n/w, so it must have connection address. that is net.Conn

type Peer interface {
	Send([]byte) error
	net.Conn
	WgDone()
}

// Transport interface will handle the communication between 2 nodes in the p2p network.

// + `ListenAndAccept` => handles the connection
// + `ConsumePayload` => transfer the payload (path + data) over channel to server to store the payload
// + `ConsumePath` => transfer the path over channel to server to read the data from disk and return
// + `CloseConnection` => closes the connection
type Transport interface {
	ListenAndAccept() error
	ConsumePayload() <-chan Payload
	ConsumePath() <-chan string
	CloseConnection() error
}
