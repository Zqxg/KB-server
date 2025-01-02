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
