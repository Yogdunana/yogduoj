package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"gorm.io/gorm"
)

type StatsHandler struct {
	db *gorm.DB
}

func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	var userCount int64
	var problemCount int64
	var submissionCount int64
	var contestCount int64
	var acceptedSubmissionCount int64

	h.db.Model(&model.User{}).Count(&userCount)
	h.db.Model(&model.Problem{}).Where("status = ?", "public").Count(&problemCount)
	h.db.Model(&model.Submission{}).Count(&submissionCount)
	h.db.Model(&model.Contest{}).Count(&contestCount)
	h.db.Model(&model.Submission{}).Where("judge_result = ?", "AC").Count(&acceptedSubmissionCount)

	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total_users":          userCount,
			"total_problems":       problemCount,
			"total_submissions":    submissionCount,
			"total_contests":       contestCount,
			"accepted_submissions": acceptedSubmissionCount,
		},
	})
}
