package rsa

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

func ReadPrivateKey(c *gin.Context) (*rsa.PrivateKey, error) {
	privateKeyFile, _, err := c.Request.FormFile("privateKeyFile")
	if err != nil {
		return nil, fmt.Errorf("키 파일이 필요합니다.")
	}
	defer privateKeyFile.Close()

	privateKeyFileContent, err := io.ReadAll(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("키 파일을 읽는 데 실패했습니다.")
	}

	privateKeyContent := strings.ReplaceAll(string(privateKeyFileContent), " ", "")
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyContent)
	if err != nil {
		return nil, fmt.Errorf("키 파일 변환 중 실패했습니다.")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("키 파일을 파싱하는 데 실패했습니다.")
	}

	return privateKey.(*rsa.PrivateKey), nil
}

func DecryptAESKey(privateKey *rsa.PrivateKey, encryptedAESKey []byte) ([]byte, error) {
	aesKey, err := rsa.DecryptPKCS1v15(nil, privateKey, encryptedAESKey)
	if err != nil {
		return nil, fmt.Errorf("AES 키 복호화에 실패했습니다.")
	}
	return aesKey, nil
}
