package rest

import (
	"net/http"

	"github.com/azmiagr/jalinusa-hackathon/model"
	"github.com/azmiagr/jalinusa-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *Rest) Login(c *gin.Context) {
	var param model.UserLoginRequest
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	result, err := r.service.UserService.Login(param)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to login user", result)
}
