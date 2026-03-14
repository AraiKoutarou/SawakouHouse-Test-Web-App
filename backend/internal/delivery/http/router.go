// internal/delivery/http/router.go: ルーティング定義。
package http

import (
	"github.com/arakou0812/backend/internal/delivery/http/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter: Ginエンジンを生成し、ルートを登録します。
func SetupRouter(locationHandler *handler.LocationHandler) *gin.Engine {
	r := gin.Default()

	// CORS設定: フロントエンド(localhost:3000)からのアクセスを許可します。
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// ルート登録
	r.GET("/locations", locationHandler.GetLocations)
	r.POST("/locations", locationHandler.CreateLocation)

	return r
}
