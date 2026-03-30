package main

import (
	"fmt"
	"net"

	"github.com/Atennop1/secure-vault/internal/generator"
	"github.com/Atennop1/secure-vault/pkg/config"
	"github.com/Atennop1/secure-vault/proto/generatorpb"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	config.AddPaths("config")
	err := config.LoadEnv("ports")
	if err != nil {
		panic(fmt.Errorf("cmd: failed to load config/ports.env: %w", err))
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("GENERATOR_PORT")))
	if err != nil {
		panic(fmt.Errorf("cmd: failed to create a listener on port %d: %w", viper.GetInt("GENERATOR_PORT"), err))
	}

	repo := generator.NewRepository()
	service := generator.NewService(repo)
	handler := generator.NewHandler(service)

	server := grpc.NewServer()
	generatorpb.RegisterGeneratorServiceServer(server, handler)

	err = server.Serve(l)
	if err != nil {
		panic(fmt.Errorf("cmd: failed to serve: %w", err))
	}
}
