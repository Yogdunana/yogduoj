package service

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"sort"
	"time"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"github.com/Yogdunana/yogduoj/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrContestNotFound    = errors.New("contest not found")
	ErrContestNotRunning  = errors.New("contest is not running")
	ErrContestEnded       = errors.New("contest has ended")
	ErrContestFull        = errors.New("contest signup limit reached")
	ErrAlreadySignedUp    = errors.New("already signed up for this contest")
	ErrNotSignedUp        = errors.New("not signed up for this contest")
	ErrInvalidContestTime = errors.New("end time must be after start time")
	ErrInvalidContestStatus = errors.New("invalid contest status")
)

// ContestFilter holds filter parameters for listing contests.
type ContestFilter struct {
	Status   string
	Type     string
	Category string
	RuleType string
	Search   string
	SortBy   string
}

// CreateContestRequest is the request body for creating a contest.
type CreateContestRequest struct {
	Title           string               `json:"title" binding:"required,max=256"`
	ContestType     string               `json:"contest_type" binding:"omitempty,oneof=individual team"`
	Category        string               `json:"category" binding:"omitempty,oneof=programming algorithm ai_assisted ctf"`
	RuleType        string               `json:"rule_type" binding:"omitempty,oneof=acm oi ioi cf ctf awd isw diy"`
	DIYTemplateID   uint                 `json:"diy_template_id"`
	StartTime       time.Time            `json:"start_time" binding:"required"`
	EndTime         time.Time            `json:"end_time" binding:"required"`
	FreezeTime      *time.Time           `json:"freeze_time"`
	Description     string               `json:"description"`
	RuleDescription string               `json:"rule_description"`
	MaxTeamSize     int                  `json:"max_team_size"`
	SignupLimit     int                  `json:"signup_limit"`
	AllowViewOthers bool                 `json:"allow_view_others"`
	ShowRealtimeRank bool                `json:"show_realtime_rank"`
	EnableAIHint    bool                 `json:"enable_ai_hint"`
	DIYRules        string               `json:"diy_rules"`
	Problems        []ContestProblemReq  `json:"problems"`
}

// UpdateContestRequest is the request body for updating a contest.
type UpdateContestRequest struct {
	Title           *string     `json:"title" binding:"omitempty,max=256"`
	ContestType     *string     `json:"contest_type" binding:"omitempty,oneof=individual team"`
	Category        *string     `json:"category" binding:"omitempty,oneof=programming algorithm ai_assisted ctf"`
	RuleType        *string     `json:"rule_type" binding:"omitempty,oneof=acm oi ioi cf ctf awd isw diy"`
	DIYTemplateID   *uint       `json:"diy_template_id"`
	StartTime       *time.Time  `json:"start_time"`
	EndTime         *time.Time  `json:"end_time"`
	FreezeTime      *time.Time  `json:"freeze_time"`
	Description     *string     `json:"description"`
	RuleDescription *string     `json:"rule_description"`
	MaxTeamSize     *int        `json:"max_team_size"`
	SignupLimit     *int        `json:"signup_limit"`
	AllowViewOthers *bool       `json:"allow_view_others"`
	ShowRealtimeRank *bool      `json:"show_realtime_rank"`
	EnableAIHint    *bool       `json:"enable_ai_hint"`
	DIYRules        *string     `json:"diy_rules"`
}

// ContestProblemReq is a request for adding a problem to a contest.
type ContestProblemReq struct {
	ProblemID    uint    `json:"problem_id" binding:"required"`
	Score        float64 `json:"score"`
	DisplayOrder int     `json:"display_order"`
	ProblemLabel string  `json:"problem_label"`
}

// ContestDetailResponse is the enriched contest detail response.
type ContestDetailResponse struct {
	model.Contest
	IsSignedUp bool `json:"is_signed_up"`
}

// RankingEntry represents a single entry in the ranking.
type RankingEntry struct {
	Rank           int                    `json:"rank"`
	UserID         uint                   `json:"user_id"`
	Username       string                 `json:"username"`
	Avatar         string                 `json:"avatar"`
	SolvedCount    int                    `json:"solved_count"`
	TotalScore     float64                `json:"total_score"`
	Penalty        int                    `json:"penalty"`
	ProblemResults map[uint]*ProblemResult `json:"problem_results"`
}

