package service

import (
	"context"
	"design-pattern/internal/entity"
	"design-pattern/internal/repository"
	"design-pattern/pkg/cache"
	"design-pattern/pkg/token"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserService interface {
	FindAll(ctx context.Context) ([]entity.User, error)
	Login(ctx context.Context, username, password string) (string, error)
}

type userService struct {
	userRepository repository.UserRepository
	tokenUseCase   token.TokenUseCase
	cacheable      cache.Cacheable
}

func NewUserService(
	userRepository repository.UserRepository,
	tokenUseCase token.TokenUseCase,
	cacheable cache.Cacheable,
) UserService {
	return &userService{userRepository, tokenUseCase, cacheable}
}

func (s *userService) FindAll(ctx context.Context) (result []entity.User, err error) {
	keyFindAll := "design-pattern-api:users:find-all"
	data := s.cacheable.Get(keyFindAll)
	if data == "" {
		result, err = s.userRepository.FindAll(ctx)
		if err != nil {
			return nil, err
		}

		marshalledData, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		err = s.cacheable.Set(keyFindAll, marshalledData, 5*time.Minute)
		if err != nil {
			return nil, err
		}
	} else {
		err = json.Unmarshal([]byte(data), &result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (s *userService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepository.FindByUsername(ctx, username)
	if err != nil {
		log.Println(err.Error())
		return "", errors.New("username or password invalid")
	}

	if user.Password != password {
		return "", errors.New("username or password invalid")
	}

	expiredTime := time.Now().Local().Add(time.Minute * 10)

	claims := token.JwtCustomClaims{
		Username: user.Username,
		Role:     user.Role,
		FullName: user.FullName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "design-pattern",
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	token, err := s.tokenUseCase.GenerateAccessToken(claims)
	if err != nil {
		return "", errors.New("ada kesalahan di server")
	}

	return token, nil
}
