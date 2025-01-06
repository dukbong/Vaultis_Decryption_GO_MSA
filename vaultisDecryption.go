package main

import (
	"fmt"
	"log"
	"net"
	"vaultis-go-module/module/decryption"

	"github.com/gin-gonic/gin"
	eureka "github.com/xuanbo/eureka-client"
)

func main() {
	// 사용 가능한 포트 동적 할당
	listener, err := net.Listen("tcp", ":0") // 포트 0으로 바인딩하면 OS가 사용 가능한 포트를 자동으로 할당
	if err != nil {
		log.Fatalf("포트 바인딩 오류: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	listener.Close()

	// Eureka 클라이언트 설정
	client := eureka.NewClient(&eureka.Config{
		DefaultZone:           "http://127.0.0.1:8761/eureka/",
		App:                   "decryption-go",
		Port:                  port, // 동적으로 할당된 포트 사용
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
	log.Printf("서버가 포트 %d에서 실행 중입니다.", port)
	r.Run(fmt.Sprintf(":%d", port))
}
