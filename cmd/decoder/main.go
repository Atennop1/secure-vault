package main

import (
	"fmt"

	"github.com/Atennop1/secure-vault/internal/decoder"
	"github.com/Atennop1/secure-vault/pkg/config"
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

	storageConn, err := grpc.NewClient(fmt.Sprintf("storage:%d", viper.GetInt("STORAGE_PORT")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(fmt.Errorf("cmd: failed to open grpc connection on port %d: %w", viper.GetInt("STORAGE_PORT"), err))
	}

	service, err := decoder.NewService([]byte(viper.GetString("AES256_SECRET")), storagepb.NewStorageServiceClient(storageConn))
	if err != nil {
		panic(fmt.Errorf("cmd: failed to create service: %w", err))
	}

	r := gin.Default()
	handler := decoder.NewHandler(service)
	r.GET("/decode/:slug", handler.Decode)

	err = r.Run(fmt.Sprintf(":%d", viper.GetInt("DECODER_PORT")))
	if err != nil {
		panic(fmt.Errorf("cmd: failed to run gin engine on port %d: %w", viper.GetInt("DECODER_PORT"), err))
	}
}
