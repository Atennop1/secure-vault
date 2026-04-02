package main

import (
	"fmt"
	"net"

	"github.com/Atennop1/secure-vault/internal/storage"
	"github.com/Atennop1/secure-vault/pkg/config"
	"github.com/Atennop1/secure-vault/proto/storagepb"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
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

	redisOpts, err := redis.ParseURL("redis://vault-redis:6379/0")
	if err != nil {
		panic(fmt.Errorf("cmd: failed to create a redis connection on port 6379: %w", err))
	}

	repo := storage.NewRepository(redis.NewClient(redisOpts))
	service := storage.NewService(repo)
	handler := storage.NewHandler(service)

	server := grpc.NewServer()
	storagepb.RegisterStorageServiceServer(server, handler)

	err = server.Serve(l)
	if err != nil {
		panic(fmt.Errorf("cmd: failed to serve: %w", err))
	}
}
