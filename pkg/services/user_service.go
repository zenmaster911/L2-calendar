package services

import "github.com/znmaster911/L2-calendar/pkg/repositories"

type UserService struct {
	repo repositories.Users
}

func NewUserService(repo repositories.Users) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) NewUser() (int, error) {
	return s.repo.NewUser()
}
