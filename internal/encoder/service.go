package encoder

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/Atennop1/secure-vault/proto/generatorpb"
	"github.com/Atennop1/secure-vault/proto/storagepb"
)

type Service struct {
	secret    []byte
	generator generatorpb.GeneratorServiceClient
	storage   storagepb.StorageServiceClient
}

func NewService(secret []byte, generator generatorpb.GeneratorServiceClient, storage storagepb.StorageServiceClient) (*Service, error) {
	return &Service{
		secret:    secret,
		generator: generator,
		storage:   storage,
	}, nil
}

func (s *Service) Encode(ctx context.Context, content string) (string, error) {
	block, err := aes.NewCipher(s.secret)
	if err != nil {
		return "", fmt.Errorf("encoder: failed to generate AES block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("encoder: failed to generate GSM block: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("encoder: failed to generate nonce: %w", err)
	}

	generatorResp, err := s.generator.Generate(ctx, &generatorpb.GenerateRequest{Length: 10})
	if err != nil {
		return "", fmt.Errorf("encoder: failed to generate slug: %w", err)
	}

	_, err = s.storage.Store(ctx, &storagepb.StoreRequest{
		Key:   generatorResp.Slug,
		Value: gcm.Seal(nonce, nonce, []byte(content), nil),
	})

	if err != nil {
		return "", fmt.Errorf("encoder: failed to store encoded data: %w", err)
	}

	return generatorResp.Slug, nil
}
