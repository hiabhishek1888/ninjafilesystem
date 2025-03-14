package main

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransform: CASPathTransformFunc,
	}
	s := NewStore(opts)
	filedata := bytes.NewReader([]byte("it is rose image bytes"))
	filepath := "myspecialpictures"

	err := s.writeStream(":3000", filepath, filedata)
	if err != nil {
		log.Fatal(err)
	}
	r, err := s.readStream(":3000", filepath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s \n", r)
	// err := s.deleteStream(filepath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("file deleteed successfully")
}