// ProblemResult represents a user's result for a single problem in ranking.
type ProblemResult struct {
	IsSolved       bool      `json:"is_solved"`
	Score          float64   `json:"score"`
	SubmitCount    int       `json:"submit_count"`
	WrongAttempts  int       `json:"wrong_attempts"`
	AcceptedTime   *time.Time `json:"accepted_time,omitempty"`
}

// ContestService handles contest-related business logic.
type ContestService interface {
	ListContests(ctx context.Context, offset, limit int, filter ContestFilter) ([]model.Contest, int64, error)
	GetContestDetail(ctx context.Context, contestID uint, userID uint) (*ContestDetailResponse, error)
	CreateContest(ctx context.Context, adminID uint, req CreateContestRequest) (*model.Contest, error)
	UpdateContest(ctx context.Context, contestID uint, req UpdateContestRequest) (*model.Contest, error)
	DeleteContest(ctx context.Context, contestID uint) error
	SignupUser(ctx context.Context, contestID uint, userID uint) error
	WithdrawUser(ctx context.Context, contestID uint, userID uint) error
	GetContestProblems(ctx context.Context, contestID uint, userID uint) ([]model.ContestProblem, error)
	SubmitToContest(ctx context.Context, contestID uint, userID uint, req CreateSubmissionRequest, ipAddress string) (*model.Submission, error)
	GetRanking(ctx context.Context, contestID uint) ([]RankingEntry, error)
	GetFrozenRanking(ctx context.Context, contestID uint) ([]RankingEntry, error)
	UpdateContestStatus(ctx context.Context, contestID uint, status string) error
	CheckAndUpdateContestStatuses(ctx context.Context) error
	// Admin helpers
	AddContestProblem(ctx context.Context, contestID uint, req ContestProblemReq) error
	RemoveContestProblem(ctx context.Context, contestID uint, problemID uint) error
	GetContestSignups(ctx context.Context, contestID uint) ([]model.ContestSignup, error)
	// DIY templates
	CreateDIYTemplate(ctx context.Context, adminID uint, req CreateDIYTemplateRequest) (*model.DIYContestTemplate, error)
	GetDIYTemplate(ctx context.Context, id uint) (*model.DIYContestTemplate, error)
	UpdateDIYTemplate(ctx context.Context, id uint, req UpdateDIYTemplateRequest) (*model.DIYContestTemplate, error)
	DeleteDIYTemplate(ctx context.Context, id uint) error
	ListDIYTemplates(ctx context.Context, offset, limit int) ([]model.DIYContestTemplate, int64, error)
}

// CreateDIYTemplateRequest is the request for creating a DIY template.
type CreateDIYTemplateRequest struct {
	Name         string `json:"name" binding:"required,max=128"`
	ScoringRule  string `json:"scoring_rule"`
	PenaltyRule  string `json:"penalty_rule"`
	RankingRule  string `json:"ranking_rule"`
	Description  string `json:"description"`
}

// UpdateDIYTemplateRequest is the request for updating a DIY template.
type UpdateDIYTemplateRequest struct {
	Name         *string `json:"name" binding:"omitempty,max=128"`
	ScoringRule  *string `json:"scoring_rule"`
	PenaltyRule  *string `json:"penalty_rule"`
	RankingRule  *string `json:"ranking_rule"`
	Description  *string `json:"description"`
}

type contestService struct {
	contestRepo    repository.ContestRepository
	submissionRepo repository.SubmissionRepository
}

func NewContestService(contestRepo repository.ContestRepository, submissionRepo repository.SubmissionRepository) ContestService {
	return &contestService{
		contestRepo:    contestRepo,
		submissionRepo: submissionRepo,
	}
}

// detectStatus returns the current status based on time.
func detectStatus(contest *model.Contest) string {
	now := time.Now()
	if now.Before(contest.StartTime) {
		return "pending"
	}
	if now.After(contest.EndTime) {
		return "ended"
	}
	return "running"
}

