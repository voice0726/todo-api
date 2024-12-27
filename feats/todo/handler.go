package todo

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"go.uber.org/zap"
)

type Handler struct {
	e    *echo.Echo
	lg   *zap.Logger
	repo Repository
}

func NewHandler(e *echo.Echo, lg *zap.Logger, repo Repository) *Handler {
	return &Handler{e: e, lg: lg, repo: repo}
}

func (h *Handler) Find(c echo.Context) error {
	if !authorization.IsAuthorized(c.Request().Context()) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
	}
	todo, err := h.repo.FindByID(c.Request().Context(), 1)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "todo not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}
	return c.JSON(http.StatusOK, todo)
}

func (h *Handler) FindAll(c echo.Context) error {
	if !authorization.IsAuthorized(c.Request().Context()) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
	}
	userID := authorization.UserID(c.Request().Context())
	todos, err := h.repo.FindAllByUserID(c.Request().Context(), userID)
	if err != nil {
		h.lg.Error("failed to find all todos", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}
	return c.JSON(http.StatusOK, todos)
}
