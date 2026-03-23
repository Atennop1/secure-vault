package main

import (
	"github.com/Atennop1/secure-vault/internal/decrypt"
	"github.com/Atennop1/secure-vault/internal/encrypt"
	"github.com/Atennop1/secure-vault/internal/generate"
	"github.com/Atennop1/secure-vault/internal/repository"
	"github.com/Atennop1/secure-vault/internal/server"
	"github.com/Atennop1/secure-vault/pkg/config"
	"github.com/spf13/viper"
)

func main() {
	err := config.LoadEnv("secret", "config")
	if err != nil {
		panic(err)
	}

	secret := []byte(viper.GetString("AES256_SECRET"))

	repo := repository.New()
	gen := generate.New()

	encryptService := encrypt.NewService(secret, repo, gen)
	encrypter := encrypt.NewHandler(encryptService)

	decryptService := decrypt.NewService(secret, repo)
	decrypter := decrypt.NewHandler(decryptService)

	s := server.New(8080, encrypter, decrypter)
	if err = s.Run(); err != nil {
		panic(err)
	}
}
