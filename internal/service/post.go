package service

import (
	"errors"
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

type IPostService interface {
	CreatePost(param model.CreatePostRequest, userID uuid.UUID) (*model.CreatePostResponse, error)
	GetAllPosts() (*model.GetAllPosts, error)
	BindingDevice(param model.BindPostRequest) (*model.BindPostResponse, error)
	GetPost(postID uuid.UUID) (*model.GetPostResponse, error)
}

type PostService struct {
	db         *gorm.DB
	postRepo   repository.IPostRepository
	userRepo   repository.IUserRepository
	deviceRepo repository.IDeviceRepository
}

func NewPostService(postRepo repository.IPostRepository, userRepo repository.IUserRepository, deviceRepo repository.IDeviceRepository) IPostService {
	return &PostService{
		db:         mariadb.Connection,
		postRepo:   postRepo,
		userRepo:   userRepo,
		deviceRepo: deviceRepo,
	}
}

func (s *PostService) CreatePost(param model.CreatePostRequest, userID uuid.UUID) (*model.CreatePostResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	post, err := s.postRepo.GetPost(tx, model.GetPost{
		PostName: param.PostName,
	})
	if post != nil {
		return nil, apperrors.BadRequest("post name already exist")
	}

	if err != gorm.ErrRecordNotFound {
		return nil, apperrors.InternalServer("failed to get post")
	}

	lastPost, err := s.postRepo.GetLastPost(tx)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get last post")
	}

	var code string
	if lastPost == nil {
		code = "PSK00001"
	} else {
		numStr := strings.TrimPrefix(lastPost.PostCode, "PSK")
		number, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, apperrors.InternalServer("failed to generate code")
		}

		code = fmt.Sprintf("PSK%05d", number+1)
	}

	newPost := &entity.Post{
		PostID:     uuid.New(),
		UserID:     userID,
		PostCode:   code,
		PostName:   param.PostName,
		Coordinate: param.Coordinate,
		Capacity:   param.Capacity,
		Status:     "active",
	}

	err = s.postRepo.CreatePost(tx, newPost)
	if err != nil {
		return nil, apperrors.InternalServer("failed to create post")
	}

	user, err := s.userRepo.GetUser(tx, model.GetUserParam{
		UserID: userID,
	})
	if err != nil {
		return nil, apperrors.InternalServer("failed to get user")
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, apperrors.InternalServer("failed to commit transaction")
	}

	return &model.CreatePostResponse{
		PostID:     newPost.PostID,
		CreatedBy:  user.Username,
		PostCode:   newPost.PostCode,
		PostName:   newPost.PostName,
		Coordinate: newPost.Coordinate,
		Capacity:   newPost.Capacity,
		Status:     newPost.Status,
	}, nil

}

func (s *PostService) GetAllPosts() (*model.GetAllPosts, error) {
	var post []model.PostsResponse

	posts, err := s.postRepo.GetAllPost(s.db)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get posts")
	}

	for _, p := range posts {

		count, err := s.postRepo.CountKiosk(s.db, p.PostID)
		if err != nil {
			return nil, apperrors.InternalServer("failed to count kiosk")
		}

		post = append(post, model.PostsResponse{
			PostID:     p.PostID,
			PostCode:   p.PostCode,
			PostName:   p.PostName,
			KioskCount: int(count),
			Capacity:   p.Capacity,
		})
	}

	return &model.GetAllPosts{
		Posts: post,
	}, nil
}

func (s *PostService) BindingDevice(param model.BindPostRequest) (*model.BindPostResponse, error) {
	tx := s.db.Begin()
	defer tx.Rollback()

	user, err := s.userRepo.GetUser(tx, model.GetUserParam{
		Username: param.Username,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.BadRequest("invalid username or password")
		}
		return nil, apperrors.InternalServer("failed to get user")
	}

	post, err := s.postRepo.GetPost(tx, model.GetPost{
		PostID: param.PostID,
	})
	if err != nil {
		return nil, apperrors.InternalServer("failed to get post")
	}

	device := &entity.Device{
		DeviceID:   uuid.New(),
		DeviceName: param.DeviceName,
	}

	deviceBinding := &entity.DeviceBinding{
		DeviceBindingID: uuid.New(),
		DeviceID:        device.DeviceID,
		PostID:          post.PostID,
		Status:          "success",
		BoundBy:         user.UserID,
	}

	_, err = s.deviceRepo.FindOrCreateDevice(tx, device)
	if err != nil {
		return nil, apperrors.InternalServer("failed to create device")
	}

	err = s.deviceRepo.CreateDeviceBinding(tx, deviceBinding)
	if err != nil {
		return nil, apperrors.InternalServer("failed to bind device")
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, apperrors.InternalServer("failed to commit tx")
	}

	return &model.BindPostResponse{
		PostID:     post.PostID,
		Status:     deviceBinding.Status,
		DeviceName: device.DeviceName,
		BoundBy:    user.Username,
		BountAt:    deviceBinding.BoundAt,
	}, nil

}

func (s *PostService) GetPost(postID uuid.UUID) (*model.GetPostResponse, error) {
	post, err := s.postRepo.GetPost(s.db, model.GetPost{
		PostID: postID,
	})
	if err != nil {
		return nil, apperrors.InternalServer("failed to get post detail")
	}

	return &model.GetPostResponse{
		PostID:   post.PostID,
		PostCode: post.PostCode,
		PostName: post.PostName,
	}, nil
}
