// internal/domain/post/repository.go: 永続化の契約（インターフェース）を定義します。
package post

// Repository: 投稿リポジトリのインターフェース。
type Repository interface {
	GetAll() ([]Post, error)
	GetByID(id int) (Post, error)
	Create(title, content, author string) (Post, error)
	Delete(id int) (bool, error)
}
