package repository

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"gorm.io/gorm"
)

type AnnouncementRepository interface {
	Create(ctx context.Context, announcement *model.Announcement) error
	GetByID(ctx context.Context, id uint) (*model.Announcement, error)
	Update(ctx context.Context, announcement *model.Announcement) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]model.Announcement, int64, error)
}

type announcementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) AnnouncementRepository {
	return &announcementRepository{db: db}
}

func (r *announcementRepository) Create(ctx context.Context, announcement *model.Announcement) error {
	return r.db.WithContext(ctx).Create(announcement).Error
}

func (r *announcementRepository) GetByID(ctx context.Context, id uint) (*model.Announcement, error) {
	var announcement model.Announcement
	if err := r.db.WithContext(ctx).Preload("Creator").First(&announcement, id).Error; err != nil {
		return nil, err
	}
	return &announcement, nil
}

func (r *announcementRepository) Update(ctx context.Context, announcement *model.Announcement) error {
	return r.db.WithContext(ctx).Save(announcement).Error
}

func (r *announcementRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Announcement{}, id).Error
}

func (r *announcementRepository) List(ctx context.Context, offset, limit int) ([]model.Announcement, int64, error) {
	var announcements []model.Announcement
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.Announcement{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Preload("Creator").
		Order("is_pinned DESC, created_at DESC").
		Offset(offset).Limit(limit).
		Find(&announcements).Error; err != nil {
		return nil, 0, err
	}

	return announcements, total, nil
}
