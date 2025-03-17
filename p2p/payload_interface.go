package p2p

// Payload interface holds path (of string type) and data (of []byte type) which can be send between 2 nodes (peers) in the network OVER network/transport interface.
type Payload struct {
	Path string
	Data []byte
}
