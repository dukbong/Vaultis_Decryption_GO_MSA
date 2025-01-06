package decryption

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func getRequestData(c *gin.Context) (*RequestData, error) {
	var requestData RequestData
	if err := c.ShouldBind(&requestData); err != nil {
		return nil, fmt.Errorf("유효하지 않은 폼 데이터입니다.")
	}
	return &requestData, nil
}

func extractDecryptionParameters(encryptedMessage string) ([]byte, []byte, []byte, error) {
	parts := strings.Split(encryptedMessage, ".")
	if len(parts) != 3 {
		return nil, nil, nil, fmt.Errorf("잘못된 암호문 형식입니다.")
	}

	encryptedAESKey, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("암호화된 AES 키 디코딩에 실패했습니다.")
	}

	iv, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("IV 디코딩에 실패했습니다.")
	}

	encryptedContentBytes, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("암호화된 본문 디코딩에 실패했습니다.")
	}

	return encryptedAESKey, iv, encryptedContentBytes, nil
}
