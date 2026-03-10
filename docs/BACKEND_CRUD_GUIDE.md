# Go (バックエンド) CRUD 実装・詳細解説ガイド

このガイドでは、実際のコード例を使いながら、掲示板アプリの「投稿作成（Create）」を例に、データが各層をどう流れるか詳しく解説します。

現在のバックエンドは **モジュラーモノリス + クリーンアーキテクチャ** を採用しており、`modules/post/` の中に以下の4層があります。

> アーキテクチャの設計思想については [ARCHITECTURE_GUIDE.md](./ARCHITECTURE_GUIDE.md) を参照してください。

---

## 層の対応表（旧 → 新）

| 旧の名前 | 新の名前 | ファイルの場所 |
|---------|---------|---------------|
| Model | **domain** | `modules/post/domain/post.go` |
| Repository | **persistence** | `modules/post/persistence/post_repository.go` |
| Service | **usecase** | `modules/post/usecase/post_usecase.go` |
| Handler | **handler** | `modules/post/handler/post_handler.go` |
| main.go での DI | **module.go** | `modules/post/module.go` |

---

## 1. domain（ドメイン）：データの形とルールを決める

`backend/modules/post/domain/post.go`

すべての層で共有される「データの設計図」とビジネスエラーを定義します。
**他のどのパッケージにも依存しない**ことが最大の特徴です。

```go
type Post struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Content   string    `json:"content"`
    Author    string    `json:"author"`
    CreatedAt time.Time `json:"created_at"`
}

// ビジネスエラーはdomain層で定義する
var (
    ErrNotFound = errors.New("投稿が見つかりません")
    ErrNGWord   = errors.New("NGワードが含まれています")
)
```

`backend/modules/post/domain/repository.go`

「何ができるか」だけを定義したインターフェースです。
「どうやるか（SQL）」は書きません。

```go
type PostRepository interface {
    GetAll() ([]Post, error)
    GetByID(id int) (Post, error)
    Create(title, content, author string) (Post, error)
    Delete(id int) (bool, error)
}
```

*   **構造体タグ**: `` `json:"id"` `` の部分は、フロントエンドにデータを送る際の名前を指定しています。
*   **インターフェース**: usecase層はこのインターフェース越しにDBを操作するため、テスト時にモック（偽のDB）に差し替えられます。

---

## 2. persistence（パーシステンス）：DB操作の専門家

`backend/modules/post/persistence/post_repository.go`

`domain.PostRepository` インターフェースを実装します。
SQLを書く唯一の場所です。

```go
// --- 【C】Create (作成) ---
func (r *PostgresPostRepository) Create(title, content, author string) (domain.Post, error) {
    var p domain.Post
    err := r.db.QueryRow(
        // RETURNING で保存直後のデータ（IDや作成日時）を取得できる
        `INSERT INTO posts (title, content, author) VALUES ($1, $2, $3)
         RETURNING id, title, content, author, created_at`,
        title, content, author,
    ).Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt)
    return p, err
}

// --- 【D】Delete (削除) ---
func (r *PostgresPostRepository) Delete(id int) (bool, error) {
    result, err := r.db.Exec(`DELETE FROM posts WHERE id = $1`, id)
    if err != nil {
        return false, err
    }
    rows, _ := result.RowsAffected() // 実際に何行削除されたか確認
    return rows > 0, err
}
```

> **注意**：`$1`, `$2` のプレースホルダーを必ず使ってください。
> 文字列を直接繋げると **SQLインジェクション** という重大なセキュリティ脆弱性になります。

---

## 3. usecase（ユースケース）：ビジネスロジックの守護神

`backend/modules/post/usecase/post_usecase.go`

「保存していいデータか？」を判断する重要な場所です。
`domain.PostRepository` **インターフェース**に依存しており、SQLの書き方は知りません。

