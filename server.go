package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/intogit/ninjafilesystem/p2p"
)

//	type FileServerOpts struct {
//		PathTransform PathTransformHandler
//		Transport     p2p.TCPTransport
//	}
type FileServer struct {
	transport      *p2p.TCPTransport
	store          *Store
	quitchannel    chan struct{}
	bootstrapNodes []string

	peerLock sync.Mutex
	peers    map[string]p2p.Peer
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p

	log.Printf("peer local address is: %s and remote address is: %s", s.transport.ListenAddress, p.RemoteAddr().String())
	return nil
}

var readdata io.Reader = nil
var readdataerr error = nil
var recievedData = make(chan bool)

func (s *FileServer) RecievePathAndReturnData() {
	defer s.transport.CloseConnection()
	for {
		// fmt.Println("hey, i am here !!")
		select {
		case path := <-s.transport.ConsumePath():
			addr := s.transport.TCPTransportOptions.ListenAddress
			if s.store.HasPath(addr, path) {
				readdata, readdataerr = s.store.Read(addr, path)
			}
			// fmt.Println("readdata is: ", readdata)
			// time.Sleep(10 * time.Second)
			recievedData <- true
		case <-s.quitchannel:
			return
		}
	}
}

func (s *FileServer) RecieveDataAndStore() {
	// addr := s.transport.TCPTransportOptions.ListenAddress
	defer s.transport.CloseConnection()
	for {
		select {
		case msg := <-s.transport.ConsumePayload():
			// var p p2p.Payload
			// err := s.transport.Decoder.Decode(bytes.NewReader(msg.Payload), &p)
			// if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&p); err != nil {
			// 	log.Fatal(err)
			// }
			// if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(&p); err != nil {
			// 	log.Fatal(err)
			// }

			// fmt.Println("msge is: \n", msg)
			// fmt.Printf("decoded path: is %+v \n", string(msg.Path))
			// fmt.Printf("decoded data: is %+v \n", string(msg.Data))

			addr := s.transport.TCPTransportOptions.ListenAddress
			if err := s.store.Write(addr, msg.Path, bytes.NewReader(msg.Data)); err != nil {
				log.Panic(err)
			}

			fmt.Printf("%s node data stored the recieved data \n", addr)

			// peer := &p2p.TCPPeer{}
			// peer.GetPeerforWaitGroup().Wg.Done()
			// peer.WgDone()

		case <-s.quitchannel:
			return
		}
	}
}
func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.bootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			if err := s.transport.Dail(addr); err != nil {
				log.Println("dial error", err)
			}
		}(addr)
	}
	return nil
}
func (s *FileServer) StartConn(wg *sync.WaitGroup) error {
	defer func() {
		log.Println("file server stopped due to user action")
		s.transport.CloseConnection()
	}()
	if err := s.transport.ListenAndAccept(); err != nil {
		return err
	}
	s.bootstrapNetwork()
	wg.Done()
	go s.RecieveDataAndStore()
	go s.RecievePathAndReturnData()
	return nil
}

// type Payload struct {
// 	Path string
// 	Data []byte
// }

func (s *FileServer) gatherData(path string) error {
	p := p2p.Payload{
		Path: path,
		Data: nil,
	}
	// p := []byte(path)
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}
	for _, peer := range peers {
		if err := gob.NewEncoder(peer).Encode(p); err != nil {
			return err
		}
	}
	return nil
}
func (s *FileServer) GetData(path string) (io.Reader, error) {
	// checking if file path exist on our disk, if yes then we return Read
	// var wg sync.WaitGroup
	// wg.Add(1)
	// go s.RecievePathAndReturnData(&wg)

	// go s.RecievePathAndReturnData()
	// time.Sleep(5 * time.Second)

	addr := s.transport.TCPTransportOptions.ListenAddress
	if s.store.HasPath(addr, path) {
		return s.store.Read(addr, path)
	}
	fmt.Printf("we do not have the file (%s) locally, fetching from peer networks \n", path)

	// Requesting peers to check and share the file, if they have.
	s.gatherData(path)
	fmt.Println("here waiting for data to be read and transfered to channel.. ")
	// fmt.Println("waiting for data to be recieved")

	// wg.Wait()
	// time.Sleep(10 * time.Second)
	<-recievedData
	if readdataerr != nil {
		return nil, readdataerr
	}

	return readdata, nil
}

func (s *FileServer) broadcastData(p *p2p.Payload) error {
	fmt.Println("INDISE BROADCAST")
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}
	mw := io.MultiWriter(peers...)
	if err := gob.NewEncoder(mw).Encode(p); err != nil {
		return err
	}
	fmt.Println("outside BROADCAST")
	return nil
}
func (s *FileServer) StoreData(path string, r io.Reader) error {
	// 1. Store the file to disk
	// 2. Broadcast this file to all known peers in n/w

	// copying the reader data because we need to consume twice,
	// once to write in our disk and second to broadcast the data

	buf := new(bytes.Buffer)
	// _, err := io.Copy(buf, r)
	// if err != nil {
	// 	return err
	// }
	tee := io.TeeReader(r, buf)
	addr := s.transport.TCPTransportOptions.ListenAddress
	if err := s.store.Write(addr, path, tee); err != nil {
		return err
	}
	p := &p2p.Payload{
		Path: path,
		Data: buf.Bytes(),
	}

	// fmt.Println(buf.Bytes())
	// fmt.Println(buf.String())
	// fmt.Println("sending payload to remote node: ", p)
	return s.broadcastData(p)
}
func (s *FileServer) StopSendingData() {
	close(s.quitchannel)
}

// func (s *FileServer) CloseConn() error {
// 	return s.transport.CloseConnection()
// }
