package service

import (
	"projectName/internal/repository"
	"projectName/pkg/jwt"
	"projectName/pkg/log"
	"projectName/pkg/sid"
)

type Service struct {
	Logger *log.Logger
	Sid    *sid.Sid
	Jwt    *jwt.JWT
	Tm     repository.Transaction
}

func NewService(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
	jwt *jwt.JWT,
) *Service {
	return &Service{
		Logger: logger,
		Sid:    sid,
		Jwt:    jwt,
		Tm:     tm,
	}
}

// page初始化
func InitPage(pageIndex int, pageSize int) (int, int) {
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 10 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageIndex, pageSize
}
