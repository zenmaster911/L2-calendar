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

func (s *UserService) NewUser(username string) (int64, error) {
	return s.repo.NewUser(username)
}

func (s *UserService) LogIn(username string) (int64, error) {
	return s.repo.LogIn(username)
}
