package generator

import (
	"context"

	"github.com/Atennop1/secure-vault/proto/generatorpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	generatorpb.UnimplementedGeneratorServiceServer

	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Generate(ctx context.Context, req *generatorpb.GenerateRequest) (*generatorpb.GenerateResponse, error) {
	if req.Length <= 0 {
		return nil, status.Error(codes.InvalidArgument, "length must be more than 0")
	}

	slug := h.service.Generate(int(req.Length))

	return &generatorpb.GenerateResponse{
		Slug: slug,
	}, nil
}
