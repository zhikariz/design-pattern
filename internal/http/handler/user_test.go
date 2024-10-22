package handler_test

import (
	"design-pattern/internal/entity"
	"design-pattern/internal/http/handler"
	mock_service "design-pattern/test/mock/service"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UserTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	userService *mock_service.MockUserService
	userHandler handler.UserHandler
}

func (s *UserTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.userService = mock_service.NewMockUserService(s.ctrl)
	s.userHandler = handler.NewUserHandler(s.userService)
}

func TestUser(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) TestFindAll() {
	users := make([]entity.User, 0)
	s.Run("error when calling service", func() {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		s.userService.EXPECT().FindAll(c.Request().Context()).Return(nil, errors.New("error"))
		s.userHandler.FindAll(c)
		s.Equal(http.StatusInternalServerError, rec.Code)
	})
	s.Run("successfully fetch all users", func() {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		s.userService.EXPECT().FindAll(c.Request().Context()).Return(users, nil)
		s.userHandler.FindAll(c)
		s.Equal(http.StatusOK, rec.Code)
	})
}

func (s *UserTestSuite) TestLogin() {
	s.Run("error login due to error binding", func() {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`invalid json format`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		s.userHandler.Login(c)
		s.Equal(http.StatusBadRequest, rec.Code)
	})
	s.Run("error login due to calling service", func() {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username": "test", "password": "test"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		s.userService.EXPECT().Login(c.Request().Context(), "test", "test").Return("", errors.New("error"))
		s.userHandler.Login(c)
		s.Equal(http.StatusUnauthorized, rec.Code)
	})
	s.Run("succesfully login", func() {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"username": "test", "password": "test"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		s.userService.EXPECT().Login(c.Request().Context(), "test", "test").Return("token", nil)
		s.userHandler.Login(c)
		s.Equal(http.StatusOK, rec.Code)
	})
}
