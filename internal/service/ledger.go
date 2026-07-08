package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/azmiagr/jalinusa-hackathon/entity"
	"github.com/azmiagr/jalinusa-hackathon/internal/repository"
	"github.com/azmiagr/jalinusa-hackathon/model"
	"github.com/azmiagr/jalinusa-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/jalinusa-hackathon/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ILedgerService interface {
	CreateResourceRequest(param model.CreateResourceRequest, postID uuid.UUID) (*model.CreateResourceResponse, error)
	ConfirmResource(param model.ConfirmResource) (*model.ConfirmResourceResponse, error)
	GetResourcesList() (*model.ResourceRequestList, error)
	GetResourceDetail(ledgerID uuid.UUID) (*model.GetResourceDetail, error)
	UpdateResourceStatus(ledgerID, userID uuid.UUID) error
	PublicDashboardStatistic() (*model.PublicDashboard, error)
	GetRequestStatistic() (*model.RequestStatistic, error)
	GetAuditLog() ([]*model.AuditLogResponse, error)
	GetPublicLedger() ([]*model.PublicLedger, error)
}

type LedgerService struct {
	db               *gorm.DB
	ledgerRepo       repository.ILogisticLedgerRepository
	distributionRepo repository.IDistributionRepository
	itemRepo         repository.ILedgerItemRepository
	userRepo         repository.IUserRepository
	postRepo         repository.IPostRepository
	auditRepo        repository.IAuditLogRepository
}

func NewLedgerService(ledgerRepo repository.ILogisticLedgerRepository, distributionRepo repository.IDistributionRepository, itemRepo repository.ILedgerItemRepository, userRepo repository.IUserRepository, postRepo repository.IPostRepository, auditRepo repository.IAuditLogRepository) ILedgerService {
	return &LedgerService{
		db:               mariadb.Connection,
		ledgerRepo:       ledgerRepo,
		distributionRepo: distributionRepo,
		itemRepo:         itemRepo,
		userRepo:         userRepo,
		postRepo:         postRepo,
		auditRepo:        auditRepo,
	}
}

const genesisTransactionHash = "0000000000000000000000000000000000000000000000000000000000000000"

func (s *LedgerService) CreateResourceRequest(param model.CreateResourceRequest, postID uuid.UUID) (*model.CreateResourceResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	post, err := s.postRepo.GetPost(tx, model.GetPost{
		PostID: postID,
	})
	if err != nil {
		return nil, apperrors.InternalServer("failed to get post")
	}

	latestTransaction, err := s.ledgerRepo.GetLatestLedgerForUpdate(tx, post.PostID)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get latest ledger")
	}

	prevHash := genesisTransactionHash

	if latestTransaction != nil {
		prevHash = latestTransaction.CurrentHash
	}

	var blockNumber string
	if latestTransaction == nil {
		blockNumber = "TX00001"
	} else {
		numStr := strings.TrimPrefix(latestTransaction.BlockNumber, "TX")
		number, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, apperrors.InternalServer("failed to generate tx code")
		}

		blockNumber = fmt.Sprintf("TX%05d", number+1)
	}

	ledger := &entity.LogisticLedger{
		LedgerID:    uuid.New(),
		PostID:      post.PostID,
		PrevHash:    prevHash,
		BlockNumber: blockNumber,
	}

	ledger.CurrentHash = buildLedgerHash(ledger)

	err = s.ledgerRepo.CreateLedger(tx, ledger)

	for _, i := range param.Resource {
		item := &entity.LedgerItem{
			LedgerItemID: uuid.New(),
			LedgerID:     ledger.LedgerID,
			Name:         i.Name,
			Quantity:     i.Quantity,
			Unit:         i.Unit,
		}

		err := s.itemRepo.CreateLedgerItem(tx, item)
		if err != nil {
			return nil, apperrors.InternalServer("failed to create ledger items")
		}
	}

	lastDistribution, err := s.distributionRepo.GetLastDistribution(tx)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get latest distribution")
	}

	var code string
	if lastDistribution == nil {
		code = "JLNS00001"
	} else {
		numStr := strings.TrimPrefix(lastDistribution.Code, "JLNS")
		number, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, apperrors.InternalServer("failed to generate distribution code")
		}

		code = fmt.Sprintf("JLNS%05d", number+1)
	}

	distribution := &entity.Distribution{
		DistributionID: uuid.New(),
		LedgerID:       ledger.LedgerID,
		Status:         "diajukan",
		UserID:         nil,
		Code:           code,
	}

	err = s.distributionRepo.CreateDistribution(tx, distribution)
	if err != nil {
		return nil, apperrors.InternalServer("failed to create distribution")
	}

	audit := &entity.AuditLog{
		AuditID: uuid.New(),
		UserID:  nil,
		Action:  "creating resources request",
	}

	_ = s.auditRepo.CreateAuditLog(tx, audit)

	err = tx.Commit().Error
	if err != nil {
		return nil, apperrors.InternalServer("failed to commit tx")
	}

	return &model.CreateResourceResponse{
		LedgerID:         ledger.LedgerID,
		DistributionCode: distribution.Code,
		BlockNumber:      ledger.BlockNumber,
	}, nil
}

