package repository

import (
	"context"

	"github.com/Yogdunana/yogduoj/backend/internal/model"
	"gorm.io/gorm"
)

type ProblemRepository interface {
	Create(ctx context.Context, problem *model.Problem) error
	GetByID(ctx context.Context, id uint) (*model.Problem, error)
	Update(ctx context.Context, problem *model.Problem) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int, filters map[string]interface{}, search string, sort, order string) ([]model.Problem, int64, error)
	GetByContestID(ctx context.Context, contestID uint) ([]model.Problem, error)
	GetTags(ctx context.Context, problemID uint) ([]model.Tag, error)
	AddTag(ctx context.Context, problemID, tagID uint) error
	RemoveTag(ctx context.Context, problemID, tagID uint) error
	GetProblemTags(ctx context.Context, problemIDs []uint) (map[uint][]model.Tag, error)
	GetSamples(ctx context.Context, problemID uint) ([]model.Sample, error)
	CreateSample(ctx context.Context, sample *model.Sample) error
	GetTestData(ctx context.Context, problemID uint) ([]model.TestData, error)
	CreateTestData(ctx context.Context, testData *model.TestData) error
	DeleteTestData(ctx context.Context, dataID uint) error
	GetTestDataByID(ctx context.Context, dataID uint) (*model.TestData, error)
	IncrementSubmitCount(ctx context.Context, problemID uint) error
	IncrementAcceptedCount(ctx context.Context, problemID uint) error
	GetAttachments(ctx context.Context, problemID uint) ([]model.ProblemAttachment, error)
}

type problemRepository struct {
	db *gorm.DB
}

func NewProblemRepository(db *gorm.DB) ProblemRepository {
	return &problemRepository{db: db}
}

func (r *problemRepository) Create(ctx context.Context, problem *model.Problem) error {
	return r.db.WithContext(ctx).Create(problem).Error
}

func (r *problemRepository) GetByID(ctx context.Context, id uint) (*model.Problem, error) {
	var problem model.Problem
	if err := r.db.WithContext(ctx).Preload("Tags").Preload("Creator").First(&problem, id).Error; err != nil {
		return nil, err
	}
	return &problem, nil
}

func (r *problemRepository) Update(ctx context.Context, problem *model.Problem) error {
	return r.db.WithContext(ctx).Save(problem).Error
}

func (r *problemRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete associated problem_tags
		if err := tx.Exec("DELETE FROM problem_tags WHERE problem_id = ?", id).Error; err != nil {
			return err
		}
		// Delete associated samples
		if err := tx.Where("problem_id = ?", id).Delete(&model.Sample{}).Error; err != nil {
			return err
		}
		// Delete associated test data
		if err := tx.Where("problem_id = ?", id).Delete(&model.TestData{}).Error; err != nil {
			return err
		}
		// Delete the problem itself
		return tx.Delete(&model.Problem{}, id).Error
	})
}

// List returns a paginated list of problems with filtering, searching, and sorting.
func (r *problemRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}, search string, sort, order string) ([]model.Problem, int64, error) {
	var problems []model.Problem
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Problem{})

	// Apply filters
	for key, value := range filters {
		if value != nil && value != "" {
			query = query.Where(key, value)
		}
	}

	// Apply search by title
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	sortColumn := "id"
	switch sort {
	case "accepted_count", "submit_count", "difficulty", "created_at":
		sortColumn = sort
	}
	sortOrder := "ASC"
	if order == "desc" {
		sortOrder = "DESC"
	}
	query = query.Order(sortColumn + " " + sortOrder)

	// Fetch results
	if err := query.Preload("Tags").Preload("Creator").Offset(offset).Limit(limit).Find(&problems).Error; err != nil {
		return nil, 0, err
	}

	return problems, total, nil
}

func (r *problemRepository) GetByContestID(ctx context.Context, contestID uint) ([]model.Problem, error) {
	var problems []model.Problem
	if err := r.db.WithContext(ctx).
		Joins("JOIN contest_problems ON contest_problems.problem_id = problems.id").
		Where("contest_problems.contest_id = ?", contestID).
		Find(&problems).Error; err != nil {
		return nil, err
	}
	return problems, nil
}

