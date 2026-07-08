package rest

import (
	"net/http"

	"github.com/azmiagr/jalinusa-hackathon/model"
	"github.com/azmiagr/jalinusa-hackathon/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Rest) CreateResource(c *gin.Context) {
	postID, err := uuid.Parse(c.Param("postID"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post ID", err)
		return
	}

	var param model.CreateResourceRequest
	err = c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	result, err := r.service.LedgerService.CreateResourceRequest(param, postID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "success to create resource request", result)
}

func (r *Rest) ConfirmResource(c *gin.Context) {
	var param model.ConfirmResource
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "failed to bind input", err)
		return
	}

	result, err := r.service.LedgerService.ConfirmResource(param)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to confirm resource", result)
}
