package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
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
}

type LedgerService struct {
	db               *gorm.DB
	ledgerRepo       repository.ILogisticLedgerRepository
	distributionRepo repository.IDistributionRepository
	itemRepo         repository.ILedgerItemRepository
	userRepo         repository.IUserRepository
	postRepo         repository.IPostRepository
}

func NewLedgerService(ledgerRepo repository.ILogisticLedgerRepository, distributionRepo repository.IDistributionRepository, itemRepo repository.ILedgerItemRepository, userRepo repository.IUserRepository, postRepo repository.IPostRepository) ILedgerService {
	return &LedgerService{
		db:               mariadb.Connection,
		ledgerRepo:       ledgerRepo,
		distributionRepo: distributionRepo,
		itemRepo:         itemRepo,
		userRepo:         userRepo,
		postRepo:         postRepo,
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

	return &model.ConfirmResourceResponse{
		Resource: itemsResponse,
	}, nil

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
