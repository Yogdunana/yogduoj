package service

import (
	"context"
	"errors"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrAnnouncementNotFound = errors.New("announcement not found")
)

// CreateAnnouncementRequest is the request body for creating an announcement.
type CreateAnnouncementRequest struct {
	Title   string `json:"title" binding:"required,max=256"`
	Content string `json:"content" binding:"required"`
	IsPinned bool  `json:"is_pinned"`
}

// UpdateAnnouncementRequest is the request body for updating an announcement.
type UpdateAnnouncementRequest struct {
	Title   *string `json:"title" binding:"omitempty,max=256"`
	Content *string `json:"content"`
	IsPinned *bool  `json:"is_pinned"`
}

// AnnouncementService handles announcement-related business logic.
type AnnouncementService interface {
	ListAnnouncements(ctx context.Context, offset, limit int) ([]model.Announcement, int64, error)
	GetAnnouncement(ctx context.Context, id uint) (*model.Announcement, error)
	CreateAnnouncement(ctx context.Context, adminID uint, req CreateAnnouncementRequest) (*model.Announcement, error)
	UpdateAnnouncement(ctx context.Context, id uint, req UpdateAnnouncementRequest) (*model.Announcement, error)
	DeleteAnnouncement(ctx context.Context, id uint) error
}

type announcementService struct {
	announcementRepo repository.AnnouncementRepository
}

func NewAnnouncementService(announcementRepo repository.AnnouncementRepository) AnnouncementService {
	return &announcementService{announcementRepo: announcementRepo}
}

// ListAnnouncements returns a paginated list of announcements, pinned first.
func (s *announcementService) ListAnnouncements(ctx context.Context, offset, limit int) ([]model.Announcement, int64, error) {
	return s.announcementRepo.List(ctx, offset, limit)
}

// GetAnnouncement returns an announcement by ID.
func (s *announcementService) GetAnnouncement(ctx context.Context, id uint) (*model.Announcement, error) {
	announcement, err := s.announcementRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAnnouncementNotFound
		}
		return nil, err
	}
	return announcement, nil
}

// CreateAnnouncement creates a new announcement.
func (s *announcementService) CreateAnnouncement(ctx context.Context, adminID uint, req CreateAnnouncementRequest) (*model.Announcement, error) {
	announcement := &model.Announcement{
		Title:     req.Title,
		Content:   req.Content,
		IsPinned:  req.IsPinned,
		CreatedBy: adminID,
	}

	if err := s.announcementRepo.Create(ctx, announcement); err != nil {
		return nil, err
	}

	return announcement, nil
}

// UpdateAnnouncement updates an existing announcement.
func (s *announcementService) UpdateAnnouncement(ctx context.Context, id uint, req UpdateAnnouncementRequest) (*model.Announcement, error) {
	announcement, err := s.announcementRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAnnouncementNotFound
		}
		return nil, err
	}

	if req.Title != nil {
		announcement.Title = *req.Title
	}
	if req.Content != nil {
		announcement.Content = *req.Content
	}
	if req.IsPinned != nil {
		announcement.IsPinned = *req.IsPinned
	}

	if err := s.announcementRepo.Update(ctx, announcement); err != nil {
		return nil, err
	}

	return announcement, nil
}

// DeleteAnnouncement deletes an announcement.
func (s *announcementService) DeleteAnnouncement(ctx context.Context, id uint) error {
	_, err := s.announcementRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAnnouncementNotFound
		}
		return err
	}

	return s.announcementRepo.Delete(ctx, id)
}
