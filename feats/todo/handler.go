package todo

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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
	userID := c.Get("user_id")
	if userID == nil {
		h.lg.Error("user_id is nil")
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
	}
	var uid uuid.UUID
	var ok bool
	if uid, ok = userID.(uuid.UUID); !ok {
		h.lg.Error("user_id is not uuid", zap.Any("user_id", userID))
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "unauthorized"})
	}
	todos, err := h.repo.FindAllByUserID(c.Request().Context(), uid)
	if err != nil {
		h.lg.Error("failed to find all todos", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
	}
	return c.JSON(http.StatusOK, todos)
}
