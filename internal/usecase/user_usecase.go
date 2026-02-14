package usecase

import (
	"errors"
	"test-elabram/internal/domain"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) GetAllUsers() ([]domain.User, error) {
	return u.userRepo.GetAll()
}

func (u *userUsecase) GetUserByID(id int) (*domain.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid ID")
	}
	return u.userRepo.GetByID(id)
}

func (u *userUsecase) CreateUser(user *domain.User) error {
	if user.Username == "" || user.Email == "" {
		return errors.New("username and email are required")
	}
	return u.userRepo.Create(user)
}
