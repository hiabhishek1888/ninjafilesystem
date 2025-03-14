package p2p

type HandshakeHandler func(Peer) error

func HandshakeHandlerFunc(Peer) error { return nil }

type OnPeerHandler func(Peer) error

func OnPeerHandlerFunc(Peer) error { return nil }
