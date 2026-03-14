// internal/delivery/http/router.go: ルーティングと共通設定を担当します。
package http

import (
	"github.com/arakou0812/backend/internal/delivery/http/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter: Ginエンジンを生成し、ルートを登録します。
func SetupRouter(postHandler *handler.PostHandler) *gin.Engine {
	r := gin.Default()

	// CORS設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// ルート登録
	r.GET("/posts", postHandler.GetPosts)
	r.POST("/posts", postHandler.CreatePost)
	r.GET("/posts/:id", postHandler.GetPost)
	r.DELETE("/posts/:id", postHandler.DeletePost)

	return r
}
