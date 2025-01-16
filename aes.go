package main

import (
	"crypto/aes"
)

func DecryptBlocks(data, key []byte, length int) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	decrypted := make([]byte, min(aes.BlockSize*length, len(data)))

	for bs, be := 0, aes.BlockSize; bs < len(decrypted); bs, be = bs+aes.BlockSize, be+aes.BlockSize {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted
}
