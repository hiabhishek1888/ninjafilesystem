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

	ciphertext := gcm.Seal(nonce, nonce, rawDataByte, nil)
	fmt.Println("ciphertext is: ", ciphertext)
	fmt.Println("ciphertext string is: ", hex.EncodeToString(ciphertext))
	return ciphertext, nil
}
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
	decryptedData, err := gcm.Open(nil, encryptedDataByte[:gcm.NonceSize()], encryptedDataByte[gcm.NonceSize():], nil)
	if err != nil {
		fmt.Println("error decrypting data", err)
		return nil, err
	}
	fmt.Println("Decrypted data:", decryptedData)
	fmt.Println("Decrypted data in string is:", string(decryptedData))
	return decryptedData, nil
}

func (x *Encdec) ExecuteCrypto() *Encdec {
	x.cryptofunction()

	key := []byte("2o3n07oek5q58u293035wthumma1n61x")
	ciphertext, _ := EncryptData([]byte("this is data"), key)
	DecryptData(ciphertext, key)

	return x
}
