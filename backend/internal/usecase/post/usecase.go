// internal/usecase/post/usecase.go: ビジネスロジックを実装します。
package post

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/arakou0812/backend/internal/domain/post"
)

var ngWords = []string{"spam", "abuse"}

// Usecase: 投稿機能の司令塔。
type Usecase struct {
	repo post.Repository
}

// NewUsecase: コンストラクタ。
func NewUsecase(repo post.Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) GetAll() ([]post.Post, error) {
	return u.repo.GetAll()
}

func (u *Usecase) GetByID(id int) (post.Post, error) {
	p, err := u.repo.GetByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		return post.Post{}, post.ErrNotFound
	}
	return p, err
}

func (u *Usecase) Create(title, content, author string) (post.Post, error) {
	if err := checkNGWords(title, content); err != nil {
		return post.Post{}, err
	}

	if len([]rune(title)) > 100 {
		return post.Post{}, fmt.Errorf("タイトルは100文字以内にしてください")
	}

	return u.repo.Create(title, content, author)
}

func (u *Usecase) Delete(id int) error {
	deleted, err := u.repo.Delete(id)
	if err != nil {
		return err
	}
	if !deleted {
		return post.ErrNotFound
	}
	return nil
}

func checkNGWords(fields ...string) error {
	for _, field := range fields {
		lower := strings.ToLower(field)
		for _, word := range ngWords {
			if strings.Contains(lower, word) {
				return fmt.Errorf("%w: %q", post.ErrNGWord, word)
			}
		}
	}
	return nil
}
