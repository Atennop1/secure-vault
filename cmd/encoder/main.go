package main

import (
	"fmt"

	"github.com/Atennop1/secure-vault/internal/encoder"
	"github.com/Atennop1/secure-vault/pkg/config"
	"github.com/Atennop1/secure-vault/proto/generatorpb"
	"github.com/Atennop1/secure-vault/proto/storagepb"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	generatorConn, err := grpc.NewClient(fmt.Sprintf("generator:%d", viper.GetInt("GENERATOR_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("cmd: failed to open grpc connection on port %d: %w", viper.GetInt("GENERATOR_PORT"), err))
	}

	storageConn, err := grpc.NewClient(fmt.Sprintf("storage:%d", viper.GetInt("STORAGE_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("cmd: failed to open grpc connection on port %d: %w", viper.GetInt("STORAGE_PORT"), err))
	}

	service, err := encoder.NewService([]byte(viper.GetString("AES256_SECRET")), generatorpb.NewGeneratorServiceClient(generatorConn), storagepb.NewStorageServiceClient(storageConn))
	if err != nil {
		panic(fmt.Errorf("cmd: failed to create service: %w", err))
	}

	r := gin.Default()
	handler := encoder.NewHandler(service)
	r.POST("/encode", handler.Encode)

	err = r.Run(fmt.Sprintf(":%d", viper.GetInt("ENCODER_PORT")))
	if err != nil {
		panic(fmt.Errorf("cmd: failed to run gin engine on port %d: %w", viper.GetInt("ENCODER_PORT"), err))
	}
}
