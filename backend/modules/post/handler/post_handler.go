// handler層: HTTPリクエストを受け取り、JSONレスポンスを返す窓口です。
// ビジネスロジックは持たず、usecase層に処理を依頼します。
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arakou0812/backend/modules/post/domain"
	"github.com/arakou0812/backend/modules/post/usecase"
	"github.com/gin-gonic/gin"
)

// PostHandler: 投稿に関するHTTPハンドラーの本体です。
type PostHandler struct {
	uc *usecase.PostUsecase
}

// NewPostHandler: Usecaseを受け取ってHandlerを生成します。
func NewPostHandler(uc *usecase.PostUsecase) *PostHandler {
	return &PostHandler{uc: uc}
}

// GetPosts: GET /posts - 投稿一覧を返します。
func (h *PostHandler) GetPosts(c *gin.Context) {
	posts, err := h.uc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// GetPost: GET /posts/:id - 指定IDの投稿を1件返します。
func (h *PostHandler) GetPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID"})
		return
	}

	post, err := h.uc.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

// CreatePost: POST /posts - 新しい投稿を作成します。
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

	post, err := h.uc.Create(input.Title, input.Content, input.Author)
	if err != nil {
		if errors.Is(err, domain.ErrNGWord) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// DeletePost: DELETE /posts/:id - 指定IDの投稿を削除します。
func (h *PostHandler) DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID"})
		return
	}

	if err := h.uc.Delete(id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "削除しました"})
}
