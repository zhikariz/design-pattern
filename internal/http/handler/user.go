package handler

import (
	"design-pattern/internal/service"
	"design-pattern/pkg/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return UserHandler{userService}
}

func (h *UserHandler) FindAll(ctx echo.Context) error {
	users, err := h.userService.FindAll(ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError,
			response.ErrorResponse(http.StatusInternalServerError, err.Error()))
	}
	return ctx.JSON(http.StatusOK, response.SuccessResponse("successfully fetch all users", users))
}

func (h *UserHandler) Login(ctx echo.Context) error {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := ctx.Bind(&loginRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest,
			response.ErrorResponse(http.StatusBadRequest, err.Error()))
	}

	token, err := h.userService.Login(ctx.Request().Context(), loginRequest.Username, loginRequest.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, response.ErrorResponse(http.StatusUnauthorized, err.Error()))
	}

	return ctx.JSON(http.StatusOK, response.SuccessResponse("successfully login", map[string]interface{}{
		"token": token,
	}))
}
