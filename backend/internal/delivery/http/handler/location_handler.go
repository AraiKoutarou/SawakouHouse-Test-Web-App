// internal/delivery/http/handler/location_handler.go: ハンドラの更新。
package handler

import (
	"errors"
	"net/http"

	"github.com/arakou0812/backend/internal/domain/location"
	usecase "github.com/arakou0812/backend/internal/usecase/location"
	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	uc *usecase.Usecase
}

func NewLocationHandler(uc *usecase.Usecase) *LocationHandler {
	return &LocationHandler{uc: uc}
}

func (h *LocationHandler) GetLocations(c *gin.Context) {
	locations, err := h.uc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locations)
}

func (h *LocationHandler) CreateLocation(c *gin.Context) {
	var input struct {
		PlaceID    string  `json:"place_id"`
		Title      string  `json:"title" binding:"required"`
		Address    string  `json:"address"`
		Prefecture string  `json:"prefecture"` // 追加
		Category   string  `json:"category"`
		Comment    string  `json:"comment"`
		Color      string  `json:"color"`
		Latitude   float64 `json:"latitude" binding:"required"`
		Longitude  float64 `json:"longitude" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力が不正です"})
		return
	}

	loc, err := h.uc.Create(input.PlaceID, input.Title, input.Address, input.Prefecture, input.Category, input.Comment, input.Color, input.Latitude, input.Longitude)
	if err != nil {
		if errors.Is(err, location.ErrInvalidCoordinate) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, loc)
}
