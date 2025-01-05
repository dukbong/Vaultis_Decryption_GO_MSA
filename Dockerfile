# 1단계: Go 빌드 환경 설정
FROM golang:1.23 AS build

# 작업 디렉토리 설정
WORKDIR /app

# Go 모듈 파일 복사
COPY go.mod go.sum ./

# Go 모듈 설치
RUN go mod tidy

# 프로젝트 소스 코드 복사
COPY . .

# Go 애플리케이션 빌드
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 2단계: 실행 환경 설정
FROM alpine:latest

# 필요한 라이브러리 설치 (만약 필요하다면)
RUN apk --no-cache add ca-certificates

# 빌드된 Go 애플리케이션 복사
COPY --from=build /app/main /main

# 애플리케이션 실행
CMD ["/main"]