// mainだよ
package main

import (
	"github.com/arakou0812/backend/db"
	"github.com/arakou0812/backend/handler"
	"github.com/arakou0812/backend/repository"
	"github.com/arakou0812/backend/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// DB初期化
	db.Init()

	// 依存関係の組み立て: repository → service → handler の順に注入する
	postRepo := repository.NewPostRepository(db.DB)
	postSvc := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postSvc)

	r := gin.Default()

	// CORS設定（フロントエンド localhost:3000 からのアクセスを許可）
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// ルーティング
	r.GET("/posts", postHandler.GetPosts)
	r.POST("/posts", postHandler.CreatePost)
	r.GET("/posts/:id", postHandler.GetPost)
	r.DELETE("/posts/:id", postHandler.DeletePost)

	r.Run(":8080")
}
