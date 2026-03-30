package decoder

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"github.com/Atennop1/secure-vault/proto/storagepb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	ErrNotFound = errors.New("there is no ciphered text with such slug")
)

type Service struct {
	key []byte

	storage storagepb.StorageServiceClient
}

func NewService(key []byte, storagePort int) (*Service, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", storagePort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("decoder: failed to open grpc connection on port %d: %w", storagePort, err)
	}

	return &Service{
		key:     key,
		storage: storagepb.NewStorageServiceClient(conn),
	}, nil
}

func (s *Service) Decode(ctx context.Context, slug string) (string, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", fmt.Errorf("decoder: failed to generate AES block: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("decoder: failed to generate GSM block: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("decoder: failed to generate nonce: %w", err)
	}

	resp, err := s.storage.Load(ctx, &storagepb.LoadRequest{Key: slug})
	if err != nil {
		return "", fmt.Errorf("decoder: failed to load ciphered text from storage: %w", err)
	}

	if !resp.Found {
		return "", fmt.Errorf("decoder: slug %s: %w", slug, ErrNotFound)
	}

	nonceSize := gcm.NonceSize()
	data, err := gcm.Open(nil, resp.Value[:nonceSize], resp.Value[nonceSize:], nil)
	if err != nil {
		return "", fmt.Errorf("decoder: failed to decode ciphered content: %w", err)
	}

	return string(data), nil
}
