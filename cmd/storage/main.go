package main

import (
	"fmt"
	"github.com/Atennop1/secure-vault/internal/storage"
	"github.com/Atennop1/secure-vault/pkg/config"
	"github.com/Atennop1/secure-vault/proto/storagepb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
)

func main() {
	config.AddPaths("config")
	err := config.LoadEnv("ports")
	if err != nil {
		panic(fmt.Errorf("cmd: failed to load config/ports.env: %w", err))
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("STORAGE_PORT")))
	if err != nil {
		panic(fmt.Errorf("cmd: failed to create a listener on port %d: %w", viper.GetInt("STORAGE_PORT"), err))
	}

	repo := storage.NewRepository()
	service := storage.NewService(repo)
	handler := storage.NewHandler(service)

	server := grpc.NewServer()
	storagepb.RegisterStorageServiceServer(server, handler)

	err = server.Serve(l)
	if err != nil {
		panic(fmt.Errorf("cmd: failed to serve: %w", err))
	}
}
