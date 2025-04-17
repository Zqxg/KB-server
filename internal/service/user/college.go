package user

import (
	"context"
	"projectName/internal/model"
	"projectName/internal/repository"
	"projectName/internal/service"
)

type CollegeService interface {
	GetCollege(ctx context.Context, id uint) (*model.College, error)
	GetCollegeList(ctx context.Context) ([]*model.College, error)
}

func NewCollegeService(
	service *service.Service,
	collegeRepository repository.CollegeRepository,
) CollegeService {
	return &collegeService{
		Service:           service,
		collegeRepository: collegeRepository,
	}
}

type collegeService struct {
	*service.Service
	collegeRepository repository.CollegeRepository
}

func (s *collegeService) GetCollege(ctx context.Context, collegeId uint) (*model.College, error) {
	return s.collegeRepository.GetCollegeByCollegeId(ctx, collegeId)
}

func (s *collegeService) GetCollegeList(ctx context.Context) ([]*model.College, error) {
	return s.collegeRepository.GetCollegeList(ctx)
}
