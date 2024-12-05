package main

import (
	"crypto/aes"
)

func Decrypt(data, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	decrypted := make([]byte, len(data))
	size := 16

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted
}

func DecryptBlock(data, key []byte, blockId int) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	decrypted := make([]byte, 16)

	start := 0 + 16*blockId
	end := start + 16

	block.Decrypt(decrypted, data[start:end])

	return decrypted
}
