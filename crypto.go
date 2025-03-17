package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type Encdec struct {
	enc           string
	decryptedData string
}

// https://www.twilio.com/en-us/blog/encrypt-and-decrypt-data-in-go-with-aes-256
func (x *Encdec) cryptofunction() {
	data := "private employee records restricted to personnel only"
	plaintext := []byte(data)
	key := make([]byte, 32)

	if _, err := rand.Reader.Read(key); err != nil {
		fmt.Println("error generating random encryption key ", err)
		return
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("error creating aes block cipher", err)
		return
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("error setting gcm mode", err)
		return
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("error generating the nonce ", err)
		return
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	x.enc = hex.EncodeToString(ciphertext)
	fmt.Println("original data:", data)
	fmt.Println("encrypted data:", x.enc)
	decodedCipherText, err := hex.DecodeString(x.enc)
	if err != nil {
		fmt.Println("error decoding hex", err)
		return
	}

	decryptedData, err := gcm.Open(nil, decodedCipherText[:gcm.NonceSize()], decodedCipherText[gcm.NonceSize():], nil)
	if err != nil {
		fmt.Println("error decrypting data", err)
		return
	}
	x.decryptedData = string(decryptedData)
	fmt.Println("Decrypted data:", string(decryptedData))
}

// Encrypt data using AES-256 => EncryptData takes raw data (of type []byte) and key (of type []byte) and encrypt the data with "nonce generated with key" using gcm and return the encrypted data (of type [byte]) and error.
//
// Read more: https://www.twilio.com/en-us/blog/encrypt-and-decrypt-data-in-go-with-aes-256
func EncryptData(rawDataByte []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("error creating aes block cipher", err)
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("error setting gcm mode", err)
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("error generating the nonce ", err)
		return nil, err
	}

	// encryptedDataByte is cipher text that is encrypted with "nonce generated with key" using gcm.
	encryptedDataByte := gcm.Seal(nonce, nonce, rawDataByte, nil)
	fmt.Println("cipherDataByte is: ", encryptedDataByte)
	fmt.Println("cipherDataByte string is: ", hex.EncodeToString(encryptedDataByte))
	return encryptedDataByte, nil
}

// Decrypt data using AES-256. => DecryptData takes encrypted data (of type []Byte) and key (of type []byte) and decrypt the encrypted data using gcm and nonce size and return the decrypted data (of type [] byte) and error
//
// Read more: https://www.twilio.com/en-us/blog/encrypt-and-decrypt-data-in-go-with-aes-256
func DecryptData(encryptedDataByte []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("error creating aes block cipher", err)
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("error setting gcm mode", err)
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println("error generating the nonce ", err)
		return nil, err
	}
	decryptedDataByte, err := gcm.Open(nil, encryptedDataByte[:gcm.NonceSize()], encryptedDataByte[gcm.NonceSize():], nil)
	if err != nil {
		fmt.Println("error decrypting data", err)
		return nil, err
	}
	fmt.Println("Decrypted data:", decryptedDataByte)
	fmt.Println("Decrypted data in string is:", string(decryptedDataByte))
	return decryptedDataByte, nil
}

func (x *Encdec) ExecuteCrypto() *Encdec {
	x.cryptofunction()

	key := []byte("2o3n07oek5q58u293035wthumma1n61x")
	ciphertext, _ := EncryptData([]byte("this is data"), key)
	DecryptData(ciphertext, key)

	return x
}
