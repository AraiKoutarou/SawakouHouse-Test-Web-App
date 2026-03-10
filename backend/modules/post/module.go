// module.go: postモジュールの「入口」です。
// このファイルの役割は2つあります。
//  1. モジュール内の依存性注入（DI）を組み立てる
//  2. 外部（router）に公開するAPIを定義する（RegisterRoutesのみ）
//
// 外部からはこのファイルだけを見れば良く、内部実装の詳細は隠蔽されます。
package post

import (
	"database/sql"

	"github.com/arakou0812/backend/modules/post/handler"
	"github.com/arakou0812/backend/modules/post/persistence"
	"github.com/arakou0812/backend/modules/post/usecase"
	"github.com/gin-gonic/gin"
)

// Module: postモジュールの公開インターフェースです。
type Module struct {
	handler *handler.PostHandler
}

// NewModule: DB接続を受け取り、モジュール内の部品を組み立てて返します。
// persistence → usecase → handler の順に依存が解決されます。
func NewModule(db *sql.DB) *Module {
	repo := persistence.NewPostRepository(db)
	repo.Migrate() // このモジュールが必要とするテーブルを作成します

	uc := usecase.NewPostUsecase(repo)
	h := handler.NewPostHandler(uc)

	return &Module{handler: h}
}

// RegisterRoutes: ルーターにpostモジュールのエンドポイントを登録します。
// gin.IRouter を使うことで、*gin.Engine にも RouterGroup にも対応できます。
func (m *Module) RegisterRoutes(r gin.IRouter) {
	r.GET("/posts", m.handler.GetPosts)
	r.POST("/posts", m.handler.CreatePost)
	r.GET("/posts/:id", m.handler.GetPost)
	r.DELETE("/posts/:id", m.handler.DeletePost)
}