// ListContests returns a paginated list of contests with auto-detected status.
func (s *contestService) ListContests(ctx context.Context, offset, limit int, filter ContestFilter) ([]model.Contest, int64, error) {
	filters := make(map[string]interface{})

	if filter.Status != "" {
		filters["status"] = filter.Status
	}
	if filter.Type != "" {
		filters["contest_type"] = filter.Type
	}
	if filter.Category != "" {
		filters["category"] = filter.Category
	}
	if filter.RuleType != "" {
		filters["rule_type"] = filter.RuleType
	}

	contests, total, err := s.contestRepo.List(ctx, offset, limit, filters, filter.Search, filter.SortBy)
	if err != nil {
		return nil, 0, err
	}

	// Auto-detect status based on time for each contest
	for i := range contests {
		contests[i].Status = detectStatus(&contests[i])
	}

	return contests, total, nil
}

// GetContestDetail returns contest detail with signup status.
func (s *contestService) GetContestDetail(ctx context.Context, contestID uint, userID uint) (*ContestDetailResponse, error) {
	contest, err := s.contestRepo.GetByID(ctx, contestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContestNotFound
		}
		return nil, err
	}

	// Auto-detect status
	contest.Status = detectStatus(contest)

	resp := &ContestDetailResponse{
		Contest: *contest,
	}

	if userID > 0 {
		signedUp, err := s.contestRepo.IsSignedUp(ctx, contestID, userID)
		if err == nil {
			resp.IsSignedUp = signedUp
		}
	}

	return resp, nil
}

// CreateContest creates a new contest with problems.
func (s *contestService) CreateContest(ctx context.Context, adminID uint, req CreateContestRequest) (*model.Contest, error) {
	if req.EndTime.Before(req.StartTime) {
		return nil, ErrInvalidContestTime
	}

	contest := &model.Contest{
		Title:            req.Title,
		ContestType:      defaultIfEmpty(req.ContestType, "individual"),
		Category:         defaultIfEmpty(req.Category, "programming"),
		RuleType:         defaultIfEmpty(req.RuleType, "acm"),
		DIYTemplateID:    req.DIYTemplateID,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		FreezeTime:       req.FreezeTime,
		Description:      req.Description,
		RuleDescription:  req.RuleDescription,
		MaxTeamSize:      defaultInt(req.MaxTeamSize, 3),
		SignupLimit:      req.SignupLimit,
		AllowViewOthers:  true,
		ShowRealtimeRank: true,
		EnableAIHint:     false,
		DIYRules:         req.DIYRules,
		Status:           "pending",
		CreatedBy:        adminID,
	}

	if err := s.contestRepo.Create(ctx, contest); err != nil {
		return nil, err
	}

	// Add problems if provided
	for _, p := range req.Problems {
		cp := &model.ContestProblem{
			ContestID:    contest.ID,
			ProblemID:    p.ProblemID,
			Score:        p.Score,
			DisplayOrder: p.DisplayOrder,
			ProblemLabel: p.ProblemLabel,
		}
		if err := s.contestRepo.AddProblem(ctx, cp); err != nil {
			// Log but don't fail the whole operation
			continue
		}
	}

	// Reload contest with problems
	contest, _ = s.contestRepo.GetByID(ctx, contest.ID)

	return contest, nil
}

// UpdateContest updates an existing contest.
func (s *contestService) UpdateContest(ctx context.Context, contestID uint, req UpdateContestRequest) (*model.Contest, error) {
	contest, err := s.contestRepo.GetByID(ctx, contestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContestNotFound
		}
		return nil, err
	}

	if req.Title != nil {
		contest.Title = *req.Title
	}
	if req.ContestType != nil {
		contest.ContestType = *req.ContestType
	}
	if req.Category != nil {
		contest.Category = *req.Category
	}
	if req.RuleType != nil {
		contest.RuleType = *req.RuleType
	}
	if req.DIYTemplateID != nil {
		contest.DIYTemplateID = *req.DIYTemplateID
	}
	if req.StartTime != nil {
		contest.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		contest.EndTime = *req.EndTime
	}
	if req.FreezeTime != nil {
		contest.FreezeTime = req.FreezeTime
	}
	if req.Description != nil {
		contest.Description = *req.Description
	}
	if req.RuleDescription != nil {
		contest.RuleDescription = *req.RuleDescription
	}
	if req.MaxTeamSize != nil {
		contest.MaxTeamSize = *req.MaxTeamSize
	}
	if req.SignupLimit != nil {
		contest.SignupLimit = *req.SignupLimit
	}
	if req.AllowViewOthers != nil {
		contest.AllowViewOthers = *req.AllowViewOthers
	}
	if req.ShowRealtimeRank != nil {
		contest.ShowRealtimeRank = *req.ShowRealtimeRank
	}
	if req.EnableAIHint != nil {
		contest.EnableAIHint = *req.EnableAIHint
	}
	if req.DIYRules != nil {
		contest.DIYRules = *req.DIYRules
	}

	// Validate times
	if contest.EndTime.Before(contest.StartTime) {
		return nil, ErrInvalidContestTime
	}

	if err := s.contestRepo.Update(ctx, contest); err != nil {
		return nil, err
	}

	return contest, nil
}

