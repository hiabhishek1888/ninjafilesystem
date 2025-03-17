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

type FileServer struct {
	transport      *p2p.TCPTransport
	store          *Store
	quitchannel    chan struct{}
	bootstrapNodes []string

	peerLock sync.Mutex
	peers    map[string]p2p.Peer
}

// Adding all "connected peer n/w" in peer map for a single machine/device
func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p

	// log.Printf("peer local address is: %s and remote address is: %s", s.transport.ListenAddress, p.RemoteAddr().String())
	return nil
}

var readdata io.Reader
var readdataerr error
var recievedData (chan bool)

func (s *FileServer) RecievePathAndReturnData() {
	defer s.transport.CloseConnection()
	for {
		fmt.Println("hey there, RecievePathAndReturnData is active !!")
		select {
		case path := <-s.transport.ConsumePath():
			addr := s.transport.TCPTransportOptions.ListenAddress
			if s.store.HasPath(addr, path) {
				readdata, readdataerr = s.store.Read(addr, path)
			}
			fmt.Printf("%s remote node data returned/saved, recieved over channel \n", addr)
			recievedData <- true
		case <-s.quitchannel:
			return
		}
	}
}

func (s *FileServer) RecieveDataAndStore() {
	defer s.transport.CloseConnection()
	for {
		fmt.Println("hey there, RecieveDataAndStore is active !!")
		select {
		case msg := <-s.transport.ConsumePayload():
			addr := s.transport.TCPTransportOptions.ListenAddress
			if err := s.store.Write(addr, msg.Path, bytes.NewReader(msg.Data)); err != nil {
				log.Panic(err)
			}
			fmt.Printf("%s remote node data stored, recieved over channel from  \n", addr)
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

func (s *FileServer) gatherData(path string) error {
	p := p2p.Payload{
		Path: path,
		Data: nil,
	}
	// p := []byte(path)
	peers := []io.Writer{} // peer should implement Writer to send the data over the network
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}
	// Sending path to each peer...
	for _, peer := range peers {
		// Will encode the data in gob format and send over the network
		if err := gob.NewEncoder(peer).Encode(p); err != nil {
			return err
		}
	}
	return nil
}
func (s *FileServer) GetData(path string) (io.Reader, error) {
	// 1. Check if file path exist on requestor disk, if yes then we return Read data
	// 2. else, Requesting peers to check and share the file if the file exist with them.

	// var wg sync.WaitGroup
	// wg.Add(1)
	// go s.RecievePathAndReturnData(&wg)
	// go s.RecievePathAndReturnData()
	// WAIT GROUP IS NOT WORKING INSIDE THE RecievePathAndReturnData

	fmt.Println("GetData initiated for path ", path)
	// 1. checking file locally
	addr := s.transport.TCPTransportOptions.ListenAddress
	if s.store.HasPath(addr, path) {
		return s.store.Read(addr, path)
	}
	fmt.Printf("(%s) file do not exist locally, fetching from peer networks \n", path)

	// 2. Requesting peers to share the file, if they have.

	// initialising the global variables, helps to get the data and wait till data recieved
	readdata = nil                 // recieve the data
	readdataerr = nil              // recieve error, if any
	recievedData = make(chan bool) // wait till data is consumed from channel and written to above vars.

	s.gatherData(path)

	fmt.Println("Waiting for data to be recieved and close the -recievedData- channel")
	<-recievedData

	if readdataerr != nil {
		return nil, readdataerr
	}
	// decrypting the data
	buf := new(bytes.Buffer)
	buf.ReadFrom(readdata)
	key := []byte("2o3n07oek5q58u293035wthumma1n61x")
	decryptedData, err := DecryptData(buf.Bytes(), key)
	if err != nil {
		fmt.Println("could not decrypt the data")
	}
	fmt.Println("decrypted file:=> ", decryptedData)
	// return readdata, nil
	return bytes.NewReader(decryptedData), nil
}

func (s *FileServer) broadcastData(p *p2p.Payload) error {
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}
	// mw will create a single multiwriter encoder for all peers
	mw := io.MultiWriter(peers...)
	// Will encode the data in gob format and send over the network
	if err := gob.NewEncoder(mw).Encode(p); err != nil {
		return err
	}
	return nil
}
func (s *FileServer) StoreData(path string, r io.Reader) error {
	// 1. Store the file to requestor disk
	// 2. Broadcast this file to all known peers in n/w

	// Copying the reader data because we need to consume twice,
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
	// encrypt the data
	key := []byte("2o3n07oek5q58u293035wthumma1n61x")
	encryptedData, err := EncryptData(buf.Bytes(), key)
	if err != nil {
		fmt.Println("could not encrypt the data")
	}
	p := &p2p.Payload{
		Path: path,
		// Data: buf.Bytes(),
		Data: encryptedData,
	}
	// fmt.Println("sending payload to remote node: ", p)
	return s.broadcastData(p)
}
func (s *FileServer) StopSendingData() {
	close(s.quitchannel)
}
