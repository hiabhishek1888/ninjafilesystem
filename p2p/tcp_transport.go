package p2p

import (
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

type TCPPeer struct {
	net.Conn
	// if we dial and retrive a conn => outbound == true
	// if we accept and retrive a conn => outbound == false
	outbound bool
	Wg       *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		Wg:       &sync.WaitGroup{},
	}
}
func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

func (p *TCPPeer) WgDone() {
	p.Wg.Done()
}

type TCPTransportOptions struct {
	ListenAddress string
	Handshake     HandshakeHandler
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOptions
	listner     net.Listener
	payloadChan chan Payload
	pathChan    chan string
}

func NewTCPTransport(opts TCPTransportOptions) *TCPTransport {
	return &TCPTransport{
		TCPTransportOptions: opts,
		payloadChan:         make(chan Payload),
		pathChan:            make(chan string),
	}
}

// consumePayload implements the transport interface which will return read-only channel
// for reading the Payload recieved from another peer in network for STORING THE DATA.
func (t *TCPTransport) ConsumePayload() <-chan Payload {
	return t.payloadChan
}

// consumePath implements the transport interface which will return read-only channel
// for reading the path recieved from another peer in network for GETTING THE DATA.
func (t *TCPTransport) ConsumePath() <-chan string {
	return t.pathChan
}

// no longer needed
func (p *TCPPeer) GetPeerforWaitGroup() *TCPPeer {
	return p
}

func (t *TCPTransport) CloseConnection() error {
	return t.listner.Close()
}

func (t *TCPTransport) Dail(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	fmt.Println(t.ListenAddress, " established connection with: ", addr)
	go t.HandleConnection(conn, true)
	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listner, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()
	log.Printf("TCP transport listening on port: %s \n", t.ListenAddress)
	return nil
}
func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listner.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s", err)
		}
		go t.HandleConnection(conn, false)
	}
}

func (t *TCPTransport) HandleConnection(conn net.Conn, outbound bool) {
	peer := NewTCPPeer(conn, outbound)
	defer conn.Close()
	if err := t.Handshake(peer); err != nil {
		fmt.Printf("TCP handshake error, dropping the peer connection: %s \n", err)
		return
	}
	if outbound && t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
		fmt.Println(t.ListenAddress, " added peer: ", conn.RemoteAddr().String())
	}
	// read the data - read loop

	for {
		fmt.Println("inside connhandler read loop")
		var p Payload
		if err := gob.NewDecoder(conn).Decode(&p); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("payload is: %+v", p)
		// fmt.Printf("payload path is: %+v", p.Path)

		// 1. if payload only contains path, means it is for reading the data from this peer.
		// 2. else if payload contains both path and data, means it is to write the data to this peer.
		if len(p.Data) == 0 {
			t.pathChan <- p.Path
			fmt.Println("path read and send to channel")
		} else {
			t.payloadChan <- p
			fmt.Printf("INSIDE CONN HANDLER: i am %+v and from %+v, recieved message is: %+v \n", conn.LocalAddr(), conn.RemoteAddr(), string(p.Data))
		}
	}

}
