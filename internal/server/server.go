package server

import (
	"github.com/abdurrahimagca/go-api-starter/internal/api"
	"github.com/abdurrahimagca/go-api-starter/internal/auth"
	"github.com/abdurrahimagca/go-api-starter/internal/labubu"
)

type Server struct {
	authService   auth.Service
	labubuService labubu.Service
}

func NewServer(authService auth.Service, labubuService labubu.Service) *Server {
	return &Server{
		authService:   authService,
		labubuService: labubuService,
	}
}

var _ api.StrictServerInterface = (*Server)(nil)