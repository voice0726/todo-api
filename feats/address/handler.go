package address

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	e    *echo.Echo
	repo Repository
}

func NewHandler(e *echo.Echo, repo Repository) *Handler {
	return &Handler{e: e, repo: repo}
}

func (h *Handler) Find(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "id is required")
	}
	var uid uuid.UUID
	var err error
	if uid, err = uuid.Parse(id); err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}
	address, err := h.repo.FindByID(c.Request().Context(), uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, "address not found")
		}
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, address)
}

func (h *Handler) FindAllByUserID(c echo.Context) error {
	uid := c.Get("user_id")
	if uid == nil {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	if _, ok := uid.(uuid.UUID); !ok {
		return c.JSON(http.StatusUnauthorized, "unauthorized")
	}
	addresses, err := h.repo.FindAllByUserID(c.Request().Context(), uid.(uuid.UUID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "internal server error")
	}
	return c.JSON(http.StatusOK, addresses)
}
