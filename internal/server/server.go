package server

import (
	"fmt"

	"github.com/Atennop1/secure-vault/internal/decrypt"
	"github.com/Atennop1/secure-vault/internal/encrypt"
	"github.com/gin-gonic/gin"
)

type Server struct {
	port   int
	engine *gin.Engine
}

func New(port int, encrypter *encrypt.Handler, decrypter *decrypt.Handler) *Server {
	r := gin.Default()

	r.POST("/encrypt", encrypter.Encrypt)
	r.GET("/decrypt/:slug", decrypter.Decrypt)

	return &Server{
		port:   port,
		engine: r,
	}
}

func (s *Server) Run() error {
	err := s.engine.Run(fmt.Sprintf("localhost:%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to run gin server: %w", err)
	}

	return nil
}
