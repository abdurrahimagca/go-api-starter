package server

import (
	"context"

	"github.com/abdurrahimagca/go-api-starter/internal/api"
)

// Login implements the /login endpoint
func (s *Server) Login(ctx context.Context, request api.LoginRequestObject) (api.LoginResponseObject, error) {
	loginResponse, err := s.authService.Login(ctx)
	if err != nil {
		return nil, err
	}

	return api.Login200JSONResponse{
		AccessToken:  loginResponse.AccessToken,
		RefreshToken: loginResponse.RefreshToken,
	}, nil
}