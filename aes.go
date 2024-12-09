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

	for bs, be := 0, aes.BlockSize; bs < len(data); bs, be = bs+aes.BlockSize, be+aes.BlockSize {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted
}

func DecryptBlock(data, key []byte, blockId int) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	decrypted := make([]byte, aes.BlockSize)

	start := 0 + aes.BlockSize*blockId
	end := start + aes.BlockSize

	block.Decrypt(decrypted, data[start:end])

	return decrypted
}
