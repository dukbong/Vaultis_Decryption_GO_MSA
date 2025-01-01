package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/decryption", decryption)
	r.Run(":8080")
}

type RequestData struct {
	Content string `form:"content" json:"content"`
}

func decryption(c *gin.Context) {
	// privateKey 파일을 읽어오는 부분을 별도의 함수로 분리
	privateKey, err := readPrivateKey(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 암호문을 받아오는 부분을 함수로 분리
	requestData, err := getRequestData(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// AES 키와 IV, 암호화된 본문 추출
	aesKey, iv, encryptedContentBytes, err := extractDecryptionParameters(requestData.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// AES 키 복호화
	aesKey, err = rsaDecryptAESKey(privateKey, aesKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 본문 복호화
	plaintext, err := decryptAESContent(aesKey, iv, encryptedContentBytes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(">>>")
	fmt.Println(string(plaintext))
	fmt.Println(">>>")

	// 복호화된 데이터 출력 및 응답
	c.JSON(http.StatusOK, gin.H{"decrypted_content": string(plaintext)})
}

func readPrivateKey(c *gin.Context) (*rsa.PrivateKey, error) {
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

func getRequestData(c *gin.Context) (*RequestData, error) {
	var requestData RequestData
	if err := c.ShouldBind(&requestData); err != nil {
		return nil, fmt.Errorf("유효하지 않은 폼 데이터입니다.")
	}
	return &requestData, nil
}

func extractDecryptionParameters(content string) ([]byte, []byte, []byte, error) {
	parts := strings.Split(content, ".")
	if len(parts) != 3 {
		return nil, nil, nil, fmt.Errorf("잘못된 암호문 형식입니다.")
	}

	// AES key 추출 및 디코딩
	encryptedAESKey, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("암호화된 AES 키 디코딩에 실패했습니다.")
	}

	// IV 추출 및 디코딩
	iv, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("IV 디코딩에 실패했습니다.")
	}

	// 암호화된 본문 추출 및 디코딩
	encryptedContentBytes, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("암호화된 본문 디코딩에 실패했습니다.")
	}

	return encryptedAESKey, iv, encryptedContentBytes, nil
}

func rsaDecryptAESKey(privateKey *rsa.PrivateKey, encryptedAESKey []byte) ([]byte, error) {
	aesKey, err := rsa.DecryptPKCS1v15(nil, privateKey, encryptedAESKey)
	if err != nil {
		return nil, fmt.Errorf("AES 키 복호화에 실패했습니다.")
	}
	return aesKey, nil
}

func decryptAESContent(aesKey, iv, encryptedContent []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("AES 암호화 블록 생성에 실패했습니다")
	}

	cbc := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(encryptedContent))
	cbc.CryptBlocks(plaintext, encryptedContent)

	return plaintext, nil
}
