package aes

import "bytes"

var (
	SecretKey = []byte("2985BCFDB5FE43129843DB59825F8647")
)

func PKCS5Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}
