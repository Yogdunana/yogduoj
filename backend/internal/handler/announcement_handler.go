package handler

import (
	"strconv"

	"github.com/Yogdunana/yogduoj/backend/internal/pkg/pagination"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/response"
	"github.com/Yogdunana/yogduoj/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AnnouncementHandler struct {
	announcementService service.AnnouncementService
}

func NewAnnouncementHandler(announcementService service.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{announcementService: announcementService}
}

// ListAnnouncements returns a paginated list of announcements.
// GET /api/v1/announcements?page=1&page_size=20
func (h *AnnouncementHandler) ListAnnouncements(c *gin.Context) {
	p := pagination.GetPagination(c)

	announcements, total, err := h.announcementService.ListAnnouncements(c.Request.Context(), p.Offset(), p.PageSize)
	if err != nil {
		response.InternalError(c, "failed to get announcements")
		return
	}

	response.PaginatedResponse(c, announcements, total, p.Page, p.PageSize)
}

// GetAnnouncement returns an announcement by ID.
// GET /api/v1/announcements/:id
func (h *AnnouncementHandler) GetAnnouncement(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid announcement id")
		return
	}

	announcement, err := h.announcementService.GetAnnouncement(c.Request.Context(), uint(id))
	if err != nil {
		response.NotFound(c, "announcement not found")
		return
	}

	response.Success(c, announcement)
}
