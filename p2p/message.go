package p2p

// message holds any type of data which can be send over
// each transport between 2 nodes (peers) in the network

type Payload struct {
	Path string // added
	Data []byte // earlier it was payload
}
