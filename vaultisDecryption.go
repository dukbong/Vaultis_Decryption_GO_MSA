package main

import (
	"vaultis-go-module/module/decryption"

	"github.com/gin-gonic/gin"
	eureka "github.com/xuanbo/eureka-client"
)

func main() {
	client := eureka.NewClient(&eureka.Config{
		DefaultZone:           "http://127.0.0.1:8761/eureka/",
		App:                   "decryption-go",
		Port:                  8999,
		RenewalIntervalInSecs: 10,
		DurationInSecs:        30,
		Metadata: map[string]interface{}{
			"VERSION":              "0.1.0",
			"NODE_GROUP_ID":        0,
			"PRODUCT_CODE":         "DEFAULT",
			"PRODUCT_VERSION_CODE": "DEFAULT",
			"PRODUCT_ENV_CODE":     "DEFAULT",
			"SERVICE_VERSION_CODE": "DEFAULT",
		},
	})

	client.Start()

	// Gin 라우터 설정
	r := gin.Default()
	r.POST("/decryption", decryption.Decrypt)
	r.Run(":8999") // 포트 8999에서 HTTP 서버 실행
}
