package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StatsHandler struct {
	db *gorm.DB
}

func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	var userCount, problemCount, submissionCount, contestCount, acceptedCount int64

	h.db.Raw("SELECT COUNT(*) FROM users").Scan(&userCount)
	h.db.Raw("SELECT COUNT(*) FROM problems WHERE status = ?", "public").Scan(&problemCount)
	h.db.Raw("SELECT COUNT(*) FROM submissions").Scan(&submissionCount)
	h.db.Raw("SELECT COUNT(*) FROM contests").Scan(&contestCount)
	h.db.Raw("SELECT COUNT(*) FROM submissions WHERE judge_result = ?", "AC").Scan(&acceptedCount)

	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total_users":          userCount,
			"total_problems":       problemCount,
			"total_submissions":    submissionCount,
			"total_contests":       contestCount,
			"accepted_submissions": acceptedCount,
		},
	})
}
