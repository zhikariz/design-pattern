package service_test

import (
	"context"
	"design-pattern/internal/entity"
	"design-pattern/internal/service"
	mock_cache "design-pattern/test/mock/pkg/cache"
	mock_token "design-pattern/test/mock/pkg/token"
	mock_repository "design-pattern/test/mock/repository"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UserTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	repo         *mock_repository.MockUserRepository
	tokenUseCase *mock_token.MockTokenUseCase
	cache        *mock_cache.MockCacheable
	userService  service.UserService
}

func (s *UserTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mock_repository.NewMockUserRepository(s.ctrl)
	s.tokenUseCase = mock_token.NewMockTokenUseCase(s.ctrl)
	s.cache = mock_cache.NewMockCacheable(s.ctrl)
	s.userService = service.NewUserService(s.repo, s.tokenUseCase, s.cache)
}

func TestUser(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) TestFindAll() {
	keyFindAll := "design-pattern-api:users:find-all"
	users := make([]entity.User, 0)
	marshalledData, _ := json.Marshal(users)

	s.Run("error find all user from db due to failed to get data from db", func() {
		s.cache.EXPECT().Get(keyFindAll).Return("")
		s.repo.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("error"))
		result, err := s.userService.FindAll(context.Background())
		s.NotNil(err)
		s.Nil(result)
	})
	s.Run("successfully find all user from db", func() {
		s.cache.EXPECT().Get(keyFindAll).Return("")
		s.repo.EXPECT().FindAll(gomock.Any()).Return(users, nil)
		s.cache.EXPECT().Set(keyFindAll, marshalledData, 5*time.Minute)
		result, err := s.userService.FindAll(context.Background())
		s.Nil(err)
		s.NotNil(result)
	})
	s.Run("successfully find all user from redis", func() {
		s.cache.EXPECT().Get(keyFindAll).Return(string(marshalledData))
		json.Unmarshal([]byte(marshalledData), &users)
		result, err := s.userService.FindAll(context.Background())
		s.Nil(err)
		s.NotNil(result)
	})
}

func (s *UserTestSuite) TestLogin() {
	username, password := "test", "test"
	user := new(entity.User)
	user.Password = password
	s.Run("failed to login due to error find by username", func() {
		s.repo.EXPECT().FindByUsername(gomock.Any(), username).Return(nil, errors.New("error"))
		token, err := s.userService.Login(context.Background(), username, password)
		s.NotNil(err)
		s.Empty(token)
	})
	s.Run("failed to login due to error password not match", func() {
		s.repo.EXPECT().FindByUsername(gomock.Any(), username).Return(user, nil)
		token, err := s.userService.Login(context.Background(), username, "invalid")
		s.NotNil(err)
		s.Empty(token)
	})
	s.Run("failed to login due to error generate access token", func() {
		s.repo.EXPECT().FindByUsername(gomock.Any(), username).Return(user, nil)
		s.tokenUseCase.EXPECT().GenerateAccessToken(gomock.Any()).Return("", errors.New("error"))
		token, err := s.userService.Login(context.Background(), username, password)
		s.NotNil(err)
		s.Empty(token)
	})
	s.Run("successfully login", func() {
		s.repo.EXPECT().FindByUsername(gomock.Any(), username).Return(user, nil)
		s.tokenUseCase.EXPECT().GenerateAccessToken(gomock.Any()).Return("token", nil)
		token, err := s.userService.Login(context.Background(), username, password)
		s.Nil(err)
		s.NotNil(token)
	})
}
