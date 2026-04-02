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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound = errors.New("there is no ciphered text with such slug")
)

type Service struct {
	secret  []byte
	storage storagepb.StorageServiceClient
}

func NewService(secret []byte, storage storagepb.StorageServiceClient) (*Service, error) {
	return &Service{
		secret:  secret,
		storage: storage,
	}, nil
}

func (s *Service) Decode(ctx context.Context, slug string) (string, error) {
	block, err := aes.NewCipher(s.secret)
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
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return "", fmt.Errorf("decoder: slug %s: %w", slug, ErrNotFound)
		}

		return "", fmt.Errorf("decoder: failed to load ciphered text from storage: %w", err)
	}

	nonceSize := gcm.NonceSize()
	data, err := gcm.Open(nil, resp.Value[:nonceSize], resp.Value[nonceSize:], nil)
	if err != nil {
		return "", fmt.Errorf("decoder: failed to decode ciphered content: %w", err)
	}

	return string(data), nil
}
