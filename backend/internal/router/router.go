package router

import (
	"github.com/Yogdunana/yogduoj/backend/internal/handler"
	"github.com/Yogdunana/yogduoj/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine            *gin.Engine
	authMiddleware    *middleware.AuthMiddleware
	authHandler       *handler.AuthHandler
	userHandler       *handler.UserHandler
	teamHandler       *handler.TeamHandler
	problemHandler    *handler.ProblemHandler
	submissionHandler *handler.SubmissionHandler
	contestHandler    *handler.ContestHandler
	announcementHandler *handler.AnnouncementHandler
	adminHandler       *handler.AdminHandler
	statsHandler       *handler.StatsHandler
}

func NewRouter(
	engine *gin.Engine,
	authMiddleware *middleware.AuthMiddleware,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	teamHandler *handler.TeamHandler,
	problemHandler *handler.ProblemHandler,
	submissionHandler *handler.SubmissionHandler,
	contestHandler *handler.ContestHandler,
	announcementHandler *handler.AnnouncementHandler,
	adminHandler *handler.AdminHandler,
	statsHandler *handler.StatsHandler,
) *Router {
	return &Router{
		engine:            engine,
		authMiddleware:    authMiddleware,
		authHandler:       authHandler,
		userHandler:       userHandler,
		teamHandler:       teamHandler,
		problemHandler:    problemHandler,
		submissionHandler: submissionHandler,
		contestHandler:    contestHandler,
		announcementHandler: announcementHandler,
		adminHandler:      adminHandler,
		statsHandler:      statsHandler,
	}
}

