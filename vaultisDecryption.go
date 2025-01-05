package main

import (
	"vaultis-go-module/module/decryption"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/decryption", decryption.Decrypt)
	r.Run(":8080")
}
