// router: 全モジュールのルートを一か所で統合します。
// 新しいモジュールが増えたら、ここに1行追加するだけです。
package router

import (
	"database/sql"

	post "github.com/arakou0812/backend/modules/post"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Setup: Ginエンジンを生成し、CORS設定と全モジュールのルートを登録して返します。
func Setup(db *sql.DB) *gin.Engine {
	r := gin.Default()

	// CORS設定: フロントエンド(localhost:3000)からのアクセスを許可します。
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// 各モジュールのルートを登録します。
	// モジュールが増えたらここに追加するだけです。
	post.NewModule(db).RegisterRoutes(r)
	// user.NewModule(db).RegisterRoutes(r)    // 将来ユーザーモジュールを追加する場合
	// comment.NewModule(db).RegisterRoutes(r) // 将来コメントモジュールを追加する場合

	return r
}
