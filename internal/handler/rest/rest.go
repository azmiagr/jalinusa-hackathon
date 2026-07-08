package rest

import (
	"fmt"
	"os"

	"github.com/azmiagr/jalinusa-hackathon/internal/service"
	"github.com/azmiagr/jalinusa-hackathon/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type Rest struct {
	router     *gin.Engine
	service    *service.Service
	middleware middleware.Interface
}

func NewRest(service *service.Service, middleware middleware.Interface) *Rest {
	return &Rest{
		router:     gin.Default(),
		service:    service,
		middleware: middleware,
	}
}

func (r *Rest) MountEndpoint() {
	r.router.Use(r.middleware.Cors())
	baseURL := r.router.Group("/api/v1")

	auth := baseURL.Group("/auth")
	auth.POST("/login", r.Login)

	admin := baseURL.Group("/admin")
	admin.Use(r.middleware.AuthenticateUser)
	admin.GET("/posts", r.GetAllPosts)
	admin.POST("/posts", r.CreatePost)

}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}
