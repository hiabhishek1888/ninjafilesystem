package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type PathTransformHandler func(string) (string, string)

// `CASPathTransformFunc` return a hashed and formatted filepath, which is easy to os operations in it.
//
//	`CASPathTransformFunc` takes string type file path and converts it to hash and then encode the hash to string and split it with `blocksize` length, ie: 20 and then joins them using '/'.
func CASPathTransformFunc(path string) (string, string) {
	hash := sha1.Sum([]byte(path))
	hashStr := hex.EncodeToString(hash[:])
	// converts to string of len 40
	blocksize := 20
	slicelen := len(hashStr) / blocksize // 8
	paths := make([]string, slicelen)
	for i := 0; i < slicelen; i++ {
		from, to := i*blocksize, i*blocksize+blocksize
		paths[i] = hashStr[from:to]
	}
	return strings.Join(paths, "/"), hashStr
}

type StoreOpts struct {
	PathTransform PathTransformHandler
}
type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	return &Store{StoreOpts: opts}
}

// public function of below implementations

// Read will take the addr(machine address, will act as root path for each machine) and filepath. It goes to read from that machine disk and return data or error, if any.
func (s *Store) Read(addr string, filepath string) (io.Reader, error) {
	return s.readStream(addr, filepath)
}

// Write will take the addr(machine address, will act as root path for each machine) and filepath and reader. It goes on to read data from reader and store that data in machine disk at location of filepath and return only error, if any.
func (s *Store) Write(addr string, filepath string, r io.Reader) error {
	return s.writeStream(addr, filepath, r)
}

// Delete will take filepath. It goes to check if that machine disk have that file and delete it and return error if any.
func (s *Store) Delete(filepath string) error {
	return s.deleteStream(filepath)
}

// private readstream, writeStream and deleteStream fxn/api's below

func (s *Store) deleteStream(filepath string) error {
	// TYPE 1:
	// transformedFilePath, fileName := s.PathTransform(filepath)
	// Delete only the file, not the directory/folder structure we need fullfilepath
	//
	// fullFilePath := transformedFilePath + "/" + fileName
	// err := os.RemoveAll(fullFilePath)
	// if err != nil {
	// 	fmt.Printf("Unable to delete the file: %s ", filepath)
	// 	return err
	// }
	// fmt.Printf("File: %s, deleted successfully", filepath)
	// return nil

	//TYPE 2:
	// Delete whole folder and all its child folder and files we need root folder which is transformedfile- first part
	transformedFilePath, _ := s.PathTransform(filepath)
	rootFolder := strings.Split(transformedFilePath, "/")[0]
	err := os.RemoveAll(rootFolder)
	if err != nil {
		fmt.Printf("Unable to delete the root of file: %s , i.e: %s", filepath, rootFolder)
		return err
	}
	fmt.Printf("Rootfolder: %s, of file %s, deleted successfully \n", rootFolder, filepath)
	return nil
}
func (s *Store) HasPath(addr string, filepath string) bool {
	transformedFilePath, fileName := s.PathTransform(filepath)
	transformedFilePath = addr[1:] + "/" + transformedFilePath
	fullFilePath := transformedFilePath + "/" + fileName

	_, err := os.Stat(fullFilePath)
	if os.IsNotExist(err) {
		log.Println("filepath do not exist in disk")
		return false
	}
	return true
}
func (s *Store) readStream(addr string, filepath string) (io.Reader, error) {

	// Checking if filepath exist
	// if !s.HasPath(addr, filepath) {
	// 	err := errors.New("path do not exist")
	// 	return nil, err
	// }

	transformedFilePath, fileName := s.PathTransform(filepath)
	transformedFilePath = addr[1:] + "/" + transformedFilePath

	fullFilePath := transformedFilePath + "/" + fileName
	f, err := os.Open(fullFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	n, err := io.Copy(buf, f)
	if err != nil {
		return nil, err
	}
	fmt.Printf("read %d byte from disk => filename is %s \n Read data is: %s", n, fullFilePath, buf)
	return buf, nil
}
func (s *Store) writeStream(addr string, filepath string, r io.Reader) error {
	transformedFilePath, fileName := s.PathTransform(filepath)
	transformedFilePath = addr[1:] + "/" + transformedFilePath
	// fmt.Println("path name is \n", transformedFilePath)
	// fmt.Println("file name is \n", fileName)
	if err := os.MkdirAll(transformedFilePath, os.ModePerm); err != nil {
		return err
	}

	// TODO:
	// CAS snapshot feature support !!
	// Current implementation is overwriting the existing data if we try to modify the data. (path of data is same)
	// See more: todo.txt

	fullFilePath := transformedFilePath + "/" + fileName
	f, err := os.Create(fullFilePath)
	if err != nil {
		return err
	}
	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}
	fmt.Printf("written %d byte to disk => filename is %s \n", n, fullFilePath)
	return nil
}
