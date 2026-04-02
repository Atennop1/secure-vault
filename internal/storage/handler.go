package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/Atennop1/secure-vault/proto/storagepb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	storagepb.UnimplementedStorageServiceServer

	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Store(ctx context.Context, req *storagepb.StoreRequest) (*emptypb.Empty, error) {
	err := h.service.Store(ctx, req.Key, req.Value)
	if err != nil {
		return &emptypb.Empty{}, fmt.Errorf("storage: failed to store: %w", err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) Load(ctx context.Context, req *storagepb.LoadRequest) (*storagepb.LoadResponse, error) {
	value, err := h.service.Load(ctx, req.Key)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, status.Error(codes.NotFound, "key not found")
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &storagepb.LoadResponse{Value: value}, nil
}
