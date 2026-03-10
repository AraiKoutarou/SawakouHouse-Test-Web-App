// usecase層: 「ビジネスロジック（システムが守るべきルールや手順）」を担当します。
// 例：「不適切な言葉が含まれていたら保存しない」「文字数が多すぎたら拒否する」など。
// DBの具体的な操作（SQL）は知らず、domain.PostRepository（設計図）に「保存して」と頼むだけです。
package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/arakou0812/backend/modules/post/domain"
)

// 【ルール】投稿を禁止したいNGワードのリストです。
var ngWords = []string{"spam", "abuse"}

// PostUsecase: 投稿機能の「司令塔」となる構造体です。
type PostUsecase struct {
	// 【DI: 依存性注入】具体的な実装(PostgreSQL用など)に直接頼らず、
	// 「投稿を保存・取得できる機能」という抽象的なインターフェースを使います。
	// これにより、テスト時に本物のDBを使わずにテストをすることも可能になります。
	repo domain.PostRepository
}

// NewPostUsecase: 「司令塔」を作成するための関数です。
func NewPostUsecase(repo domain.PostRepository) *PostUsecase {
	return &PostUsecase{repo: repo}
}

// GetAll: 全投稿を取得する司令を出し、結果をそのまま返します。
func (u *PostUsecase) GetAll() ([]domain.Post, error) {
	return u.repo.GetAll()
}

// GetByID: IDを指定して1件取得する司令を出します。
func (u *PostUsecase) GetByID(id int) (domain.Post, error) {
	p, err := u.repo.GetByID(id)
	// 【エラー変換】DBから「データなし」と言われたとき、
	// そのまま返すと「DBのエラー」を外（Handlerなど）に漏らすことになります。
	// ここで「投稿が見つかりませんでした」という独自のエラーに変換して返します。
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Post{}, domain.ErrNotFound
	}
	return p, err
}

// Create: 「新しい投稿を作ってほしい」という依頼を受け取り、ルールに従って処理します。
func (u *PostUsecase) Create(title, content, author string) (domain.Post, error) {
	// 【ビジネスルール1】タイトルや本文にNGワードが含まれていないかチェックします。
	// このルールはUI（フロントエンド）でもチェックできますが、
	// セキュリティやデータの整合性を守るために、バックエンドでも必ずチェックを行います。
	if err := checkNGWords(title, content); err != nil {
		return domain.Post{}, err
	}

	// 【ビジネスルール2】タイトルの文字数制限です。
	// 日本語（マルチバイト文字）を正しく数えるために []rune に変換してからカウントしています。
	if len([]rune(title)) > 100 {
		return domain.Post{}, fmt.Errorf("タイトルは100文字以内にしてください")
	}

	// すべてのルール（チェック）を通過した「正しいデータ」だけを、
	// 最後にリポジトリ（DB担当）に渡して保存してもらいます。
	return u.repo.Create(title, content, author)
}

// Delete: 投稿を削除する司令を出し、成功したかどうかを判断します。
func (u *PostUsecase) Delete(id int) error {
	deleted, err := u.repo.Delete(id)
	if err != nil {
		return err
	}
	// もし削除対象が存在しなかった場合（削除件数が0件）はエラーを返します。
	if !deleted {
		return domain.ErrNotFound
	}
	return nil
}

// checkNGWords: 入力された文字列にNGワードが含まれているかを1つずつ調べます。
func checkNGWords(fields ...string) error {
	for _, field := range fields {
		lower := strings.ToLower(field) // 全て小文字にして大文字小文字を区別せずチェック
		for _, word := range ngWords {
			if strings.Contains(lower, word) {
				// %w を使って domain.ErrNGWord をラップして返します。
				return fmt.Errorf("%w: %q", domain.ErrNGWord, word)
			}
		}
	}
	return nil
}
