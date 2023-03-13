package services

import (
	"web-api/app/dummyAuth/dao"
)

type AuthService struct {
	dao dao.AuthDataAccess
}

func NewAuthService(dao dao.AuthDataAccess) *AuthService {
	return &AuthService{
		dao: dao,
	}
}
