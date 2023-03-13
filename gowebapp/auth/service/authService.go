package services

import (
	"auth/dao"
)

type AuthService struct {
	dao dao.AuthDataAccess
}

func NewAuthService(dao dao.AuthDataAccess) *AuthService {
	return &AuthService{
		dao: dao,
	}
}
