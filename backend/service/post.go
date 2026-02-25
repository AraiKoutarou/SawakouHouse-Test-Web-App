// service層: ビジネスロジックを担当する
// 「NGワードチェック」「文字数制限」などのルールをここに集める。
// DBの詳細は知らず、repositoryのメソッドを呼ぶだけ。
package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/arakou0812/backend/model"
	"github.com/arakou0812/backend/repository"
)

// ErrNotFound: 投稿が存在しないことを表すエラー
var ErrNotFound = errors.New("投稿が見つかりません")

// ErrNGWord: NGワードが含まれている場合のエラー
var ErrNGWord = errors.New("NGワードが含まれています")

// ngWords: 禁止ワード一覧
var ngWords = []string{"spam", "abuse"}

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) GetAll() ([]model.Post, error) {
	return s.repo.GetAll()
}

func (s *PostService) GetByID(id int) (model.Post, error) {
	p, err := s.repo.GetByID(id)
	if errors.Is(err, sql.ErrNoRows) {
		return model.Post{}, ErrNotFound
	}
	return p, err
}

// Create: NGワードチェック → 文字数チェック → 保存
func (s *PostService) Create(title, content, author string) (model.Post, error) {
	if err := checkNGWords(title, content); err != nil {
		return model.Post{}, err
	}
	if len([]rune(title)) > 100 {
		return model.Post{}, fmt.Errorf("タイトルは100文字以内にしてください")
	}
	return s.repo.Create(title, content, author)
}

func (s *PostService) Delete(id int) error {
	deleted, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrNotFound
	}
	return nil
}

// checkNGWords: タイトル・本文にNGワードが含まれていないか確認する
func checkNGWords(fields ...string) error {
	for _, field := range fields {
		lower := strings.ToLower(field)
		for _, word := range ngWords {
			if strings.Contains(lower, word) {
				return fmt.Errorf("%w: %q", ErrNGWord, word)
			}
		}
	}
	return nil
}
