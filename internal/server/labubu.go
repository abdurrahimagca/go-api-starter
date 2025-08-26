package server

import (
	"context"

	"github.com/abdurrahimagca/go-api-starter/internal/api"
	"github.com/abdurrahimagca/go-api-starter/internal/labubu"
)

// CreateLabubu implements the POST /labubu endpoint
func (s *Server) CreateLabubu(ctx context.Context, request api.CreateLabubuRequestObject) (api.CreateLabubuResponseObject, error) {
	req := labubu.CreateLabubuRequest{
		Text: request.Body.Text,
	}

	result, err := s.labubuService.CreateLabubu(ctx, req)
	if err != nil {
		return api.CreateLabubu400Response{}, nil
	}

	return api.CreateLabubu200JSONResponse{
		Id:   result.ID,
		Text: result.Text,
	}, nil
}

// GetLabubu implements the GET /labubu endpoint
func (s *Server) GetLabubu(ctx context.Context, request api.GetLabubuRequestObject) (api.GetLabubuResponseObject, error) {
	results, err := s.labubuService.GetAllLabubu(ctx)
	if err != nil {
		return nil, err
	}

	var labubuItems api.GetLabubu200JSONResponse
	
	for _, item := range results {
		labubuItems = append(labubuItems, struct {
			Id   int    `json:"id"`
			Text string `json:"text"`
		}{
			Id:   item.ID,
			Text: item.Text,
		})
	}

	return labubuItems, nil
}