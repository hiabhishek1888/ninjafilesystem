package p2p

import (
	"encoding/gob"
	"fmt"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *Payload) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, v any) error {
	return gob.NewDecoder(r).Decode(v)
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, msg *Payload) error {
	// Read data from the connection
	buffer := make([]byte, 1024) // Create a buffer to store incoming data
	n, err := r.Read(buffer)
	if err != nil {
		fmt.Println("Error reading data:", err)
		return err
	}
	fmt.Println("Decoded msg is:", string(buffer[:n])) // Convert and print received data
	// fmt.Printf("message length is: %+v\n", n)
	// no need to convert buffer data to string.. will increase overhead
	msg.Data = buffer[:n]
	return nil
}