// DeleteContest soft-deletes a contest.
func (s *contestService) DeleteContest(ctx context.Context, contestID uint) error {
	_, err := s.contestRepo.GetByID(ctx, contestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContestNotFound
		}
		return err
	}
	return s.contestRepo.Delete(ctx, contestID)
}

// SignupUser signs up a user for a contest.
func (s *contestService) SignupUser(ctx context.Context, contestID uint, userID uint) error {
	contest, err := s.contestRepo.GetByID(ctx, contestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContestNotFound
		}
		return err
	}

	// Check if contest has ended
	if time.Now().After(contest.EndTime) {
		return ErrContestEnded
	}

	// Check if already signed up
	signedUp, err := s.contestRepo.IsSignedUp(ctx, contestID, userID)
	if err != nil {
		return err
	}
	if signedUp {
		return ErrAlreadySignedUp
	}

	// Check signup limit
	if contest.SignupLimit > 0 {
		count, err := s.contestRepo.CountSignups(ctx, contestID)
		if err != nil {
			return err
		}
		if int(count) >= contest.SignupLimit {
			return ErrContestFull
		}
	}

	signup := &model.ContestSignup{
		ContestID:  contestID,
		UserID:     userID,
		SignupTime: time.Now(),
	}

	if err := s.contestRepo.Signup(ctx, signup); err != nil {
		return err
	}

	// Increment participant count
	_ = s.contestRepo.IncrementParticipantCount(ctx, contestID)

	return nil
}

// WithdrawUser withdraws a user from a contest.
func (s *contestService) WithdrawUser(ctx context.Context, contestID uint, userID uint) error {
	contest, err := s.contestRepo.GetByID(ctx, contestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContestNotFound
		}
		return err
	}

	// Cannot withdraw from running or ended contest
	status := detectStatus(contest)
	if status == "running" || status == "ended" {
		return ErrContestNotRunning
	}

	if err := s.contestRepo.Withdraw(ctx, contestID, userID); err != nil {
		return err
	}

	_ = s.contestRepo.DecrementParticipantCount(ctx, contestID)

	return nil
}

// GetContestProblems returns contest problems (only if signed up and contest is running).
func (s *contestService) GetContestProblems(ctx context.Context, contestID uint, userID uint) ([]model.ContestProblem, error) {
	contest, err := s.contestRepo.GetByID(ctx, contestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContestNotFound
		}
		return nil, err
	}

	// If userID is 0, allow viewing (for public contests that have ended)
	if userID > 0 {
		signedUp, err := s.contestRepo.IsSignedUp(ctx, contestID, userID)
		if err != nil {
			return nil, err
		}
		if !signedUp {
			return nil, ErrNotSignedUp
		}
	}

	problems, err := s.contestRepo.GetProblems(ctx, contestID)
	if err != nil {
		return nil, err
	}

	return problems, nil
}

