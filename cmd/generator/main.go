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
	err := config.LoadEnv("ports", "config")
	if err != nil {
		panic(err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("GENERATOR_PORT")))
	if err != nil {
		panic(err)
	}

	repo := generator.NewRepository()
	service := generator.NewService(repo)
	handler := generator.NewHandler(service)

	server := grpc.NewServer()
	generatorpb.RegisterGeneratorServiceServer(server, handler)

	err = server.Serve(l)
	if err != nil {
		panic(err)
	}
}
