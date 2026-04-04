package migration

import (
	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"gorm.io/gorm"
)

// Migrate runs auto-migration and seeds initial data.
func Migrate(db *gorm.DB) error {
	if err := AutoMigrate(db); err != nil {
		return err
	}
	return SeedData(db)
}

// AutoMigrate runs GORM auto-migration for all models.
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// User related
		&model.User{},
		&model.LoginAttempt{},

		// Team related
		&model.Team{},
		&model.TeamMember{},
		&model.TeamInvitation{},

		// Problem related
		&model.Problem{},
		&model.TestData{},
		&model.Sample{},
		&model.Tag{},
		&model.ProblemTag{},

		// Submission
		&model.Submission{},

		// Contest related
		&model.Contest{},
		&model.ContestProblem{},
		&model.ContestSignup{},
		&model.DIYContestTemplate{},

		// Announcement
		&model.Announcement{},

		// Anti-cheat
		&model.CheatRecord{},

		// AI records
		&model.AIProblemRecord{},
		&model.AITestdataRecord{},

		// Import
		&model.ImportRecord{},

		// CTF
		&model.CTFResource{},

		// System
		&model.SystemConfig{},
	)
}