```go
func (u *PostUsecase) Create(title, content, author string) (domain.Post, error) {
    // 【チェック1】NGワードが含まれていないか？
    if err := checkNGWords(title, content); err != nil {
        return domain.Post{}, err
    }
    // 【チェック2】タイトルが100文字以内か？
    if len([]rune(title)) > 100 {
        return domain.Post{}, fmt.Errorf("タイトルは100文字以内にしてください")
    }
    // すべてのチェックを通過したらpersistence層に保存を依頼する
    return u.repo.Create(title, content, author)
}

func (u *PostUsecase) GetByID(id int) (domain.Post, error) {
    p, err := u.repo.GetByID(id)
    // DBのエラー(sql.ErrNoRows)をドメインエラーに変換するのがusecase層の役割
    if errors.Is(err, sql.ErrNoRows) {
        return domain.Post{}, domain.ErrNotFound
    }
    return p, err
}

func (u *PostUsecase) Delete(id int) error {
    deleted, err := u.repo.Delete(id)
    if err != nil { return err }
    if !deleted { return domain.ErrNotFound }
    return nil
}
```

---

## 4. handler（ハンドラー）：外の世界との窓口

`backend/modules/post/handler/post_handler.go`

ブラウザからのリクエストを受け取り、JSONで結果を返します。
ビジネスロジックは持たず、usecase層に処理を依頼します。

```go
// --- 【C】Create (作成) ---
func (h *PostHandler) CreatePost(c *gin.Context) {
    // 1. リクエストのJSONをパースする
    var input struct {
        Title   string `json:"title" binding:"required"`
        Content string `json:"content" binding:"required"`
        Author  string `json:"author" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    // 2. usecaseに処理を依頼する
    post, err := h.uc.Create(input.Title, input.Content, input.Author)
    if err != nil {
        // 3. ドメインエラー → HTTPステータスコードに変換する
        if errors.Is(err, domain.ErrNGWord) {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    // 4. 成功レスポンスを返す
    c.JSON(http.StatusCreated, post)
}

// --- 【D】Delete (削除) ---
func (h *PostHandler) DeletePost(c *gin.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "無効なID"})
        return
    }
    if err := h.uc.Delete(id); err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "削除しました"})
}
```

---

## 5. module.go：部品の組み立てと公開

`backend/modules/post/module.go`

モジュール内の依存性注入（DI）を一か所で行います。
外部（router）に公開するのは `RegisterRoutes` のみです。

```go
func NewModule(db *sql.DB) *Module {
    // persistence → usecase → handler の順に組み立てる
    repo := persistence.NewPostRepository(db)
    repo.Migrate() // このモジュールが使うテーブルを作成
    uc   := usecase.NewPostUsecase(repo)
    h    := handler.NewPostHandler(uc)
    return &Module{handler: h}
}

func (m *Module) RegisterRoutes(r gin.IRouter) {
    r.GET("/posts", m.handler.GetPosts)
    r.POST("/posts", m.handler.CreatePost)
    r.GET("/posts/:id", m.handler.GetPost)
    r.DELETE("/posts/:id", m.handler.DeletePost)
}
```

`backend/router/router.go` では、全モジュールをここで登録するだけです。

```go
func Setup(db *sql.DB) *gin.Engine {
    r := gin.Default()
    r.Use(cors.New(...)) // CORS設定

    post.NewModule(db).RegisterRoutes(r)
    // user.NewModule(db).RegisterRoutes(r) ← 新モジュールはここに1行追加するだけ

    return r
}
```

---

## まとめ：データのバケツリレー

投稿作成（POST /posts）を例に、データの流れを追います。

```
ブラウザ
  │ POST /posts {"title":"...", "content":"...", "author":"..."}
  ▼
handler（CreatePost）
  │ JSONをパース → usecase.Create() を呼ぶ
  ▼
usecase（Create）
  │ NGワードチェック・文字数チェック → repo.Create() を呼ぶ
  ▼
persistence（Create）
  │ INSERT INTO posts ... を実行 → 保存されたデータを返す
  ▼
usecase → handler
  │ 結果をそのまま上に返す
  ▼
ブラウザ
  201 Created {"id":1, "title":"...", ...}
```

各層は**自分の責任だけ**を持ち、他の層の内部実装を知りません。これにより、例えばDBをPostgreSQLからMySQLに変えても `persistence` 層だけを修正すれば済みます。
