// cmd/api/main.go: アプリケーションの起動とDI。
package main

import (
	"log"

	"github.com/arakou0812/backend/internal/delivery/http"
	"github.com/arakou0812/backend/internal/delivery/http/handler"
	"github.com/arakou0812/backend/internal/infrastructure/db"
	persistence "github.com/arakou0812/backend/internal/infrastructure/persistence/post"
	usecase "github.com/arakou0812/backend/internal/usecase/post"
)

func main() {
	// 1. インフラ層の初期化
	database := db.Connect()
	defer database.Close()

	// 2. リポジトリの初期化
	postRepo := persistence.NewPostgresRepository(database)
	postRepo.Migrate()

	// 3. ユースケース層の初期化
	postUC := usecase.NewUsecase(postRepo)

	// 4. ハンドラ層の初期化
	postHandler := handler.NewPostHandler(postUC)

	// 5. サーバー起動
	r := http.SetupRouter(postHandler)
	
	log.Println("Server starting on :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
