// service層: 「ビジネスロジック（計算やチェック）」を担当します。
// 「NGワードが入っていないか？」「文字数が長すぎないか？」といったルールをここに集約します。
// データベースの具体的な書き方は知らず、repository層に依頼します。
package service

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/arakou0812/backend/model"
	"github.com/arakou0812/backend/repository"
)

// 【カスタムエラーの定義】
// プログラムの中で「何が原因で失敗したか」を分かりやすくするためにエラーを定義します。
var ErrNotFound = errors.New("投稿が見つかりません")
var ErrNGWord = errors.New("NGワードが含まれています")

// 禁止ワードのリスト。実務ではDBや外部サービスで管理することもあります。
var ngWords = []string{"spam", "abuse"}

// PostService: サービス層の本体です。
type PostService struct {
	repo *repository.PostRepository
}

// NewPostService: repositoryを受け取って新しいServiceを作ります。
func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

// GetAll: 全ての投稿を取得します。
func (s *PostService) GetAll() ([]model.Post, error) {
	// そのままrepositoryに依頼します。
	return s.repo.GetAll()
}

// GetByID: 指定されたIDの投稿を取得します。
func (s *PostService) GetByID(id int) (model.Post, error) {
	p, err := s.repo.GetByID(id)
	// repositoryから「データがないよ(sql.ErrNoRows)」と言われたら、
	// サービス独自のエラー(ErrNotFound)に変換して返します。
	if errors.Is(err, sql.ErrNoRows) {
		return model.Post{}, ErrNotFound
	}
	return p, err
}

// Create: 「チェック」をしてから「保存」します。
func (s *PostService) Create(title, content, author string) (model.Post, error) {
	// 【チェック1】NGワードが含まれていないか？
	if err := checkNGWords(title, content); err != nil {
		return model.Post{}, err
	}

	// 【チェック2】タイトルが長すぎないか？
	// 日本語（マルチバイト文字）を正しく数えるために []rune に変換しています。
	if len([]rune(title)) > 100 {
		return model.Post{}, fmt.Errorf("タイトルは100文字以内にしてください")
	}

	// すべてのチェックを通過したら、ようやくrepositoryに保存を依頼します。
	return s.repo.Create(title, content, author)
}

// Delete: 削除を依頼し、結果を確認します。
func (s *PostService) Delete(id int) error {
	deleted, err := s.repo.Delete(id)
	if err != nil {
		return err
	}
	// もし削除対象が見つからなかったら、エラーを返します。
	if !deleted {
		return ErrNotFound
	}
	return nil
}

// checkNGWords: プライベートな関数（小文字で始まる）で、共通のチェック処理を行います。
func checkNGWords(fields ...string) error {
	for _, field := range fields {
		lower := strings.ToLower(field) // 全て小文字にして比較しやすくします
		for _, word := range ngWords {
			if strings.Contains(lower, word) {
				// %w を使うことで、エラーを「ラップ」して詳細情報を付け加えられます。
				return fmt.Errorf("%w: %q", ErrNGWord, word)
			}
		}
	}
	return nil
}
