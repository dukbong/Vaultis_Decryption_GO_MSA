package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

func DecryptContent(aesKey, iv, encryptedContent []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("AES 암호화 블록 생성에 실패했습니다")
	}

	cbc := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(encryptedContent))
	cbc.CryptBlocks(plaintext, encryptedContent)

	return plaintext, nil
}
