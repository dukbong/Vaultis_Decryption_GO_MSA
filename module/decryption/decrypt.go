package decryption

import (
	"vaultis-go-module/module/aes"
	"vaultis-go-module/module/rsa"

	"github.com/gin-gonic/gin"
)

type RequestData struct {
	EncryptedMessage string `form:"encryptedMessage" json:"encryptedMessage"`
}

func Decrypt(c *gin.Context) {
	privateKey, err := rsa.ReadPrivateKey(c)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	requestData, err := getRequestData(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	aesKey, iv, encryptedContentBytes, err := extractDecryptionParameters(requestData.EncryptedMessage)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	aesKey, err = rsa.DecryptAESKey(privateKey, aesKey)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	plaintext, err := aes.DecryptContent(aesKey, iv, encryptedContentBytes)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"decrypted_content": string(plaintext)})
}
