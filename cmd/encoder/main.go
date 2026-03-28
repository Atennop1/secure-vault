package main

import (
	"fmt"

	"github.com/Atennop1/secure-vault/internal/encoder"
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

	service, err := encoder.NewService(secret, viper.GetInt("GENERATOR_PORT"), viper.GetInt("STORAGE_PORT"))
	if err != nil {
		panic(err)
	}

	handler := encoder.NewHandler(service)

	r := gin.Default()
	r.POST("/encode", handler.Encode)

	err = r.Run(fmt.Sprintf("localhost:%d", 8080))
	if err != nil {
		panic(err)
	}
}
