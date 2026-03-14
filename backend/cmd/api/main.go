// cmd/api/main.go: アプリケーションの起動とDI。
package main

import (
	"log"

	"github.com/arakou0812/backend/internal/delivery/http"
	"github.com/arakou0812/backend/internal/delivery/http/handler"
	"github.com/arakou0812/backend/internal/infrastructure/db"
	persistence "github.com/arakou0812/backend/internal/infrastructure/persistence/location"
	usecase "github.com/arakou0812/backend/internal/usecase/location"
)

func main() {
	// 1. インフラ層の初期化 (DB接続)
	database := db.Connect()
	defer database.Close()

	// 2. リポジトリの初期化とマイグレーション
	locationRepo := persistence.NewPostgresRepository(database)
	
	// スキーマ変更のため、一度テーブルをリセットする (開発用)
	database.Exec("DROP TABLE IF EXISTS locations;")
	locationRepo.Migrate()

	// 3. ユースケース層の初期化
	locationUC := usecase.NewUsecase(locationRepo)

	// 4. ハンドラ層の初期化
	locationHandler := handler.NewLocationHandler(locationUC)

	// 5. サーバー起動
	r := http.SetupRouter(locationHandler)
	
	log.Println("Map App API starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
