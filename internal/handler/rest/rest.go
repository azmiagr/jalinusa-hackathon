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

	public := baseURL.Group("/public")
	public.GET("/dashboard-statistic", r.PublicDashboardStatistic)
	public.GET("/audit-logs", r.GetAuditLog)
	public.GET("/ledgers", r.GetPublicLedger)
	public.GET("/distribution", r.GetPublicDistributionCount)

	auth := baseURL.Group("/auth")
	auth.POST("/login", r.Login)

	binding := baseURL.Group("/bind")
	binding.GET("/posts", r.GetAllPosts)
	binding.GET("/posts/:postID", r.GetPost)
	binding.POST("/posts", r.BindingDevice)

	resources := baseURL.Group("/resources")
	resources.POST("/request/:postID", r.CreateResource)
	resources.POST("/confirm", r.ConfirmResource)

	admin := baseURL.Group("/admin")
	admin.Use(r.middleware.AuthenticateUser)
	admin.GET("/posts", r.GetAllPosts)
	admin.GET("/resources", r.GetResourceList)
	admin.GET("/resources/:ledgerID", r.GetResourceDetails)
	admin.GET("/resources/statistics", r.GetRequestStatistic)
	admin.POST("/posts", r.CreatePost)
	admin.PATCH("/resources/:ledgerID", r.UpdateResourceStatus)

}

func (r *Rest) Run() {
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")

	r.router.Run(fmt.Sprintf("%s:%s", addr, port))
}
