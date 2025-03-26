package admin

import (
	"projectName/internal/repository"
	"projectName/internal/service"
)

type AdminService interface {
	//createPublicKB(ctx context.Context, req *CreatePublicKBRequest) (*CreatePublicKBResponse, error)
}

func NewAdminService(
	service *service.Service,
	userRepository repository.UserRepository,
	articleRepository repository.ArticleRepository,
	teamRepository repository.TeamRepository,
	kbRepository repository.KBRepository,
) AdminService {
	return &adminService{
		Service:           service,
		userRepository:    userRepository,
		teamRepository:    teamRepository,
		articleRepository: articleRepository,
		kbRepository:      kbRepository,
	}
}

type adminService struct {
	*service.Service
	userRepository    repository.UserRepository
	teamRepository    repository.TeamRepository
	articleRepository repository.ArticleRepository
	kbRepository      repository.KBRepository
}
