// handler層: 「窓口」としての役割を担います。
// ネットからのHTTPリクエスト（質問）を受け取り、JSONレスポンス（回答）を返すことだけに集中します。
// データベースを直接いじったり、複雑なビジネス計算は行わず、service層のメソッドに依頼します。
package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/arakou0812/backend/service"
	"github.com/gin-gonic/gin"
)

// PostHandler: 投稿に関する窓口の本体です。
// service層を使えるように、svcフィールドを持ってます。
type PostHandler struct {
	svc *service.PostService
}

// NewPostHandler: 新しい窓口を作って返す関数（コンストラクタ）です。
func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

// GetPosts: GET /posts (一覧取得) が呼ばれたときに動く関数です。
func (h *PostHandler) GetPosts(c *gin.Context) {
	// 【1】ビジネスロジック（service層）に「全部持ってきて！」と依頼します。
	posts, err := h.svc.GetAll()
	if err != nil {
		// もしエラーが起きたら、500番（サーバー内部エラー）とエラー理由をJSONで返します。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// 【2】うまくいったら、200番（OK）と、取得した投稿リストをJSONで返します。
	c.JSON(http.StatusOK, posts)
}

// GetPost: GET /posts/:id (1件取得) が呼ばれたときに動く関数です。
func (h *PostHandler) GetPost(c *gin.Context) {
	// 【1】URLの中にある「:id」の部分を抜き出して、数値(int)に変換します。
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// 数値に変換できなかったら「無効なID」として400番（リクエストエラー）を返します。
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID"})
		return
	}

	// 【2】service層に「このIDのデータを取ってきて！」と依頼します。
	post, err := h.svc.GetByID(id)
	if err != nil {
		// 見つからなかった場合(ErrNotFound)は404番（Not Found）を返します。
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		// それ以外の予期せぬエラーは500番。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 【3】成功。200番と1件の投稿データをJSONで返します。
	c.JSON(http.StatusOK, post)
}

// CreatePost: POST /posts (新規投稿作成) が呼ばれたときに動く関数です。
func (h *PostHandler) CreatePost(c *gin.Context) {
	// 【1】届いたJSONデータを受け取るための入れ物（構造体）を用意します。
	var input struct {
		Title   string `json:"title" binding:"required"`   // JSONの "title" フィールド、必須項目
		Content string `json:"content" binding:"required"` // JSONの "content" フィールド、必須項目
		Author  string `json:"author" binding:"required"`  // JSONの "author" フィールド、必須項目
	}

	// 【2】c.ShouldBindJSON でJSONの中身を input に自動でコピーします。
	if err := c.ShouldBindJSON(&input); err != nil {
		// JSONの形が変だったり、必須項目がなかったりしたら400番エラーを返します。
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 【3】service層に「新しく作って！」と依頼します。
	post, err := h.svc.Create(input.Title, input.Content, input.Author)
	if err != nil {
		// NGワードが含まれているなど、ビジネスルール違反の場合は400番。
		if errors.Is(err, service.ErrNGWord) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// その他は500番。
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 【4】成功。201番（作成済み）と、保存されたばかりのデータをJSONで返します。
	c.JSON(http.StatusCreated, post)
}

// DeletePost: DELETE /posts/:id (投稿削除) が呼ばれたときに動く関数です。
func (h *PostHandler) DeletePost(c *gin.Context) {
	// 【1】URLからIDを取り出して数値に直します。
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID"})
		return
	}

	// 【2】service層に「このIDを消して！」と依頼します。
	if err := h.svc.Delete(id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 【3】成功。200番と「削除完了メッセージ」を返します。
	c.JSON(http.StatusOK, gin.H{"message": "削除しました"})
}