// SubmitToContest creates a submission linked to a contest.
func (s *contestService) SubmitToContest(ctx context.Context, contestID uint, userID uint, req CreateSubmissionRequest, ipAddress string) (*model.Submission, error) {
	contest, err := s.contestRepo.GetByID(ctx, contestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContestNotFound
		}
		return nil, err
	}

	// Check contest is running
	status := detectStatus(contest)
	if status != "running" {
		return nil, ErrContestNotRunning
	}

	// Check user is signed up
	signedUp, err := s.contestRepo.IsSignedUp(ctx, contestID, userID)
	if err != nil {
		return nil, err
	}
	if !signedUp {
		return nil, ErrNotSignedUp
	}

	// Verify the problem is part of this contest
	problems, err := s.contestRepo.GetProblems(ctx, contestID)
	if err != nil {
		return nil, err
	}

	var found bool
	for _, cp := range problems {
		if cp.ProblemID == req.ProblemID {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrProblemNotFound
	}

	// Create submission with contest ID
	submission := &model.Submission{
		UserID:      userID,
		ProblemID:   req.ProblemID,
		ContestID:   contestID,
		Language:    req.Language,
		IPAddress:   ipAddress,
		SubmitTime:  time.Now(),
		JudgeResult: "pending",
	}

	// Handle CTF submissions
	if contest.Category == "ctf" || req.Language == "ctf" {
		if req.CTFAnswer == "" {
			return nil, ErrCTFAnswerEmpty
		}
		submission.CTFAnswer = req.CTFAnswer
		submission.Language = "ctf"
		submission.CodeLength = len(req.CTFAnswer)

		// For CTF, we need to check the flag - but that's done by the judge/submission service
		// For now, create the submission and let the judge handle it
	} else {
		if !ValidLanguages[req.Language] {
			return nil, ErrInvalidLanguage
		}
		if req.Code == "" {
			return nil, ErrCodeEmpty
		}
		submission.CodeLength = len(req.Code)
	}

	if err := s.submissionRepo.Create(ctx, submission); err != nil {
		return nil, err
	}

	return submission, nil
}

// GetRanking calculates and returns the ranking for a contest.
func (s *contestService) GetRanking(ctx context.Context, contestID uint) ([]RankingEntry, error) {
	contest, err := s.contestRepo.GetByID(ctx, contestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContestNotFound
		}
		return nil, err
	}

	submissions, err := s.contestRepo.GetContestRanking(ctx, contestID)
	if err != nil {
		return nil, err
	}

	return s.calculateRanking(contest, submissions)
}

// GetFrozenRanking calculates ranking frozen at freeze time.
func (s *contestService) GetFrozenRanking(ctx context.Context, contestID uint) ([]RankingEntry, error) {
	contest, err := s.contestRepo.GetByID(ctx, contestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContestNotFound
		}
		return nil, err
	}

	if contest.FreezeTime == nil {
		// No freeze time set, return empty ranking
		return []RankingEntry{}, nil
	}

	submissions, err := s.contestRepo.GetFrozenRanking(ctx, contestID, *contest.FreezeTime)
	if err != nil {
		return nil, err
	}

	return s.calculateRanking(contest, submissions)
}

// calculateRanking computes ranking entries based on contest rule type.
func (s *contestService) calculateRanking(contest *model.Contest, submissions []model.Submission) ([]RankingEntry, error) {
	// Get contest problems for reference
	problems, err := s.contestRepo.GetProblems(context.Background(), contest.ID)
	if err != nil {
		return nil, err
	}

	// Build a map of problem scores by problem ID
	problemScores := make(map[uint]float64)
	for _, cp := range problems {
		problemScores[cp.ProblemID] = cp.Score
	}

	// Group submissions by user
	type userData struct {
		userID       uint
		username     string
		avatar       string
		problemStats map[uint]*ProblemResult
	}

	userMap := make(map[uint]*userData)

	for _, sub := range submissions {
		uid := sub.UserID
		pid := sub.ProblemID

		if _, ok := userMap[uid]; !ok {
			userMap[uid] = &userData{
				userID:       uid,
				username:     sub.User.Username,
				avatar:       sub.User.Avatar,
				problemStats: make(map[uint]*ProblemResult),
			}
		}

		if _, ok := userMap[uid].problemStats[pid]; !ok {
			userMap[uid].problemStats[pid] = &ProblemResult{}
		}

		pr := userMap[uid].problemStats[pid]
		pr.SubmitCount++

		if sub.JudgeResult == "AC" {
			pr.IsSolved = true
			pr.Score = problemScores[pid]
			if pr.AcceptedTime == nil || sub.SubmitTime.Before(*pr.AcceptedTime) {
				pr.AcceptedTime = &sub.SubmitTime
			}
		} else if sub.JudgeResult != "pending" && sub.JudgeResult != "judging" {
			// Count wrong attempts only before first AC
			if !pr.IsSolved {
				pr.WrongAttempts++
			}
			// For OI/IOI, track max score
			if contest.RuleType == "oi" || contest.RuleType == "ioi" {
				if sub.JudgeScore > pr.Score {
					pr.Score = sub.JudgeScore
				}
			}
		}
	}

	// Build ranking entries
	var entries []RankingEntry
	for _, ud := range userMap {
		entry := RankingEntry{
			UserID:         ud.userID,
			Username:       ud.username,
			Avatar:         ud.avatar,
			ProblemResults: ud.problemStats,
		}

		for _, pr := range ud.problemStats {
			if pr.IsSolved {
				entry.SolvedCount++
				entry.TotalScore += pr.Score
			} else {
				entry.TotalScore += pr.Score // For OI/IOI, partial scores count
			}
		}

		entries = append(entries, entry)
	}

	// Sort based on rule type
	switch contest.RuleType {
	case "acm":
		s.sortACMRanking(entries, contest)
	case "oi", "ioi":
		s.sortOIRanking(entries)
	case "cf":
		s.sortCFRanking(entries, contest)
	case "ctf":
		s.sortCTFRanking(entries, contest)
	case "diy":
		s.sortDIYRanking(entries, contest)
	default:
		s.sortACMRanking(entries, contest)
	}

	// Assign ranks
	for i := range entries {
		entries[i].Rank = i + 1
	}

	return entries, nil
}

