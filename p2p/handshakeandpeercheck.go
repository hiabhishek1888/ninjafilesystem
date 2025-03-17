package p2p

type HandshakeHandler func(Peer) error

// Needs to be implemented !! It simply returns nil(no error), means handshake is always successful.
func HandshakeHandlerFunc(Peer) error { return nil }

type OnPeerHandler func(Peer) error

func OnPeerHandlerFunc(Peer) error { return nil }
