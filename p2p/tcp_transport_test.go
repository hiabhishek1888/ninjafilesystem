package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	tcpOpts := TCPTransportOptions{
		ListenAddress: ":3000",
		Handshake:     HandshakeHandlerFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(tcpOpts)
	assert.Equal(t, tr.ListenAddress, ":3000")
	// server
	assert.Nil(t, tr.ListenAndAccept())
}