func (r *problemRepository) GetTags(ctx context.Context, problemID uint) ([]model.Tag, error) {
	var tags []model.Tag
	if err := r.db.WithContext(ctx).
		Joins("JOIN problem_tags ON problem_tags.tag_id = tags.id").
		Where("problem_tags.problem_id = ?", problemID).
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *problemRepository) AddTag(ctx context.Context, problemID, tagID uint) error {
	return r.db.WithContext(ctx).Exec(
		"INSERT IGNORE INTO problem_tags (problem_id, tag_id) VALUES (?, ?)",
		problemID, tagID,
	).Error
}

func (r *problemRepository) RemoveTag(ctx context.Context, problemID, tagID uint) error {
	return r.db.WithContext(ctx).Exec(
		"DELETE FROM problem_tags WHERE problem_id = ? AND tag_id = ?",
		problemID, tagID,
	).Error
}

// GetProblemTags batch-fetches tags for multiple problems, returning a map of problemID -> []Tag.
func (r *problemRepository) GetProblemTags(ctx context.Context, problemIDs []uint) (map[uint][]model.Tag, error) {
	result := make(map[uint][]model.Tag)
	if len(problemIDs) == 0 {
		return result, nil
	}

	type problemTagRow struct {
		ProblemID uint
		TagID     uint
		TagName   string
	}

	var rows []problemTagRow
	if err := r.db.WithContext(ctx).
		Table("problem_tags").
		Select("problem_tags.problem_id, tags.id as tag_id, tags.name as tag_name").
		Joins("JOIN tags ON tags.id = problem_tags.tag_id").
		Where("problem_tags.problem_id IN ?", problemIDs).
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		tag := model.Tag{ID: row.TagID, Name: row.TagName}
		result[row.ProblemID] = append(result[row.ProblemID], tag)
	}

	return result, nil
}

func (r *problemRepository) GetSamples(ctx context.Context, problemID uint) ([]model.Sample, error) {
	var samples []model.Sample
	if err := r.db.WithContext(ctx).
		Where("problem_id = ?", problemID).
		Order("display_order ASC").
		Preload("TestData").
		Find(&samples).Error; err != nil {
		return nil, err
	}
	return samples, nil
}

func (r *problemRepository) CreateSample(ctx context.Context, sample *model.Sample) error {
	return r.db.WithContext(ctx).Create(sample).Error
}

func (r *problemRepository) GetTestData(ctx context.Context, problemID uint) ([]model.TestData, error) {
	var testData []model.TestData
	if err := r.db.WithContext(ctx).Where("problem_id = ?", problemID).Find(&testData).Error; err != nil {
		return nil, err
	}
	return testData, nil
}

func (r *problemRepository) CreateTestData(ctx context.Context, testData *model.TestData) error {
	return r.db.WithContext(ctx).Create(testData).Error
}

func (r *problemRepository) DeleteTestData(ctx context.Context, dataID uint) error {
	return r.db.WithContext(ctx).Delete(&model.TestData{}, dataID).Error
}

func (r *problemRepository) GetTestDataByID(ctx context.Context, dataID uint) (*model.TestData, error) {
	var testData model.TestData
	if err := r.db.WithContext(ctx).First(&testData, dataID).Error; err != nil {
		return nil, err
	}
	return &testData, nil
}

func (r *problemRepository) IncrementSubmitCount(ctx context.Context, problemID uint) error {
	return r.db.WithContext(ctx).Model(&model.Problem{}).Where("id = ?", problemID).
		UpdateColumn("submit_count", gorm.Expr("submit_count + 1")).Error
}

func (r *problemRepository) IncrementAcceptedCount(ctx context.Context, problemID uint) error {
	return r.db.WithContext(ctx).Model(&model.Problem{}).Where("id = ?", problemID).
		UpdateColumn("accepted_count", gorm.Expr("accepted_count + 1")).Error
}

// GetAttachments returns the parsed attachments for a problem.
func (r *problemRepository) GetAttachments(ctx context.Context, problemID uint) ([]model.ProblemAttachment, error) {
	var problem model.Problem
	if err := r.db.WithContext(ctx).Select("attachments").First(&problem, problemID).Error; err != nil {
		return nil, err
	}
	return model.ParseAttachments(problem.Attachments), nil
}
