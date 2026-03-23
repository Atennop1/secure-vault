package generate

import (
	"math/rand/v2"
	"strings"
)

type Generator struct {
}

func New() *Generator {
	return &Generator{}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (g *Generator) Generate(length int) string {
	var sb strings.Builder
	sb.Grow(length)

	for range length {
		sb.WriteByte(charset[rand.IntN(len(charset))])
	}

	return sb.String()
}
