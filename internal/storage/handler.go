package storage

import (
	"context"

	"github.com/Atennop1/secure-vault/proto/storagepb"
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
	h.service.Store(req.Key, req.Value)
	return &emptypb.Empty{}, nil
}

func (h *Handler) Load(ctx context.Context, req *storagepb.LoadRequest) (*storagepb.LoadResponse, error) {
	value, ok := h.service.Load(req.Key)
	if !ok {
		return &storagepb.LoadResponse{Found: false}, nil
	}

	return &storagepb.LoadResponse{
		Value: value,
		Found: true,
	}, nil
}
