package repository

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	v1 "projectName/api/v1"
	"projectName/internal/model"
)

type CollegeRepository interface {
	// db
	GetCollegeByCollegeId(ctx context.Context, id int64) (*model.College, error)
	GetCollegeList(ctx context.Context) ([]*model.College, error)
}

func NewCollegeRepository(
	repository *Repository,
) CollegeRepository {
	return &collegeRepository{
		Repository: repository,
	}
}

type collegeRepository struct {
	*Repository
}

// GetCollege 根据ID获取单个学院的信息
func (r *collegeRepository) GetCollegeByCollegeId(ctx context.Context, college_id int64) (*model.College, error) {
	var college model.College
	// 查询单个学院
	if err := r.DB(ctx).Table("sys_colleges").Where("college_id = ?", college_id).First(&college).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 如果未找到学院，返回未找到错误
			return nil, v1.ErrNotFound
		}
		// 记录查询错误
		r.logger.WithContext(ctx).Error("collegeRepository.GetCollege error", zap.Error(err))
		return nil, err
	}
	// 返回查询到的学院信息
	return &college, nil
}

// GetCollegeList 获取多个学院的信息，支持通过条件过滤
func (r *collegeRepository) GetCollegeList(ctx context.Context) ([]*model.College, error) {
	var collegeList []*model.College
	// 查询所有未删除的学院，条件可根据需要调整
	if err := r.DB(ctx).Table("sys_colleges").Where("is_deleted = ?", 0).Find(&collegeList).Error; err != nil {
		// 如果查询出错，记录错误并返回
		r.logger.WithContext(ctx).Error("collegeRepository.GetCollegeList error", zap.Error(err))
		return nil, err
	}
	// 返回学院列表
	return collegeList, nil
}
