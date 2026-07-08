package middleware

import (
	"github.com/azmiagr/jalinusa-hackathon/internal/service"
	"github.com/azmiagr/jalinusa-hackathon/pkg/jwt"
)

type Interface interface {
}

type middleware struct {
	service *service.Service
	jwtAuth jwt.Interface
}

func Init(service *service.Service, jwtAuth jwt.Interface) Interface {
	return &middleware{
		service: service,
		jwtAuth: jwtAuth,
	}
}
