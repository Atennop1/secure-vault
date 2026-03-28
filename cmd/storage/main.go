package main

import (
	"fmt"
	"net"

	"github.com/Atennop1/secure-vault/internal/storage"
	"github.com/Atennop1/secure-vault/pkg/config"
	"github.com/Atennop1/secure-vault/proto/storagepb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	err := config.LoadEnv("ports", "config")
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("STORAGE_PORT")))
	if err != nil {
		panic(err)
	}

	repo := storage.NewRepository()
	service := storage.NewService(repo)
	handler := storage.NewHandler(service)

	server := grpc.NewServer()
	storagepb.RegisterStorageServiceServer(server, handler)

	err = server.Serve(l)
	if err != nil {
		panic(err)
	}
}