// sortACMRanking sorts by solved count desc, then penalty asc.
// Penalty = 20 minutes per wrong submission before AC.
func (s *contestService) sortACMRanking(entries []RankingEntry, contest *model.Contest) {
	contestStart := contest.StartTime

	for i := range entries {
		penalty := 0
		for _, pr := range entries[i].ProblemResults {
			if pr.IsSolved && pr.AcceptedTime != nil {
				// Time in minutes from contest start
				minutes := int(pr.AcceptedTime.Sub(contestStart).Minutes())
				penalty += minutes
				penalty += pr.WrongAttempts * 20
			}
		}
		entries[i].Penalty = penalty
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].SolvedCount != entries[j].SolvedCount {
			return entries[i].SolvedCount > entries[j].SolvedCount
		}
		return entries[i].Penalty < entries[j].Penalty
	})
}

// sortOIRanking sorts by total score desc, then submission count asc.
func (s *contestService) sortOIRanking(entries []RankingEntry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].TotalScore != entries[j].TotalScore {
			return entries[i].TotalScore > entries[j].TotalScore
		}
		// Tie-breaker: fewer total submissions wins
		iSubs := 0
		jSubs := 0
		for _, pr := range entries[i].ProblemResults {
			iSubs += pr.SubmitCount
		}
		for _, pr := range entries[j].ProblemResults {
			jSubs += pr.SubmitCount
		}
		return iSubs < jSubs
	})
}

// sortCFRanking sorts by score desc with decay factor based on solve time.
// Score = max_score * decay_factor, where decay_factor decreases over time.
func (s *contestService) sortCFRanking(entries []RankingEntry, contest *model.Contest) {
	contestStart := contest.StartTime
	totalMinutes := contest.EndTime.Sub(contestStart).Minutes()

	for i := range entries {
		totalScore := 0.0
		for _, pr := range entries[i].ProblemResults {
			if pr.IsSolved && pr.AcceptedTime != nil {
				solveMinutes := pr.AcceptedTime.Sub(contestStart).Minutes()
				// Decay factor: starts at 1.0, decreases to 0.3 over contest duration
				decayFactor := 1.0 - 0.7*(solveMinutes/totalMinutes)
				if decayFactor < 0.3 {
					decayFactor = 0.3
				}
				totalScore += pr.Score * decayFactor
			}
		}
		entries[i].TotalScore = math.Round(totalScore*100) / 100
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].TotalScore != entries[j].TotalScore {
			return entries[i].TotalScore > entries[j].TotalScore
		}
		// Tie-breaker: fewer wrong attempts
		iWrong := 0
		jWrong := 0
		for _, pr := range entries[i].ProblemResults {
			iWrong += pr.WrongAttempts
		}
		for _, pr := range entries[j].ProblemResults {
			jWrong += pr.WrongAttempts
		}
		return iWrong < jWrong
	})
}