func (r *Router) Setup() {
	api := r.engine.Group("/api/v1")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "yogduoj-backend",
		})
	})

	// Stats (public)
	api.GET("/stats", r.statsHandler.GetStats)

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/logout", r.authMiddleware.RequireAuth(), r.authHandler.Logout)
		auth.POST("/refresh", r.authHandler.RefreshToken)
	}

	// User routes (authenticated)
	users := api.Group("/users")
	users.Use(r.authMiddleware.RequireAuth())
	{
		users.GET("/me", r.userHandler.GetMe)
		users.PUT("/me", r.userHandler.UpdateMe)
		users.PUT("/me/password", r.userHandler.UpdatePassword)
		users.GET("/me/submissions", r.userHandler.GetMySubmissions)
		users.GET("/me/contests", r.userHandler.GetMyContests)
		users.GET("/:id", r.userHandler.GetUser)
	}

	// Team routes
	teams := api.Group("/teams")
	{
		teams.GET("", r.teamHandler.ListTeams)
		teams.GET("/invitations", r.authMiddleware.RequireAuth(), r.teamHandler.GetInvitations)
		teams.POST("/invitations/:invitationId/accept", r.authMiddleware.RequireAuth(), r.teamHandler.AcceptInvitation)
		teams.POST("/invitations/:invitationId/reject", r.authMiddleware.RequireAuth(), r.teamHandler.RejectInvitation)
		teams.GET("/:id", r.teamHandler.GetTeam)
		teams.POST("", r.authMiddleware.RequireAuth(), r.teamHandler.CreateTeam)
		teams.PUT("/:id", r.authMiddleware.RequireAuth(), r.teamHandler.UpdateTeam)
		teams.DELETE("/:id", r.authMiddleware.RequireAuth(), r.teamHandler.DeleteTeam)
		teams.POST("/:id/invite", r.authMiddleware.RequireAuth(), r.teamHandler.InviteUser)
		teams.POST("/:id/leave", r.authMiddleware.RequireAuth(), r.teamHandler.LeaveTeam)
		teams.DELETE("/:id/members/:userId", r.authMiddleware.RequireAuth(), r.teamHandler.RemoveMember)
	}

	// Problem routes
	problems := api.Group("/problems")
	{
		problems.GET("", r.problemHandler.ListProblems)
		problems.GET("/:id", r.problemHandler.GetProblem)
		problems.GET("/:id/samples", r.problemHandler.GetProblemSamples)
		problems.GET("/:id/attachments/:fileId", r.problemHandler.GetAttachment)
	}

	// Submission routes (authenticated)
	submissions := api.Group("/submissions")
	submissions.Use(r.authMiddleware.RequireAuth())
	{
		submissions.POST("", r.submissionHandler.CreateSubmission)
		submissions.GET("", r.submissionHandler.ListSubmissions)
		submissions.GET("/:id", r.submissionHandler.GetSubmission)
		submissions.GET("/:id/code", r.submissionHandler.GetSubmissionCode)
		submissions.GET("/:id/judge", r.submissionHandler.JudgeWebSocket)
	}

	// Contest routes
	contests := api.Group("/contests")
	{
		contests.GET("", r.contestHandler.ListContests)
		contests.GET("/:id", r.contestHandler.GetContest)
		contests.GET("/:id/problems", r.authMiddleware.RequireAuth(), r.contestHandler.GetContestProblems)
		contests.GET("/:id/ranking", r.contestHandler.GetContestRanking)
		contests.GET("/:id/ranking/frozen", r.contestHandler.GetFrozenRanking)
		contests.POST("/:id/signup", r.authMiddleware.RequireAuth(), r.contestHandler.Signup)
		contests.POST("/:id/withdraw", r.authMiddleware.RequireAuth(), r.contestHandler.Withdraw)
		contests.POST("/:id/submissions", r.authMiddleware.RequireAuth(), r.contestHandler.SubmitToContest)
	}

	// Announcement routes (public)
	announcements := api.Group("/announcements")
	{
		announcements.GET("", r.announcementHandler.ListAnnouncements)
		announcements.GET("/:id", r.announcementHandler.GetAnnouncement)
	}

	// CTF routes
	ctf := api.Group("/ctf")
	{
		ctf.GET("/problems", r.problemHandler.ListProblems)
		ctf.GET("/resources", r.adminHandler.ListCTFResources)
		ctf.POST("/submissions", r.authMiddleware.RequireAuth(), r.submissionHandler.CreateSubmission)
	}

	// Admin routes (admin only)
	admin := api.Group("/admin")
	admin.Use(r.authMiddleware.RequireAuth(), r.authMiddleware.RequireAdmin())
	{
		// User management
		admin.GET("/users", r.adminHandler.ListUsers)
		admin.PUT("/users/:id", r.adminHandler.UpdateUser)
		admin.PUT("/users/:id/role", r.adminHandler.UpdateUserRole)
		admin.PUT("/users/:id/disable", r.adminHandler.DisableUser)
		admin.POST("/users/:id/reset-password", r.adminHandler.ResetUserPassword)

		// Problem management
		admin.POST("/problems", r.adminHandler.CreateProblem)
		admin.PUT("/problems/:id", r.adminHandler.UpdateProblem)
		admin.DELETE("/problems/:id", r.adminHandler.DeleteProblem)
		admin.POST("/problems/:id/testdata", r.adminHandler.UploadTestData)
		admin.GET("/problems/:id/testdata", r.adminHandler.GetTestDataList)
		admin.DELETE("/problems/:id/testdata/:dataId", r.adminHandler.DeleteTestData)

		// Submission management
		admin.GET("/submissions", r.adminHandler.ListAllSubmissions)
		admin.POST("/submissions/:id/rejudge", r.adminHandler.RejudgeSubmission)

		// Contest management
		admin.POST("/contests", r.adminHandler.CreateContest)
		admin.PUT("/contests/:id", r.adminHandler.UpdateContest)
		admin.DELETE("/contests/:id", r.adminHandler.DeleteContest)
		admin.PUT("/contests/:id/status", r.adminHandler.UpdateContestStatus)
		admin.POST("/contests/:id/problems", r.adminHandler.AddContestProblem)
		admin.DELETE("/contests/:id/problems/:problemId", r.adminHandler.RemoveContestProblem)
		admin.GET("/contests/:id/signups", r.adminHandler.GetContestSignups)

		// DIY template management
		admin.POST("/diy-templates", r.adminHandler.CreateDIYTemplate)
		admin.GET("/diy-templates", r.adminHandler.ListDIYTemplates)
		admin.PUT("/diy-templates/:id", r.adminHandler.UpdateDIYTemplate)
		admin.DELETE("/diy-templates/:id", r.adminHandler.DeleteDIYTemplate)

		// Announcement management
		admin.POST("/announcements", r.adminHandler.CreateAnnouncement)
		admin.PUT("/announcements/:id", r.adminHandler.UpdateAnnouncement)
		admin.DELETE("/announcements/:id", r.adminHandler.DeleteAnnouncement)

		// Anti-cheat
		admin.POST("/anti-cheat/detect", r.adminHandler.DetectCheating)
		admin.GET("/cheat-records", r.adminHandler.ListCheatRecords)
		admin.PUT("/cheat-records/:id/review", r.adminHandler.ReviewCheatRecord)

		// AI features
		admin.POST("/ai/generate-problem", r.adminHandler.GenerateAIProblem)
		admin.POST("/ai/generate-testdata", r.adminHandler.GenerateAITestdata)

		// Import
		admin.POST("/import/problems", r.adminHandler.ImportProblems)

		// System config
		admin.GET("/system/configs", r.adminHandler.GetSystemConfigs)
		admin.PUT("/system/configs", r.adminHandler.SetSystemConfig)

		// CTF resources
		admin.POST("/ctf/resources", r.adminHandler.CreateCTFResource)
	}
}
