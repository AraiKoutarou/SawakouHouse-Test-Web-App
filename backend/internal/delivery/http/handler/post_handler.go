// internal/delivery/http/handler/post_handler.go: HTTPリクエストを処理します。
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arakou0812/backend/internal/domain/post"
	usecase "github.com/arakou0812/backend/internal/usecase/post"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	uc *usecase.Usecase
}

func NewPostHandler(uc *usecase.Usecase) *PostHandler {
	return &PostHandler{uc: uc}
}

func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := h.uc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID"})
		return
	}

	p, err := h.uc.GetByID(id)
	if err != nil {
		if errors.Is(err, post.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var input struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Author  string `json:"author" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := h.uc.Create(input.Title, input.Content, input.Author)
	if err != nil {
		if errors.Is(err, post.ErrNGWord) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, p)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID"})
		return
	}

	if err := h.uc.Delete(id); err != nil {
		if errors.Is(err, post.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "削除しました"})
}