// sortCTFRanking sorts by solved count desc, then last solve time asc.
func (s *contestService) sortCTFRanking(entries []RankingEntry, contest *model.Contest) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].SolvedCount != entries[j].SolvedCount {
			return entries[i].SolvedCount > entries[j].SolvedCount
		}
		// Tie-breaker: earlier last solve time wins
		iLast := getLastSolveTime(entries[i].ProblemResults)
		jLast := getLastSolveTime(entries[j].ProblemResults)
		if iLast == nil && jLast == nil {
			return false
		}
		if iLast == nil {
			return false
		}
		if jLast == nil {
			return true
		}
		return iLast.Before(*jLast)
	})
}

// sortDIYRanking uses custom ranking rules from diy_rules JSON.
func (s *contestService) sortDIYRanking(entries []RankingEntry, contest *model.Contest) {
	// Parse DIY rules
	if contest.DIYRules == "" {
		// Default to ACM-style
		s.sortACMRanking(entries, contest)
		return
	}

	var rules struct {
		SortBy     string `json:"sort_by"`      // "score", "solved", "penalty"
		SortOrder  string `json:"sort_order"`    // "asc", "desc"
		TieBreaker string `json:"tie_breaker"`   // "time", "submissions"
	}

	if err := json.Unmarshal([]byte(contest.DIYRules), &rules); err != nil {
		s.sortACMRanking(entries, contest)
		return
	}

	switch rules.SortBy {
	case "score":
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].TotalScore != entries[j].TotalScore {
				if rules.SortOrder == "asc" {
					return entries[i].TotalScore < entries[j].TotalScore
				}
				return entries[i].TotalScore > entries[j].TotalScore
			}
			return applyTieBreaker(entries[i], entries[j], rules.TieBreaker)
		})
	case "solved":
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].SolvedCount != entries[j].SolvedCount {
				if rules.SortOrder == "asc" {
					return entries[i].SolvedCount < entries[j].SolvedCount
				}
				return entries[i].SolvedCount > entries[j].SolvedCount
			}
			return applyTieBreaker(entries[i], entries[j], rules.TieBreaker)
		})
	case "penalty":
		// Calculate penalty for all entries first
		contestStart := contest.StartTime
		for i := range entries {
			penalty := 0
			for _, pr := range entries[i].ProblemResults {
				if pr.IsSolved && pr.AcceptedTime != nil {
					minutes := int(pr.AcceptedTime.Sub(contestStart).Minutes())
					penalty += minutes
					penalty += pr.WrongAttempts * 20
				}
			}
			entries[i].Penalty = penalty
		}
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].Penalty != entries[j].Penalty {
				if rules.SortOrder == "asc" {
					return entries[i].Penalty < entries[j].Penalty
				}
				return entries[i].Penalty > entries[j].Penalty
			}
			return applyTieBreaker(entries[i], entries[j], rules.TieBreaker)
		})
	default:
		s.sortACMRanking(entries, contest)
	}
}

func applyTieBreaker(a, b RankingEntry, tieBreaker string) bool {
	switch tieBreaker {
	case "time":
		aLast := getLastSolveTime(a.ProblemResults)
		bLast := getLastSolveTime(b.ProblemResults)
		if aLast != nil && bLast != nil {
			return aLast.Before(*bLast)
		}
		return false
	case "submissions":
		aSubs := 0
		bSubs := 0
		for _, pr := range a.ProblemResults {
			aSubs += pr.SubmitCount
		}
		for _, pr := range b.ProblemResults {
			bSubs += pr.SubmitCount
		}
		return aSubs < bSubs
	default:
		return false
	}
}

func getLastSolveTime(problemResults map[uint]*ProblemResult) *time.Time {
	var lastTime *time.Time
	for _, pr := range problemResults {
		if pr.IsSolved && pr.AcceptedTime != nil {
			if lastTime == nil || pr.AcceptedTime.After(*lastTime) {
				lastTime = pr.AcceptedTime
			}
		}
	}
	return lastTime
}

