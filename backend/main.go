// main.go: アプリケーションのエントリポイント（起動地点）です。
// このファイルの役割は、アプリケーションが動き出すための「準備」を整えることだけです。
package main

import (
	"github.com/arakou0812/backend/router"
	"github.com/arakou0812/backend/shared/db"
)

func main() {
	// 【工程1】データベース(PostgreSQL)へ接続します。
	// shared/db/db.go で定義された接続ロジックを呼び出し、
	// アプリケーション全体で使い回すDB接続インスタンスを取得します。
	database := db.Connect()

	// 【工程2】ルーターを組み立てて起動します。
	// ここで「どのURLにアクセスしたら、どの処理が動くか」の地図を作ります。
	// CORS（フロントエンドからの通信許可）の設定や、
	// 各機能（postモジュールなど）の登録もここで行われます。
	r := router.Setup(database)

	// 【工程3】サーバーをポート8080で待ち受け状態にします。
	// これにより、ブラウザやフロントエンドからのリクエストを受け取れるようになります。
	r.Run(":8080")
}
