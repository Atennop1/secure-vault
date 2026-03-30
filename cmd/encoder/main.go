package main

import (
	"fmt"

	"github.com/Atennop1/secure-vault/internal/encoder"
	"github.com/Atennop1/secure-vault/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	config.AddPaths("config")
	err := config.LoadEnv("secret")
	if err != nil {
		panic(fmt.Errorf("cmd: failed to load config/secret.env: %w", err))
	}

	err = config.LoadEnv("ports")
	if err != nil {
		panic(fmt.Errorf("cmd: failed to load config/ports.env: %w", err))
	}

	service, err := encoder.NewService([]byte(viper.GetString("AES256_SECRET")), viper.GetInt("GENERATOR_PORT"), viper.GetInt("STORAGE_PORT"))
	if err != nil {
		panic(fmt.Errorf("cmd: failed to create service: %w", err))
	}

	r := gin.Default()

	handler := encoder.NewHandler(service)
	r.POST("/encode", handler.Encode)

	err = r.Run(fmt.Sprintf("localhost:%d", viper.GetInt("ENCODER_PORT")))
	if err != nil {
		panic(fmt.Errorf("cmd: failed to run gin engine on port %d: %w", viper.GetInt("ENCODER_PORT"), err))
	}
}