// UpdateContestStatus updates a contest's status.
func (s *contestService) UpdateContestStatus(ctx context.Context, contestID uint, status string) error {
	validStatuses := map[string]bool{
		"pending": true, "running": true, "ended": true, "cancelled": true,
	}
	if !validStatuses[status] {
		return ErrInvalidContestStatus
	}

	_, err := s.contestRepo.GetByID(ctx, contestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContestNotFound
		}
		return err
	}

	return s.contestRepo.UpdateStatus(ctx, contestID, status)
}

// CheckAndUpdateContestStatuses auto-updates contest statuses based on time.
func (s *contestService) CheckAndUpdateContestStatuses(ctx context.Context) error {
	contests, err := s.contestRepo.GetContestsNeedStatusUpdate(ctx)
	if err != nil {
		return err
	}

	for _, contest := range contests {
		newStatus := detectStatus(&contest)
		if newStatus != contest.Status {
			_ = s.contestRepo.UpdateStatus(ctx, contest.ID, newStatus)
		}
	}

	return nil
}

// AddContestProblem adds a problem to a contest.
func (s *contestService) AddContestProblem(ctx context.Context, contestID uint, req ContestProblemReq) error {
	_, err := s.contestRepo.GetByID(ctx, contestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContestNotFound
		}
		return err
	}

	cp := &model.ContestProblem{
		ContestID:    contestID,
		ProblemID:    req.ProblemID,
		Score:        req.Score,
		DisplayOrder: req.DisplayOrder,
		ProblemLabel: req.ProblemLabel,
	}

	return s.contestRepo.AddProblem(ctx, cp)
}

// RemoveContestProblem removes a problem from a contest.
func (s *contestService) RemoveContestProblem(ctx context.Context, contestID uint, problemID uint) error {
	_, err := s.contestRepo.GetByID(ctx, contestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContestNotFound
		}
		return err
	}

	return s.contestRepo.RemoveProblem(ctx, contestID, problemID)
}

// GetContestSignups returns all signups for a contest.
func (s *contestService) GetContestSignups(ctx context.Context, contestID uint) ([]model.ContestSignup, error) {
	_, err := s.contestRepo.GetByID(ctx, contestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContestNotFound
		}
		return nil, err
	}

	return s.contestRepo.GetSignups(ctx, contestID)
}

// ==================== DIY Template Methods ====================

func (s *contestService) CreateDIYTemplate(ctx context.Context, adminID uint, req CreateDIYTemplateRequest) (*model.DIYContestTemplate, error) {
	tmpl := &model.DIYContestTemplate{
		Name:        req.Name,
		ScoringRule: req.ScoringRule,
		PenaltyRule: req.PenaltyRule,
		RankingRule: req.RankingRule,
		Description: req.Description,
		CreatedBy:   adminID,
	}

	if err := s.contestRepo.CreateDIYTemplate(ctx, tmpl); err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (s *contestService) GetDIYTemplate(ctx context.Context, id uint) (*model.DIYContestTemplate, error) {
	tmpl, err := s.contestRepo.GetDIYTemplateByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContestNotFound
		}
		return nil, err
	}
	return tmpl, nil
}

func (s *contestService) UpdateDIYTemplate(ctx context.Context, id uint, req UpdateDIYTemplateRequest) (*model.DIYContestTemplate, error) {
	tmpl, err := s.contestRepo.GetDIYTemplateByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContestNotFound
		}
		return nil, err
	}

	if req.Name != nil {
		tmpl.Name = *req.Name
	}
	if req.ScoringRule != nil {
		tmpl.ScoringRule = *req.ScoringRule
	}
	if req.PenaltyRule != nil {
		tmpl.PenaltyRule = *req.PenaltyRule
	}
	if req.RankingRule != nil {
		tmpl.RankingRule = *req.RankingRule
	}
	if req.Description != nil {
		tmpl.Description = *req.Description
	}

	if err := s.contestRepo.UpdateDIYTemplate(ctx, tmpl); err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (s *contestService) DeleteDIYTemplate(ctx context.Context, id uint) error {
	_, err := s.contestRepo.GetDIYTemplateByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContestNotFound
		}
		return err
	}
	return s.contestRepo.DeleteDIYTemplate(ctx, id)
}

func (s *contestService) ListDIYTemplates(ctx context.Context, offset, limit int) ([]model.DIYContestTemplate, int64, error) {
	return s.contestRepo.ListDIYTemplates(ctx, offset, limit)
}
