package service

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
)

type AnnouncementService interface {
	CreateAnnouncement(ctx context.Context, announcement *model.Announcement) error
	GetAnnouncement(ctx context.Context, id uint) (*model.Announcement, error)
	UpdateAnnouncement(ctx context.Context, announcement *model.Announcement) error
	DeleteAnnouncement(ctx context.Context, id uint) error
	ListAnnouncements(ctx context.Context, offset, limit int) ([]model.Announcement, int64, error)
}

type announcementService struct {
	announcementRepo repository.AnnouncementRepository
}

func NewAnnouncementService(announcementRepo repository.AnnouncementRepository) AnnouncementService {
	return &announcementService{announcementRepo: announcementRepo}
}

func (s *announcementService) CreateAnnouncement(ctx context.Context, announcement *model.Announcement) error {
	return s.announcementRepo.Create(ctx, announcement)
}

func (s *announcementService) GetAnnouncement(ctx context.Context, id uint) (*model.Announcement, error) {
	return s.announcementRepo.GetByID(ctx, id)
}

func (s *announcementService) UpdateAnnouncement(ctx context.Context, announcement *model.Announcement) error {
	return s.announcementRepo.Update(ctx, announcement)
}

func (s *announcementService) DeleteAnnouncement(ctx context.Context, id uint) error {
	return s.announcementRepo.Delete(ctx, id)
}

func (s *announcementService) ListAnnouncements(ctx context.Context, offset, limit int) ([]model.Announcement, int64, error) {
	return s.announcementRepo.List(ctx, offset, limit)
}
