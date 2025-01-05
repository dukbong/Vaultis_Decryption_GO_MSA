package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
	"vaultis-go-module/module/decryption"

	"github.com/ArthurHlt/go-eureka-client/eureka"
	"github.com/gin-gonic/gin"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // 랜덤 시드 설정

	// Eureka 클라이언트 생성
	client := eureka.NewClient([]string{
		"http://127.0.0.1:8761/eureka", // Spring Boot 기반 Eureka 서버
	})

	// 랜덤 instance_id 생성
	applicationName := "decryption-go"
	randomValue := fmt.Sprintf("%d", rand.Intn(1000000)) // 0~999999 범위의 랜덤 숫자
	instanceID := fmt.Sprintf("%s:%s", applicationName, randomValue)

	// 인스턴스 정보 설정
	instance := eureka.NewInstanceInfo(
		instanceID,      // 인스턴스 ID (고유한 값 사용)
		applicationName, // 애플리케이션 이름
		"localhost",     // 실제 IP 주소나 고유한 호스트 이름 사용
		8080,            // 포트 번호
		30,              // Heartbeat 주기 (초 단위)
		false,           // 디버그 여부
	)

	// 메타데이터 초기화
	instance.Metadata = &eureka.MetaData{
		Map: make(map[string]string),
	}

	// Eureka에 인스턴스 등록
	err := client.RegisterInstance(applicationName, instance)
	if err != nil {
		log.Fatalf("Eureka에 인스턴스를 등록하는 중 오류 발생: %v", err)
	}
	log.Printf("Eureka에 인스턴스 등록 성공: %s", instanceID)

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
