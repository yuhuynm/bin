package service

import (
	"errors"
	"go-bin/internal/dto"
	"go-bin/internal/entity"
	"go-bin/internal/repository"

	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UserService interface {
	Create(input dto.CreateUserRequest) (*entity.User, error)
	List() ([]entity.User, error)
	GetByID(id uint) (*entity.User, error)
	Update(id uint, input dto.UpdateUserRequest) (*entity.User, error)
	Delete(id uint) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(input dto.CreateUserRequest) (*entity.User, error) {
	user := &entity.User{
		Name:  input.Name,
		Email: input.Email,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) List() ([]entity.User, error) {
	return s.repo.FindAll()
}

func (s *userService) GetByID(id uint) (*entity.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *userService) Update(id uint, input dto.UpdateUserRequest) (*entity.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	user.Name = input.Name
	user.Email = input.Email

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(id uint) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	return s.repo.Delete(user.ID)
}
