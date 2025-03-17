package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypto(t *testing.T) {
	encdec := &Encdec{}
	tr := encdec.ExecuteCrypto()
	assert.Equal(t, tr.enc, tr.decryptedData)
}
