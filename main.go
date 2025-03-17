package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/intogit/ninjafilesystem/p2p"
)

func makeServer(listenAddr string, remoteNodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOptions{
		ListenAddress: listenAddr,
		Handshake:     p2p.HandshakeHandlerFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	storeOpts := StoreOpts{
		PathTransform: CASPathTransformFunc,
	}
	st := NewStore(storeOpts)

	s := FileServer{
		transport:      tcpTransport,
		store:          st,
		quitchannel:    make(chan struct{}),
		bootstrapNodes: remoteNodes,
		peers:          make(map[string]p2p.Peer),
	}
	tcpTransport.OnPeer = s.OnPeer
	return &s

}

func main() {
	fmt.Println("Hi Ninja!!")

	fmt.Println("main - BEGINS")
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")
	s3 := makeServer(":2000", ":3000", ":4000")

	var wg sync.WaitGroup
	wg.Add(3)
	go s1.StartConn(&wg)
	go s2.StartConn(&wg)
	go s3.StartConn(&wg)
	wg.Wait()

	time.Sleep(2 * time.Second)

	fmt.Println(s1.peers)
	fmt.Println(s2.peers)
	fmt.Println(s3.peers)

	fmt.Println("can start store and get of data now..")

	// data1 := bytes.NewReader([]byte("this is small data file !!"))
	// s3.StoreData("myTempPath", data1)
	// data2 := bytes.NewReader([]byte("this is large data file !!"))
	// s2.StoreData("newTempfile", data2)

	// fmt.Println("waiting for store to be completed")
	// time.Sleep(20 * time.Second)
	x, err := s2.GetData("newTempfile")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("main: ", x)

	x, err = s3.GetData("myTempPath")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("main: ", x)

	fmt.Println("main - ENDS")
	select {}
}
