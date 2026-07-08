package rest

import (
	"net/http"

	"github.com/azmiagr/jalinusa-hackathon/model"
	"github.com/azmiagr/jalinusa-hackathon/pkg/helper"
	"github.com/azmiagr/jalinusa-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
)

func (r *Rest) CreatePost(c *gin.Context) {
	var param model.CreatePostRequest
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	userID := helper.GetAuthenticatedUserID(c)

	result, err := r.service.PostService.CreatePost(param, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "success to create post", result)

}

func (r *Rest) GetAllPosts(c *gin.Context) {
	result, err := r.service.PostService.GetAllPosts()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get all posts", result)
}

func (r *Rest) BindingDevice(c *gin.Context) {
	var param model.BindPostRequest
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	result, err := r.service.PostService.BindingDevice(param)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to bind device", result)
}
