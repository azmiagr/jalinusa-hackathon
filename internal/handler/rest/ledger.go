package rest

import (
	"net/http"

	"github.com/azmiagr/jalinusa-hackathon/model"
	"github.com/azmiagr/jalinusa-hackathon/pkg/helper"
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

func (r *Rest) GetResourceList(c *gin.Context) {
	result, err := r.service.LedgerService.GetResourcesList()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get resouces list", result)
}

func (r *Rest) GetResourceDetails(c *gin.Context) {
	ledgerID, err := uuid.Parse(c.Param("ledgerID"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid ledger ID", err)
		return
	}

	result, err := r.service.LedgerService.GetResourceDetail(ledgerID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get resource details", result)

}

func (r *Rest) UpdateResourceStatus(c *gin.Context) {
	userID := helper.GetAuthenticatedUserID(c)

	ledgerID, err := uuid.Parse(c.Param("ledgerID"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid ledger ID", err)
		return
	}

	err = r.service.LedgerService.UpdateResourceStatus(ledgerID, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to update status", nil)

}

func (r *Rest) GetRequestStatistic(c *gin.Context) {
	result, err := r.service.LedgerService.GetRequestStatistic()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get request statistic", result)
}

func (r *Rest) PublicDashboardStatistic(c *gin.Context) {
	result, err := r.service.LedgerService.PublicDashboardStatistic()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get public dashboard statistic", result)
}

func (r *Rest) GetAuditLog(c *gin.Context) {
	result, err := r.service.LedgerService.GetAuditLog()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get audit logs", result)
}

func (r *Rest) GetPublicLedger(c *gin.Context) {
	result, err := r.service.LedgerService.GetPublicLedger()
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success to get public ledger", result)
}
