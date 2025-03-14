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
	// tcpOpts := p2p.TCPTransportOptions{
	// 	ListenAddress: ":3000",
	// 	Handshake:     p2p.HandshakeHandlerFunc,
	// 	Decoder:       p2p.DefaultDecoder{},
	// 	PeerCheck:     p2p.OnPeerHandlerFunc,
	// }
	// tr := p2p.NewTCPTransport(tcpOpts)

	// go func() {
	// 	for {
	// 		msg := <-tr.Consume()
	// 		fmt.Printf("consumed data is %+v: \n", msg)
	// 	}
	// }()

	// if err := tr.ListenAndAccept(); err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("before select")
	// select {}

	// tcpTransportOpts := p2p.TCPTransportOptions{
	// 	ListenAddress: ":3000",
	// 	Handshake:     p2p.HandshakeHandlerFunc,
	// 	Decoder:       p2p.DefaultDecoder{},
	// 	PeerCheck:     p2p.PeerCheckHandlerFunc,
	// }
	// tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	// storeOpts := StoreOpts{
	// 	PathTransform: CASPathTransformFunc,
	// }
	// st := NewStore(storeOpts)

	// s := FileServer{
	// 	transport:      tcpTransport,
	// 	store:          st,
	// 	quitchannel:    make(chan struct{}),
	// 	bootstrapNodes: []string{":4000"},
	// }

	// // go func() {
	// // 	time.Sleep(time.Second * 20)
	// // 	s.StopSendingData()
	// // }()

	// if err := s.StartConn(); err != nil {
	// 	log.Fatal(err)
	// }

	// // go func() {
	// // 	time.Sleep(time.Second * 5)
	// // 	err := s.CloseConn()
	// // 	if err != nil {
	// // 		log.Fatal(err)
	// // 	}
	// // }()

	// select {}
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

	time.Sleep(5 * time.Second)

	fmt.Println(s1.peers)
	fmt.Println(s2.peers)
	fmt.Println(s3.peers)

	fmt.Println("will start sending data now..")
	// time.Sleep(10 * time.Second)

	// data1 := bytes.NewReader([]byte("this is small data file !!"))
	// s3.StoreData("myTempPath", data1)
	// data2 := bytes.NewReader([]byte("this is large data file !!"))
	// s2.StoreData("newTempfile", data2)

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
