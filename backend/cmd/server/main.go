package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Yogdunana/yogduoj/backend/internal/config"
	"github.com/Yogdunana/yogduoj/backend/internal/handler"
	"github.com/Yogdunana/yogduoj/backend/internal/middleware"
	"github.com/Yogdunana/yogduoj/backend/internal/migration"
	"github.com/Yogdunana/yogduoj/backend/internal/pkg/jwt"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
	"github.com/Yogdunana/yogduoj/backend/internal/router"
	"github.com/Yogdunana/yogduoj/backend/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Determine config path
	configPath := "config.yaml"
	if envPath := os.Getenv("YOGDUOJ_CONFIG_PATH"); envPath != "" {
		configPath = envPath
	}

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	zapLogger := initLogger(cfg)
	defer zapLogger.Sync()

	zapLogger.Info("Starting YogduOJ Backend",
		zap.String("port", fmt.Sprintf("%d", cfg.Server.Port)),
		zap.String("mode", cfg.Server.Mode),
	)

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database
	db, err := initDatabase(cfg)
	if err != nil {
		zapLogger.Fatal("Failed to connect to database", zap.Error(err))
	}
	zapLogger.Info("Database connected successfully")

	// Run migrations
	if err := migration.AutoMigrate(db); err != nil {
		zapLogger.Fatal("Failed to run migrations", zap.Error(err))
	}
	zapLogger.Info("Database migrations completed")

	// Initialize JWT manager
	jwtManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	teamRepo := repository.NewTeamRepository(db)
	problemRepo := repository.NewProblemRepository(db)
	submissionRepo := repository.NewSubmissionRepository(db)
	contestRepo := repository.NewContestRepository(db)
	announcementRepo := repository.NewAnnouncementRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtManager)
	userService := service.NewUserService(userRepo)
	teamService := service.NewTeamService(teamRepo, userRepo)
	contestService := service.NewContestService(contestRepo, submissionRepo)
	judgeService := service.NewJudgeService(cfg.Judge.GRPCAddr, cfg.Judge.Timeout, problemRepo, submissionRepo)

	// Wire up the judge service with the WebSocket notification callback.
	// The handler package registers the callback via init(), so we retrieve it here.
	if js, ok := judgeService.(interface{ SetNotifyFunc(service.JudgeNotifyFunc) }); ok {
		js.SetNotifyFunc(service.GetGlobalNotifyFunc())
	}
	announcementService := service.NewAnnouncementService(announcementRepo)
	antiCheatService := service.NewAntiCheatService()
	aiService := service.NewAIService()
	importService := service.NewImportService()
	systemService := service.NewSystemService()
	problemService := service.NewProblemService(problemRepo, submissionRepo)
	submissionService := service.NewSubmissionService(submissionRepo, problemRepo, judgeService)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	teamHandler := handler.NewTeamHandler(teamService)
	problemHandler := handler.NewProblemHandler(problemService)
	submissionHandler := handler.NewSubmissionHandler(submissionService, judgeService)
	contestHandler := handler.NewContestHandler(contestService)
	announcementHandler := handler.NewAnnouncementHandler(announcementService)
	adminHandler := handler.NewAdminHandler(
		userService,
		problemService,
		submissionService,
		contestService,
		announcementService,
		antiCheatService,
		aiService,
		importService,
		systemService,
	)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)
	rateLimiter := middleware.NewRateLimiterMiddleware(&cfg.RateLimit)

	// Initialize Gin engine
	engine := gin.New()

	// Apply global middleware
	engine.Use(middleware.Recovery(zapLogger))
	engine.Use(middleware.Logger(zapLogger))
	engine.Use(middleware.CORS(&cfg.Server))
	engine.Use(rateLimiter.Limit())

	// Setup routes
	r := router.NewRouter(
		engine,
		authMiddleware,
		authHandler,
		userHandler,
		teamHandler,
		problemHandler,
		submissionHandler,
		contestHandler,
		announcementHandler,
		adminHandler,
	)
	r.Setup()

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	zapLogger.Info("Server starting", zap.String("address", addr))

	if err := engine.Run(addr); err != nil {
		zapLogger.Fatal("Failed to start server", zap.Error(err))
	}
}

func initLogger(cfg *config.Config) *zap.Logger {
	var zapLevel zapcore.Level
	switch cfg.Log.Level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	// Add file logging if configured
	if cfg.Log.FilePath != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.Log.FilePath), 0755); err == nil {
			file, err := os.OpenFile(cfg.Log.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				core = zapcore.NewTee(
					core,
					zapcore.NewCore(
						zapcore.NewJSONEncoder(encoderConfig),
						zapcore.AddSync(file),
						zapLevel,
					),
				)
			}
		}
	}

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Second)

	return db, nil
}
