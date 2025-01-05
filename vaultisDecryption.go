package main

import (
	"log"
	"time"
	"vaultis-go-module/module/decryption"

	"github.com/ArthurHlt/go-eureka-client/eureka"
	"github.com/gin-gonic/gin"
)

func main() {
	// Eureka 클라이언트 생성
	client := eureka.NewClient([]string{
		"http://127.0.0.1:8761/eureka", // Spring Boot 기반 Eureka 서버
	})

	// 인스턴스 정보 설정
	instance := eureka.NewInstanceInfo(
		"JANG",      // 호스트 이름 (고유한 값 사용)
		"JAJANGNG",  // 애플리케이션 이름
		"localhost", // 실제 IP 주소나 고유한 호스트 이름 사용
		8080,        // 포트 번호
		30,          // Heartbeat 주기 (초 단위)
		false,       // 디버그 여부
	)

	// 메타데이터 초기화
	instance.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}

	// Eureka에 인스턴스 등록
	err := client.RegisterInstance("JAJANGNG", instance)
	if err != nil {
		log.Fatalf("Eureka에 인스턴스를 등록하는 중 오류 발생: %v", err)
	}
	log.Println("Eureka에 인스턴스 등록 성공")

	// Eureka 서버가 인스턴스를 반영할 시간을 대기
	time.Sleep(5 * time.Second)

	// Heartbeat 전송
	err = client.SendHeartbeat(instance.App, instance.HostName)
	if err != nil {
		log.Printf("Heartbeat 전송 중 오류 발생: %v", err)
	} else {
		log.Println("Heartbeat 전송 성공")
	}

	// Gin 라우터 설정
	r := gin.Default()
	r.POST("/decryption", decryption.Decrypt)
	r.Run(":8080") // 포트 8080에서 HTTP 서버 실행
}
