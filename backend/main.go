// mainだよ
package main

import (
	"github.com/arakou0812/backend/keijiban/db"
	"github.com/arakou0812/backend/keijiban/handler"
	"github.com/arakou0812/backend/keijiban/repository"
	"github.com/arakou0812/backend/keijiban/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 【1】データベースの準備
	// db/db.go で定義されている Init 関数を呼び出し、DBに接続します。
	db.Init()

	// 【2】プログラムの部品（レイヤー）を組み立てる
	// このアプリは「役割分担（レイヤー）」を明確にしています。
	// repository: データベースを直接操作する部品
	// service: データの検証や計算など、ビジネスルールを担当する部品
	// handler: ネットからのリクエストを受け取り、結果を返す窓口の部品
	postRepo := repository.NewPostRepository(db.DB)
	postSvc := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postSvc)

	// 【3】ウェブサーバー（Gin）の作成
	// Gin は Go でよく使われる、高速なウェブフレームワークです。
	r := gin.Default()

	// 【4】CORS（異なるサーバー間での通信許可）の設定
	// フロントエンド(localhost:3000)からこのバックエンド(localhost:8080)に
	// アクセスできるように許可を出しています。これがないとブラウザでエラーになります。
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	// 【5】ルーティングの設定（URLと処理の紐付け）
	// 「どのURL」に「どんな命令（GET/POST等）」が来たら「どの関数」を動かすか決めます。
	r.GET("/posts", postHandler.GetPosts)          // 投稿一覧の取得
	r.POST("/posts", postHandler.CreatePost)       // 新規投稿の作成
	r.GET("/posts/:id", postHandler.GetPost)       // 特定の投稿1件の取得
	r.DELETE("/posts/:id", postHandler.DeletePost) // 投稿の削除

	// 【6】サーバーの起動
	// 8080番ポートで待ち受けを開始します。
	r.Run(":8080")
}
