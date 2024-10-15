package service

import (
	"context"
	"design-pattern/internal/entity"
	"design-pattern/internal/repository"
	"errors"
	"log"
)

type UserService interface {
	FindAll(ctx context.Context) ([]entity.User, error)
	Login(ctx context.Context, username, password string) (*entity.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository}
}

func (s *userService) FindAll(ctx context.Context) ([]entity.User, error) {
	users, err := s.userRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) Login(ctx context.Context, username, password string) (*entity.User, error) {
	user, err := s.userRepository.FindByUsername(ctx, username)
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("username or password invalid")
	}

	if user.Password != password {
		return nil, errors.New("username or password invalid")
	}

	return user, nil
}
