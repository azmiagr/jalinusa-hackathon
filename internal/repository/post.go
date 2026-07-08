package repository

import (
	"errors"

	"github.com/azmiagr/jalinusa-hackathon/entity"
	"github.com/azmiagr/jalinusa-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPostRepository interface {
	CreatePost(tx *gorm.DB, post *entity.Post) error
	GetPost(tx *gorm.DB, param model.GetPost) (*entity.Post, error)
	GetLastPost(tx *gorm.DB) (*entity.Post, error)
	GetAllPost(tx *gorm.DB) ([]*entity.Post, error)
	CountKiosk(tx *gorm.DB, postID uuid.UUID) (int64, error)
}

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) IPostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) CreatePost(tx *gorm.DB, post *entity.Post) error {
	err := tx.Debug().Create(post).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *PostRepository) GetPost(tx *gorm.DB, param model.GetPost) (*entity.Post, error) {
	var post entity.Post
	err := tx.Debug().Where(&param).First(&post).Error
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetAllPost(tx *gorm.DB) ([]*entity.Post, error) {
	var posts []*entity.Post
	err := tx.Debug().Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetLastPost(tx *gorm.DB) (*entity.Post, error) {
	var posko entity.Post

	err := tx.
		Order("post_code DESC").
		First(&posko).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &posko, nil
}

func (r *PostRepository) CountKiosk(tx *gorm.DB, postID uuid.UUID) (int64, error) {
	var count int64
	err := tx.Debug().Model(&entity.DeviceBinding{}).Where("post_id = ?", postID).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
