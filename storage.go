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
func PathTransformHandlerFunc(path string) string {
	return path
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
func (s *Store) Read(addr string, filepath string) (io.Reader, error) {
	return s.readStream(addr, filepath)
}
func (s *Store) Write(addr string, filepath string, r io.Reader) error {
	return s.writeStream(addr, filepath, r)
}
func (s *Store) Delete(filepath string) error {
	return s.deleteStream(filepath)
}

// private readstream, writeStream and deleteStream fxn/api's below
func (s *Store) deleteStream(filepath string) error {
	transformedFilePath, fileName := s.PathTransform(filepath)
	fmt.Println("path name is \n", transformedFilePath)
	fmt.Println("file name is \n", fileName)

	// Delete only the file, not the directory/folder structure
	// we need fullfilepath
	// fullFilePath := transformedFilePath + "/" + fileName
	// err := os.RemoveAll(fullFilePath)
	// if err != nil {
	// 	fmt.Printf("Unable to delete the file: %s ", filepath)
	// 	return err
	// }
	// fmt.Printf("File: %s, deleted successfully", filepath)
	// return nil

	// Delete whole folder and all its child folder and files
	// we need root folder which is transformedfile- first part
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
		// } else if err != nil {
		// 	return false
	}
	return true
}

func (s *Store) readStream(addr string, filepath string) (io.Reader, error) {

	// checking if filepath exist
	// if !s.HasPath(addr, filepath) {
	// 	err := errors.New("path do not exist")
	// 	return nil, err
	// }

	transformedFilePath, fileName := s.PathTransform(filepath)
	transformedFilePath = addr[1:] + "/" + transformedFilePath
	// fmt.Println("path name is \n", transformedFilePath)
	// fmt.Println("file name is \n", fileName)

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

	// CANNOT DO BELOW ?? WHY ?
	// BECAUSE we wont be able to find the file, when we need to read.
	// TODO:
	// but it is CAS feature.. no duplication.
	// current implementation is overwriting the existing data if we give the same path
	// but in CAS, it create a new file and leave the existing one as it is.

	// buf := make([]byte(buffer))
	// buf := new(bytes.Buffer)
	// io.Copy(buf, r)
	// filenameByte := md5.Sum(buf.Bytes())
	// fileName := hex.EncodeToString(filenameByte[:])
	// fullFilePath := transformedFilePath + "/" + fileName

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
