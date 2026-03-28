package main

import (
	"fmt"

	"github.com/Atennop1/secure-vault/internal/decoder"
	"github.com/Atennop1/secure-vault/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	err := config.LoadEnv("secret", "config")
	if err != nil {
		panic(err)
	}

	secret := []byte(viper.GetString("AES256_SECRET"))

	err = config.LoadEnv("ports", "config")
	if err != nil {
		panic(err)
	}

	service, err := decoder.NewService(secret, viper.GetInt("STORAGE_PORT"))
	if err != nil {
		panic(err)
	}

	handler := decoder.NewHandler(service)

	r := gin.Default()
	r.GET("/decode/:slug", handler.Decode)

	err = r.Run(fmt.Sprintf("localhost:%d", 8081))
	if err != nil {
		panic(err)
	}
}