func (s *LedgerService) ConfirmResource(param model.ConfirmResource) (*model.ConfirmResourceResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	distribution, err := s.distributionRepo.GetDistribution(tx, model.DistributionParam{
		Code: param.DistributionCode,
	})
	if err != nil {
		return nil, apperrors.BadRequest("distribution not found")
	}

	if distribution.Status != "terdistribusi" {
		return nil, apperrors.BadRequest("items must have been distributed")
	}

	distribution.Status = "diterima"

	err = s.distributionRepo.UpdateDistribution(tx, distribution)
	if err != nil {
		return nil, apperrors.InternalServer("failed to update distribution")
	}

	items, err := s.itemRepo.GetLedgerItemByLedgerID(tx, distribution.LedgerID)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get ledger items")
	}

	var itemsResponse []model.ItemRequest
	for _, i := range items {
		itemsResponse = append(itemsResponse, model.ItemRequest{
			Name:     i.Name,
			Quantity: i.Quantity,
			Unit:     i.Unit,
		})
	}

	audit := &entity.AuditLog{
		AuditID: uuid.New(),
		UserID:  nil,
		Action:  "resource confirmed by user",
	}

	_ = s.auditRepo.CreateAuditLog(tx, audit)

	err = tx.Commit().Error
	if err != nil {
		return nil, apperrors.InternalServer("failed to commit tx")
	}

	return &model.ConfirmResourceResponse{
		Resource: itemsResponse,
	}, nil

}

func (s *LedgerService) GetResourcesList() (*model.ResourceRequestList, error) {
	var (
		resourcesResponse []model.ResourceResponse
		itemsResponse     []model.ItemRequest
	)

	resources, err := s.ledgerRepo.GetResourceListRequest(s.db)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get resources list")
	}

	for _, r := range resources {
		resourcesResponse = append(resourcesResponse, model.ResourceResponse{
			LedgerID: r.LedgerID,
		})

		items, err := s.itemRepo.GetLedgerItemByLedgerID(s.db, r.LedgerID)
		if err != nil {
			return nil, apperrors.InternalServer("failed to get items")
		}

		for _, i := range items {
			itemsResponse = append(itemsResponse, model.ItemRequest{
				Name:     i.Name,
				Quantity: i.Quantity,
				Unit:     i.Unit,
			})
		}

		post, err := s.postRepo.GetPost(s.db, model.GetPost{
			PostID: r.PostID,
		})
		if err != nil {
			return nil, apperrors.InternalServer("failed to get posts")
		}

		distributions, err := s.distributionRepo.GetDistributionsByLedgerID(s.db, r.LedgerID)
		if err != nil {
			return nil, apperrors.InternalServer("failed to get distributions")
		}

		for _, d := range distributions {
			resourcesResponse = append(resourcesResponse, model.ResourceResponse{
				PostName:           post.PostName,
				DistributionCode:   d.Code,
				DistributionStatus: d.Status,
				BlockNumber:        r.BlockNumber,
				Items:              itemsResponse,
			})
		}
	}

	return &model.ResourceRequestList{
		Resources: resourcesResponse,
	}, nil

}

func (s *LedgerService) GetResourceDetail(ledgerID uuid.UUID) (*model.GetResourceDetail, error) {
	var itemsResponse []model.ItemRequest

	items, err := s.itemRepo.GetLedgerItemByLedgerID(s.db, ledgerID)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get items")
	}

	for _, i := range items {
		itemsResponse = append(itemsResponse, model.ItemRequest{
			Name:     i.Name,
			Quantity: i.Quantity,
			Unit:     i.Unit,
		})
	}

	ledger, err := s.ledgerRepo.GetLedger(s.db, model.GetLedgerParam{
		LedgerID: ledgerID,
	})
	if err != nil {
		return nil, apperrors.InternalServer("failed to get ledger")
	}

	return &model.GetResourceDetail{
		Items:      itemsResponse,
		Status:     "valid",
		HashLedger: ledger.CurrentHash,
	}, nil

}

