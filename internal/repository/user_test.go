package repository_test

import (
	"context"
	"design-pattern/internal/repository"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserTestSuite struct {
	suite.Suite
	ctrl *gomock.Controller
	db   *gorm.DB
	mock sqlmock.Sqlmock
	repo repository.UserRepository
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	db, mock, err := sqlmock.New()
	if err != nil {
		s.FailNow("failed to create mock db", err)
	}

	s.db, err = gorm.Open(postgres.New(
		postgres.Config{
			Conn: db,
		}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		s.FailNow("error openinng mock db", err)
	}

	s.mock = mock
	s.repo = repository.NewUserRepository(s.db)
}

func (s *UserTestSuite) AfterTest(string, string) {
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.FailNow("error expectations : ", err)
	}
}

func (s *UserTestSuite) TestFindAll() {
	s.Run("failed to get all users", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "public"."users"`)).
			WillReturnError(errors.New("error"))
		result, err := s.repo.FindAll(context.Background())
		s.NotNil(err)
		s.Nil(result)
	})
	s.Run("success get all users", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "public"."users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		result, err := s.repo.FindAll(context.Background())
		s.Nil(err)
		s.NotNil(result)
	})
}

func (s *UserTestSuite) TestFindByUsername() {
	username := "admin"
	s.Run("error get user by username", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "public"."users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(username, 1).
			WillReturnError(errors.New("error"))

		result, err := s.repo.FindByUsername(context.Background(), username)
		s.NotNil(err)
		s.Nil(result)
	})
	s.Run("successfully get user by username", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "public"."users" WHERE username = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(username, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		result, err := s.repo.FindByUsername(context.Background(), username)
		s.Nil(err)
		s.NotNil(result)
	})
}
