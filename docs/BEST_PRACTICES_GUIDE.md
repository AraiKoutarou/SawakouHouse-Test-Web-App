# バックエンド ベストプラクティスガイド（初心者向け）

このガイドでは、このプロジェクトのコードを読み書きするうえで「なぜそう書くのか」を、
初心者がつまずきやすいポイントを中心に解説します。

---

## 目次

1. [エラーハンドリング](#1-エラーハンドリング)
2. [インターフェース](#2-インターフェース)
3. [命名規則](#3-命名規則)
4. [セキュリティ（SQLインジェクション）](#4-セキュリティsqlインジェクション)
5. [nilとゼロ値](#5-nilとゼロ値)
6. [defer（後始末）](#6-defer後始末)
7. [ポインタと値渡し](#7-ポインタと値渡し)
8. [よくやりがちなミス集](#8-よくやりがちなミス集)

---

## 1. エラーハンドリング

Goには例外（try/catch）がありません。エラーは**戻り値**として返します。

### 基本の形

```go
// Goの関数は「結果, エラー」の2つを返すのが慣習
result, err := someFunction()
if err != nil {
    // エラーが発生した場合の処理
    return err
}
// エラーがなければここに続く
```

### ❌ やってはいけない：エラーを無視する

```go
// ❌ エラーを _ で捨てている。バグの原因になる
result, _ := someFunction()

// ❌ エラーチェックをしていない
rows, err := r.db.Query(`SELECT ...`)
rows.Close() // errを確認する前にrowsを使っている。errがnilでないときにpanicする
```

### ✅ エラーは必ずチェックする

```go
rows, err := r.db.Query(`SELECT ...`)
if err != nil {
    return nil, err // 上の層にエラーを伝える
}
defer rows.Close()
```

### エラーの「ラップ」と「判定」

エラーに文脈（どこで何が起きたか）を追加するのを「ラップ」と言います。

```go
// ラップ：%w を使うとエラーに情報を付け加えられる
return fmt.Errorf("投稿の取得に失敗: %w", err)

// 判定：errors.Is() でラップされたエラーを判定する
// == で比較すると、ラップされたエラーは検出できないので必ず errors.Is() を使う
if errors.Is(err, domain.ErrNotFound) {
    c.JSON(http.StatusNotFound, ...)
}
```

### このプロジェクトでのエラーの流れ

```
persistence層         usecase層              handler層
sql.ErrNoRows   →   domain.ErrNotFound  →   HTTP 404
（DB固有のエラー）   （ドメインエラーに変換）  （HTTPに変換）
```

各層でエラーを「翻訳」することで、上の層がDBの詳細を知らなくて済みます。

---

## 2. インターフェース

インターフェースとは「このメソッドが使えることを保証する」という約束事です。

### なぜインターフェースを使うのか

```go
// ❌ 具体的な型に依存した場合
type PostUsecase struct {
    repo *persistence.PostgresPostRepository // PostgreSQL専用。他のDBに変えられない
}

// ✅ インターフェースに依存した場合
type PostUsecase struct {
    repo domain.PostRepository // 「これができれば何でもいい」という定義
}
```

インターフェースを使うと、**テスト時に本物のDBの代わりに偽物（モック）を使える**ようになります。

### モックを使ったテストの例

```go
// 本物のDB接続なしにテストできる
type MockPostRepository struct{}

func (m *MockPostRepository) GetAll() ([]domain.Post, error) {
    // 決まったデータを返すだけ
    return []domain.Post{
        {ID: 1, Title: "テスト投稿", Author: "テストユーザー"},
    }, nil
}

func (m *MockPostRepository) GetByID(id int) (domain.Post, error) { ... }
func (m *MockPostRepository) Create(...) (domain.Post, error) { ... }
func (m *MockPostRepository) Delete(id int) (bool, error) { ... }

// テストコード
func TestGetAll(t *testing.T) {
    uc := usecase.NewPostUsecase(&MockPostRepository{}) // モックを注入
    posts, err := uc.GetAll()
    // ...
}
```

### インターフェースを満たしているかコンパイル時に確認する

```go
// この1行を persistence/post_repository.go に書くと、
// インターフェースを満たしていない場合にコンパイルエラーで教えてくれる
var _ domain.PostRepository = (*PostgresPostRepository)(nil)
```

---

## 3. 命名規則

Goには明確な命名規則があります。これに従わないとコードレビューで指摘されます。

### 大文字と小文字の使い分け（最重要）

```go
// 大文字で始まる → パッケージ外から見える（公開）
type Post struct { ... }       // 他のパッケージから使える
func NewPostUsecase() { ... }  // 他のパッケージから使える

// 小文字で始まる → パッケージ内だけで使える（非公開）
type postHelper struct { ... } // このパッケージ内でしか使えない
func checkNGWords() { ... }    // このパッケージ内でしか使えない
var ngWords = []string{ ... }  // このパッケージ内でしか使えない
```

### 変数・関数名はキャメルケース

```go
// ✅ Goらしい書き方
postRepository := ...
getUserByID := ...
maxRetryCount := 3

// ❌ Goらしくない書き方
post_repository := ...  // スネークケースはGoでは使わない
get_user_by_id := ...
```

### 略語は全部大文字

```go
// ✅ 正しい
userID   := 1       // ID は全部大文字
httpURL  := "..."   // URL は全部大文字
jsonData := ...     // JSON は全部大文字

// ❌ 間違い
userId   := 1       // Id は不自然
httpUrl  := "..."
jsonData := ...     // これはOK（小文字始まりの場合）
```

### コンストラクタは `New〇〇` の形

```go
// ✅ Goの慣習：New + 型名
func NewPostUsecase(repo domain.PostRepository) *PostUsecase { ... }
func NewPostHandler(uc *usecase.PostUsecase) *PostHandler { ... }
```

### エラー変数は `Err〇〇` の形

```go
// ✅ 慣習通り
var ErrNotFound = errors.New("...")
var ErrNGWord   = errors.New("...")
```

---

## 4. セキュリティ（SQLインジェクション）

SQLインジェクションとは、悪意のあるSQL文を入力されてDBを不正操作される攻撃です。

### 危険なコードの例

```go
// ❌ 絶対にやってはいけない
// 例えば id に "1 OR 1=1" を渡されると全件削除される
query := "DELETE FROM posts WHERE id = " + id
r.db.Exec(query)

// ❌ fmt.Sprintf での文字列結合も同様に危険
query := fmt.Sprintf("SELECT * FROM posts WHERE author = '%s'", author)
// author に "'; DROP TABLE posts; --" を渡されたらテーブルが消える
```

### 安全なコードの例

```go
// ✅ プレースホルダー（$1, $2）を使う
// 値は必ず別の引数として渡す。SQLとデータを分離することで攻撃を防ぐ
r.db.QueryRow(
    `SELECT id, title FROM posts WHERE id = $1`, // SQL文にはプレースホルダーだけ
    id,                                           // 値は別に渡す
)

r.db.Exec(
    `INSERT INTO posts (title, content) VALUES ($1, $2)`,
    title, content,
)
```

**ルール：ユーザーから受け取った値を直接SQL文字列に埋め込まない。必ずプレースホルダーを使う。**

---

## 5. nilとゼロ値

Goにはnullの代わりに**ゼロ値**があります。

### 型ごとのゼロ値

```go
var i int     // 0
var s string  // ""（空文字列）
var b bool    // false
var p *Post   // nil（ポインタのゼロ値はnil）
var posts []Post // nil（スライスのゼロ値はnil）
```

### nilスライスと空スライスの違い

```go
var posts []Post    // nil（JSONにすると null になる）
posts = []Post{}    // 空スライス（JSONにすると [] になる）

// ❌ フロントエンドに null を返してしまう
return posts, nil // posts が nil の場合

// ✅ 1件もなかった場合でも [] を返す
if posts == nil {
    posts = []Post{}
}
return posts, nil
```

このプロジェクトの `GetAll()` でも同様の処理をしています。

### nilポインタのパニック

```go
// ❌ nilのポインタを参照するとpanicになる（実行時エラー）
var post *Post
fmt.Println(post.Title) // panic: runtime error: invalid memory address

// ✅ nilチェックをしてから使う
if post != nil {
    fmt.Println(post.Title)
}
```

---

## 6. defer（後始末）

`defer` は「関数が終わるときに必ず実行する」という仕組みです。

### なぜdeferが必要か

```go
// ❌ 途中でreturnするとClose()が呼ばれない
func (r *PostgresPostRepository) GetAll() ([]domain.Post, error) {
    rows, err := r.db.Query(`SELECT ...`)
    if err != nil {
        return nil, err // ここでreturnするとrows.Close()が呼ばれない
    }

    // 処理中にerrorが起きてreturnしても同様
    for rows.Next() {
        if err := rows.Scan(...); err != nil {
            return nil, err // ここでもrows.Close()が呼ばれない
        }
    }

    rows.Close() // ここまで辿り着かないケースがある
    return posts, nil
}
```

```go
// ✅ deferを使えば必ずClose()が呼ばれる
func (r *PostgresPostRepository) GetAll() ([]domain.Post, error) {
    rows, err := r.db.Query(`SELECT ...`)
    if err != nil {
        return nil, err
    }
    defer rows.Close() // ← 関数が終わるとき（どのreturnでも）必ず実行される

    for rows.Next() {
        if err := rows.Scan(...); err != nil {
            return nil, err // ここでreturnしてもrows.Close()が呼ばれる
        }
    }
    return posts, nil
}
```

**ルール：`rows`, ファイル, ネットワーク接続など「開いたものは必ず閉じる」リソースには `defer` を使う。**

---

## 7. ポインタと値渡し

Goではデータを「コピーして渡す（値渡し）」か「場所を渡す（ポインタ渡し）」か選べます。

### 値渡しとポインタ渡しの違い

```go
// 値渡し：コピーが渡される。元のデータは変わらない
func double(n int) int {
    n = n * 2
    return n
}
x := 5
double(x)
fmt.Println(x) // 5（変わっていない）

// ポインタ渡し：場所が渡される。元のデータが変わる
func doublePtr(n *int) {
    *n = *n * 2
}
x := 5
doublePtr(&x)
fmt.Println(x) // 10（変わっている）
```

### このプロジェクトでのポインタの使い方

```go
// 構造体はポインタで持つのが一般的（大きいデータのコピーを避ける）
type PostUsecase struct {
    repo domain.PostRepository
}

func NewPostUsecase(repo domain.PostRepository) *PostUsecase { // *PostUsecase（ポインタを返す）
    return &PostUsecase{repo: repo}
}

// メソッドはポインタレシーバで定義する
func (u *PostUsecase) GetAll() ([]domain.Post, error) { // *PostUsecase（ポインタレシーバ）
    ...
}
```

### ポインタを返す場合と値を返す場合

```go
// ✅ 「見つからなかった」を表現したい場合はポインタが便利
func GetByID(id int) (*domain.Post, error) {
    // 見つからなかった場合
    return nil, domain.ErrNotFound

    // 見つかった場合
    return &post, nil
}

// このプロジェクトでは値を返してエラーで判断している
func GetByID(id int) (domain.Post, error) {
    // 見つからなかった場合はエラーで通知
    return domain.Post{}, domain.ErrNotFound
}
```

---

## 8. よくやりがちなミス集

### ① エラーを上に返し忘れる

```go
// ❌ エラーを確認しているのに上に返していない
func (u *PostUsecase) Delete(id int) error {
    deleted, err := u.repo.Delete(id)
    if err != nil {
        fmt.Println("エラー:", err) // ログに出力するだけで握りつぶしている
    }
    return nil // 常にnilを返してしまう
}

// ✅ エラーは必ず呼び出し元に返す
func (u *PostUsecase) Delete(id int) error {
    deleted, err := u.repo.Delete(id)
    if err != nil {
        return err // 上に伝える
    }
    ...
}
```

### ② 層をまたいだ依存

```go
// ❌ handler層がdatabase/sqlをimportしている
import "database/sql"

func (h *PostHandler) GetPost(c *gin.Context) {
    if errors.Is(err, sql.ErrNoRows) { // handlerがDBを知っている
        c.JSON(http.StatusNotFound, ...)
    }
}

// ✅ handlerはドメインエラーだけを知っている
func (h *PostHandler) GetPost(c *gin.Context) {
    if errors.Is(err, domain.ErrNotFound) { // DBのことは知らない
        c.JSON(http.StatusNotFound, ...)
    }
}
```

### ③ ビジネスロジックをhandler層に書いてしまう

```go
// ❌ NGワードチェックをhandler層に書いている
func (h *PostHandler) CreatePost(c *gin.Context) {
    if strings.Contains(input.Title, "spam") { // ビジネスロジックはここに書かない
        c.JSON(http.StatusBadRequest, ...)
        return
    }
    h.uc.Create(...)
}

// ✅ ビジネスロジックはusecase層に書く
func (h *PostHandler) CreatePost(c *gin.Context) {
    post, err := h.uc.Create(input.Title, input.Content, input.Author) // usecaseに任せる
    if err != nil {
        if errors.Is(err, domain.ErrNGWord) { // エラーの種類に応じてHTTPステータスを決めるだけ
            c.JSON(http.StatusBadRequest, ...)
        }
    }
}
```

### ④ グローバル変数を使いすぎる

```go
// ❌ グローバル変数（どこからでも変更できる。テストが難しくなる）
var DB *sql.DB // db/db.go で定義していた旧パターン

func Init() {
    DB, err = sql.Open(...)
}

// ✅ 関数の戻り値として返す（このプロジェクトの現在のパターン）
func Connect() *sql.DB {
    db, err := sql.Open(...)
    return db // 使う側に渡す
}
```

### ⑤ 日本語文字数カウントのミス

```go
// ❌ len() はバイト数を返す。日本語は1文字3バイトなのでおかしな結果になる
if len("あいう") > 3 { // len("あいう") は 9 なので常にtrueになる
    return errors.New("3文字以上です")
}

// ✅ []rune に変換してからlen()を使う。runeは文字数を正しく数える
if len([]rune("あいう")) > 3 { // len([]rune("あいう")) は 3
    return errors.New("3文字以上です")
}
```

---

## まとめ

| テーマ | 最重要ポイント |
|--------|--------------|
| エラーハンドリング | エラーは無視せず必ず上に返す。`errors.Is()` で判定する |
| インターフェース | 依存は「具体的な型」ではなく「インターフェース」に向ける |
| 命名 | 公開/非公開は大文字/小文字で区別。`New〇〇` と `Err〇〇` の慣習を守る |
| セキュリティ | SQLにユーザー入力を直接埋め込まない。プレースホルダーを必ず使う |
| nil | スライスのゼロ値はnilなので空スライスに変換してから返す |
| defer | ファイルやDB結果セットは `defer` で必ず閉じる |
| ポインタ | 構造体はポインタで扱う。メソッドはポインタレシーバで定義する |
| 層の分離 | 各層の責任を守る。上の層の知識を下の層に持ち込まない |