func (s *LedgerService) UpdateResourceStatus(ledgerID, userID uuid.UUID) error {
	tx := s.db.Begin()
	defer tx.Rollback()

	ledger, err := s.ledgerRepo.GetLedger(tx, model.GetLedgerParam{
		LedgerID: ledgerID,
	})
	if err != nil {
		return apperrors.InternalServer("failed to get ledger")
	}

	distribution, err := s.distributionRepo.GetDistribution(tx, model.DistributionParam{
		LedgerID: ledger.LedgerID,
	})
	if err != nil {
		return apperrors.InternalServer("failed to get distribution")
	}

	switch distribution.Status {
	case "diajukan":
		distribution.Status = "diproses"
	case "diproses":
		distribution.Status = "pengiriman"
	case "pengiriman":
		distribution.Status = "terdistribusi"
	case "terdistribusi":
		return apperrors.BadRequest("must be user that confirm the distribution")
	}

	distribution.UserID = &userID

	err = s.distributionRepo.UpdateDistribution(tx, distribution)
	if err != nil {
		return apperrors.InternalServer("failed to update distribution status")
	}

	audit := &entity.AuditLog{
		AuditID: uuid.New(),
		UserID:  &userID,
		Action:  "updated resource status",
	}

	_ = s.auditRepo.CreateAuditLog(tx, audit)

	err = tx.Commit().Error
	if err != nil {
		return apperrors.InternalServer("failed to commit tx")
	}

	return nil

}

func (s *LedgerService) GetRequestStatistic() (*model.RequestStatistic, error) {
	submittedCount, err := s.distributionRepo.GetSubmittedDistribution(s.db)
	if err != nil {
		return nil, apperrors.InternalServer("failed to count submitted")
	}

	deliveredCount, err := s.distributionRepo.GetDeliveredDistribution(s.db)
	if err != nil {
		return nil, apperrors.InternalServer("failed to count delivered")
	}

	acceptedCount, err := s.distributionRepo.GetAcceptedDistribution(s.db)
	if err != nil {
		return nil, apperrors.InternalServer("failed to count accepted")
	}

	unfinishedCount, err := s.distributionRepo.GetUnfinishedDistribution(s.db)
	if err != nil {
		return nil, apperrors.InternalServer("failed to unfinished count")
	}

	totalCount, err := s.distributionRepo.GetAllDistributionCount(s.db)
	if err != nil {
		return nil, apperrors.InternalServer("failed to count total distributed")
	}

	var aidDivertion float64
	aidDivertion = (float64(unfinishedCount) / float64(totalCount)) * 100

	return &model.RequestStatistic{
		Submitted:    int(submittedCount),
		Delivered:    int(deliveredCount),
		Accepted:     int(acceptedCount),
		AidDivertion: roundToTwo(aidDivertion),
	}, nil
}

func (s *LedgerService) PublicDashboardStatistic() (*model.PublicDashboard, error) {
	totalCount, err := s.distributionRepo.GetAllDistributionCount(s.db)
	if err != nil {
		return nil, apperrors.InternalServer("failed to count total distributed")
	}

	unfinishedCount, err := s.distributionRepo.GetUnfinishedDistribution(s.db)
	if err != nil {
		return nil, apperrors.InternalServer("failed to unfinished count")
	}

	totalAccepted, err := s.distributionRepo.GetAcceptedDistribution(s.db)
	if err != nil {
		return nil, apperrors.InternalServer("failed to count total accepted")
	}

	var aidDivertion float64
	aidDivertion = (float64(unfinishedCount) / float64(totalCount)) * 100

	return &model.PublicDashboard{
		TotalRequest:      int(totalCount),
		TotalAccepted:     int(totalAccepted),
		AidDivertionRate:  roundToTwo(aidDivertion),
		HashChainValidity: "valid",
	}, nil

}

func (s *LedgerService) GetAuditLog() ([]*model.AuditLogResponse, error) {
	var response []*model.AuditLogResponse

	audits, err := s.auditRepo.GetAllAuditLog(s.db)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get audit logs")
	}

	for _, a := range audits {
		var userID uuid.UUID
		if a.UserID != nil {
			userID = *a.UserID
		}

		response = append(response, &model.AuditLogResponse{
			AuditID:   a.AuditID,
			UserID:    &userID,
			Action:    a.Action,
			CreatedAt: a.CreatedAt,
		})
	}

	return response, nil

}

func (s *LedgerService) GetPublicLedger() ([]*model.PublicLedger, error) {
	var response []*model.PublicLedger

	ledgers, err := s.ledgerRepo.GetResourceListRequest(s.db)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get ledgers")
	}

	for _, i := range ledgers {
		distribution, err := s.distributionRepo.GetDistribution(s.db, model.DistributionParam{
			LedgerID: i.LedgerID,
		})
		if err != nil {
			return nil, apperrors.InternalServer("failed to get distribution")
		}

		response = append(response, &model.PublicLedger{
			BlockNumber:      i.BlockNumber,
			DistributionCode: distribution.Code,
			CurrentHash:      i.CurrentHash,
			PrevHash:         i.PrevHash,
			Status:           distribution.Status,
		})
	}

	return response, nil

}

func buildLedgerHash(ledger *entity.LogisticLedger) string {
	payload := fmt.Sprintf(
		"%s|%s|%s|%s|%s",
		ledger.PrevHash,
		ledger.LedgerID.String(),
		ledger.PostID.String(),
		ledger.BlockNumber,
		ledger.BlockNumber,
	)

	sum := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(sum[:])
}

func roundToTwo(val float64) float64 {
	return math.Round(val*100) / 100
}
