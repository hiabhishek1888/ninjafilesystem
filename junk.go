package main

// MAIN.GO JUNKS

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

//
// RECIEVEDATASTORE METHOD JUNK

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

//
//
//
//
//
//
//
// ListenAndAccept()
//
// func (t *TCPTransport) ListenAndAccept() error {
// 	var err error
// 	t.listner, err = net.Listen("tcp", t.ListenAddress)
// 	if err != nil {
// 		return err
// 	}
// 	// fmt.Printf("server created on localhost: %s\n", t.ListenAddress)

// 	// for {
// 	// 	conn, err := t.listner.Accept()
// 	// 	if err != nil {
// 	// 		fmt.Printf("TCP accept error: %s", err)
// 	// 		continue
// 	// 	}
// 	// 	go t.handleConnection(conn)
// 	// }

// 	// THE ABOVE CODE IS BLOCKING MAIN,
// 	// BCOZ IT IS GOING TO CONTINOUSLY ACCEPT NEW CONNECTIONS
// 	// AND WE NEED TO HANDLE THEM.. SO BASICALLY IT'S INFINITE LOOP
// 	// TO HANDLE THIS.. WE NEED TO PUT THE ABOVE LOOP IN SEPERATE GOROUTINE

// 	// Approach 1:
// 	go t.startAcceptLoop()
// 	log.Printf("TCP transport listening on port: %s \n", t.ListenAddress)
// 	return nil

// 	// Approach 2:
// 	// go func() {
// 	// 	for {
// 	// 		conn, err := t.listner.Accept()
// 	// 		if err != nil {
// 	// 			fmt.Printf("TCP accept error: %s", err)
// 	// 			continue
// 	// 		}
// 	// 		go t.handleConnection(conn)
// 	// 	}
// 	// }()
// 	// return nil

// 	// Approach 3: // handleconn should be defined in another goroutine for each in anonymous func
// 	// for {
// 	// 	conn, err := t.listner.Accept()
// 	// 	if err != nil {
// 	// 		fmt.Printf("TCP accept error: %s", err)
// 	// 		continue
// 	// 	}
// 	// 	go func(c net.Conn) {
// 	// 		defer c.Close()
// 	// 		fmt.Printf("new incomming connection to be handled: %+v\n", c)
// 	// 	}(conn)
// 	// }
// }

//

//
// func (t *TCPTransport) HandleConnection(conn net.Conn, outbound bool) {
// 	peer := NewTCPPeer(conn, outbound)
// 	defer conn.Close()
// 	if err := t.Handshake(peer); err != nil {
// 		fmt.Printf("TCP handshake error, dropping the peer connection: %s \n", err)
// 		return
// 	}
// 	if outbound && t.OnPeer != nil {
// 		if err := t.OnPeer(peer); err != nil {
// 			return
// 		}
// 		fmt.Println(t.ListenAddress, " added peer: ", conn.RemoteAddr().String())
// 	}
// 	// read the data - read loop

// 	// for {
// 	// 	// Read data from the connection
// 	// 	buffer := make([]byte, 1024) // Create a buffer to store incoming data
// 	// 	msg, err := conn.Read(buffer)
// 	// 	if err != nil {
// 	// 		fmt.Println("Error reading data:", err)
// 	// 		return
// 	// 	}
// 	// 	fmt.Println("Received msg is:", string(buffer[:msg])) // Convert and print received data
// 	// 	fmt.Printf("message length is: %+v\n", msg)
// 	// }

// 	//TODO (resolve below error):
// 	// when connection is closed from remote node in the middle of data read,
// 	// IT CREATES the infinite loop of error: "TCP error: read TCP <remote ip>... : use of closed network connection"..
// 	// "An existing connection was forcibly closed by the remote host."

// 	// msg := Payload{}
// 	for {
// 		fmt.Println("inside connhandler")
// 		var p Payload
// 		if err := gob.NewDecoder(conn).Decode(&p); err != nil {
// 			log.Fatal(err)
// 		}
// 		fmt.Printf("payload is: %+v", p)
// 		fmt.Printf("payload path is: %+v", p.Path)
// 		if len(p.Data) == 0 {
// 			fmt.Println("hello whyy ?")
// 			t.pathChan <- p.Path
// 			fmt.Println(len(t.pathChan))
// 		} else {

// 			// peer.Wg.Add(1)
// 			// fmt.Println("waiting for the data to be read and write in disk")
// 			t.payloadChan <- p
// 			// peer.Wg.Wait()
// 			// fmt.Println("stream done, continuinng normal read loop")

// 			// err := t.Decoder.Decode(conn, &msg)
// 			// if err != nil {
// 			// 	fmt.Printf("Could not decode the message from the remote node: %s\n", err)
// 			// 	continue
// 			// 	// return
// 			// 	// using return when error occured... this could be a way to prevent the above error
// 			// 	// but if the error is due to something else..
// 			// 	// we cannot be able to read that again
// 			// }
// 			// // msg.From = conn.RemoteAddr()
// 			fmt.Printf("INSIDE CONN HANDLER: i am %+v and from %+v, recieved message is: %+v\n", conn.LocalAddr(), conn.RemoteAddr(), string(p.Data))
// 			// t.payloadChan <- msg
// 			// }
// 		}
// 	}

// }

//

//

//

//
//
//
//
//
// func Test() {
// 	t := NewTCPTransport(":4344")
// 	t.listner.Accept()
// }

// YOU CAN ALSO RETURN INTERFACE TYPE FROM A FUCNTION
// BUT ACCORDINGLY YOU HAVE TO MODIFY OTHER FUCNTION OUTPUT USE.
// below is the example

// func NewTCPTransport(listenAddr string) Transport {
//	// #Transport is interface
// 	return &TCPTransport{listenAddress: listenAddr}
// }

// func Test() {
// 	// t := NewTCPTransport(":4344")
//  // # now above code wont work, you need to resolve the type of output you want from an interface.
//	t := NewTCPTransport(":4344").(*TCPTransport)
// 	t.listner.Accept()
// }
