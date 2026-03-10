# バックエンドアーキテクチャガイド

このプロジェクトのバックエンドは **モジュラーモノリス + クリーンアーキテクチャ** を採用しています。
このドキュメントでは、なぜこの設計を選んだのか・各層の役割・書き方のベストプラクティスを説明します。

> エラーハンドリング・インターフェース・命名規則・セキュリティなど Go の基礎的なベストプラクティスは [BEST_PRACTICES_GUIDE.md](./BEST_PRACTICES_GUIDE.md) を参照してください。

---

## 目次

1. [全体像：なぜこの設計なのか](#1-全体像なぜこの設計なのか)
2. [ディレクトリ構成](#2-ディレクトリ構成)
3. [依存関係のルール](#3-依存関係のルール)
4. [各層の役割とベストプラクティス](#4-各層の役割とベストプラクティス)
   - [domain層](#domain層)
   - [persistence層](#persistence層)
   - [usecase層](#usecase層)
   - [handler層](#handler層)
   - [module.go（DI組み立て）](#modulego-di組み立て)
5. [モジュール間の通信ルール](#5-モジュール間の通信ルール)
6. [新しいモジュールの追加手順](#6-新しいモジュールの追加手順)

---

## 1. 全体像：なぜこの設計なのか

### モジュラーモノリスとは

1つのアプリ（モノリス）の中を、機能ごとのモジュール（post, user, comment など）に分けた構造です。

```
┌─────────────────────────────────┐
│         1つのGoバイナリ          │
│  ┌────────┐  ┌────────┐         │
│  │  post  │  │  user  │   ...   │
│  │module  │  │module  │         │
│  └────────┘  └────────┘         │
└─────────────────────────────────┘
```

**メリット**

| メリット | 説明 |
|----------|------|
| 境界が明確 | モジュール外の内部実装に触れられない |
| チーム分担しやすい | モジュール単位で担当を分けられる |
| デプロイがシンプル | 1つのバイナリで動く |
| 将来の分割が容易 | 各moduleをそのままマイクロサービスに切り出せる |

### クリーンアーキテクチャとは

「外側の層が内側の層に依存する」という1方向のルールで層を分ける設計です。

```
      ┌──────────────────┐
      │    handler層      │  ← HTTPの世界
      │  ┌────────────┐  │
      │  │ usecase層  │  │  ← ビジネスロジック
      │  │ ┌────────┐ │  │
      │  │ │domain層│ │  │  ← 核心（依存なし）
      │  │ └────────┘ │  │
      │  └────────────┘  │
      └──────────────────┘
      persistence層（外側からdomainに依存）
```

> **鉄則**：内側の層は外側の層を知らない。

---

## 2. ディレクトリ構成

```
backend/
├── main.go                  # エントリポイント（DB接続 + サーバー起動のみ）
├── go.mod
│
├── shared/                  # モジュール横断の共通処理
│   └── db/
│       └── db.go            # DB接続（*sql.DBを返す）
│
├── router/
│   └── router.go            # 全モジュールのルート統合 + CORS設定
│
└── modules/
    └── post/                # postモジュール（1機能 = 1ディレクトリ）
        ├── module.go        # ★入口：DIの組み立て・外部への公開API
        ├── domain/
        │   ├── post.go      # エンティティ・ドメインエラー
        │   └── repository.go # リポジトリインターフェース
        ├── usecase/
        │   └── post_usecase.go # ビジネスロジック
        ├── handler/
        │   └── post_handler.go # HTTPハンドラー
        └── persistence/
            └── post_repository.go # DB操作の実装
```

---

## 3. 依存関係のルール

```
handler ──→ usecase ──→ domain ←── persistence
```

- `handler` は `usecase` だけを呼ぶ
- `usecase` は `domain`（インターフェース）だけを呼ぶ
- `persistence` は `domain` のインターフェースを実装する
- **`domain` は何にも依存しない**

### なぜインターフェースを使うのか

```go
// ❌ 悪い例：具体的な実装に依存
type PostUsecase struct {
    repo *persistence.PostgresPostRepository // PostgreSQL専用になってしまう
}

// ✅ 良い例：インターフェースに依存
type PostUsecase struct {
    repo domain.PostRepository // どんな実装でも差し替え可能
}
```

インターフェースを使うと、テスト時にモック（偽物のDB）に差し替えられます。

```go
// テスト用のモック
type MockPostRepository struct{}
func (m *MockPostRepository) GetAll() ([]domain.Post, error) {
    return []domain.Post{{ID: 1, Title: "テスト投稿"}}, nil
}

// テストコードで差し替え
uc := usecase.NewPostUsecase(&MockPostRepository{})
```

---

## 4. 各層の役割とベストプラクティス

---

### domain層

**役割**：アプリの核心。エンティティとビジネスエラー、リポジトリインターフェースを定義する。

**ルール**：外部パッケージを一切インポートしない（`time` など標準ライブラリはOK）。

#### エンティティ

```go
// ✅ 良い例
type Post struct {
    ID        int
    Title     string
    Content   string
    Author    string
    CreatedAt time.Time
}
```

**ベストプラクティス**

- ✅ エンティティにドメインロジックのメソッドを持たせてよい
  ```go
  func (p *Post) IsExpired() bool {
      return time.Since(p.CreatedAt) > 30*24*time.Hour // 30日経過で期限切れ
  }
  ```
- ✅ ビジネスエラーはdomain層に定義する（usecase/serviceに書かない）
  ```go
  var (
      ErrNotFound = errors.New("投稿が見つかりません")
      ErrNGWord   = errors.New("NGワードが含まれています")
  )
  ```
- ❌ `database/sql` や `gin` などの外部パッケージをimportしない
- ❌ JSON変換や HTTP の概念を持ち込まない

#### リポジトリインターフェース

```go
// ✅ 良い例：使う側（usecase）の視点でメソッドを定義する
type PostRepository interface {
    GetAll() ([]Post, error)
    GetByID(id int) (Post, error)
    Create(title, content, author string) (Post, error)
    Delete(id int) (bool, error)
}
```

**ベストプラクティス**

- ✅ インターフェースは「使う側」（usecase）のニーズに合わせて設計する
- ✅ メソッドは少なくシンプルに（Goの格言：「小さいインターフェースを好む」）
- ❌ DB固有の型（`*sql.Rows` など）をインターフェースに含めない

---

### persistence層

**役割**：`domain.PostRepository` インターフェースの具体的な実装。SQLを書く唯一の場所。

#### 実装例

```go
type PostgresPostRepository struct {
    db *sql.DB
}

// domain.PostRepository を実装していることをコンパイル時に保証する
var _ domain.PostRepository = (*PostgresPostRepository)(nil)
```

> `var _ domain.PostRepository = (*PostgresPostRepository)(nil)` を書くと、
> インターフェースを満たしていない場合にコンパイルエラーになります。

**ベストプラクティス**

- ✅ **プレースホルダー（`$1`, `$2`）を必ず使う**（SQLインジェクション対策）
  ```go
  // ✅ 安全
  r.db.Query(`SELECT * FROM posts WHERE id = $1`, id)

  // ❌ 危険：ユーザー入力を直接埋め込んでいる
  r.db.Query(`SELECT * FROM posts WHERE id = ` + id)
  ```
- ✅ `rows.Close()` は必ず `defer` で閉じる
  ```go
  rows, err := r.db.Query(`SELECT ...`)
  if err != nil { return nil, err }
  defer rows.Close() // ← 必須
  ```
- ✅ `Migrate()` メソッドでテーブル作成を管理する（モジュールが自分のテーブルを責任を持つ）
- ❌ ビジネスロジック（NGワードチェックなど）を書かない
- ❌ `sql.ErrNoRows` をそのまま上に返さず、`domain.ErrNotFound` への変換はusecase層で行う

---

### usecase層

**役割**：ビジネスロジック（ルール・チェック・計算）の置き場所。

#### 実装例

```go
type PostUsecase struct {
    repo domain.PostRepository // インターフェースに依存
}

func (u *PostUsecase) Create(title, content, author string) (domain.Post, error) {
    // ビジネスルールのチェックはここで行う
    if err := checkNGWords(title, content); err != nil {
        return domain.Post{}, err
    }
    if len([]rune(title)) > 100 {
        return domain.Post{}, fmt.Errorf("タイトルは100文字以内にしてください")
    }
    return u.repo.Create(title, content, author)
}
```

**ベストプラクティス**

- ✅ **DBエラー（`sql.ErrNoRows`）をドメインエラーに変換する**
  ```go
  // ✅ 良い例：usecase層でエラーを翻訳する
  func (u *PostUsecase) GetByID(id int) (domain.Post, error) {
      p, err := u.repo.GetByID(id)
      if errors.Is(err, sql.ErrNoRows) {
          return domain.Post{}, domain.ErrNotFound // ← ドメインエラーに変換
      }
      return p, err
  }
  ```
- ✅ **Ginなどのフレームワークをimportしない**（usecase層はHTTPを知らない）
- ✅ 複数のリポジトリを使う処理（複数モジュールをまたぐ調整）はusecase層に書く
  ```go
  // 例：コメント追加時に投稿の存在確認が必要な場合
  type AddCommentUsecase struct {
      commentRepo domain.CommentRepository
      postChecker PostChecker // postモジュールの公開インターフェース
  }
  ```
- ❌ `c *gin.Context` のような引数を受け取らない
- ❌ JSONのパースやHTTPステータスコードを意識しない

---

### handler層

**役割**：HTTPリクエストの受付と、JSONレスポンスの返却のみ。

#### 実装例

```go
func (h *PostHandler) CreatePost(c *gin.Context) {
    // 1. リクエストのパース
    var input struct {
        Title   string `json:"title" binding:"required"`
        Content string `json:"content" binding:"required"`
        Author  string `json:"author" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 2. usecaseに処理を依頼
    post, err := h.uc.Create(input.Title, input.Content, input.Author)

    // 3. エラーをHTTPステータスコードに変換
    if err != nil {
        if errors.Is(err, domain.ErrNGWord) {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 4. レスポンスを返す
    c.JSON(http.StatusCreated, post)
}
```

**ベストプラクティス**

- ✅ **ビジネスロジックをhandler層に書かない**（「NGワードかどうか」の判断はusecase層）
- ✅ **ドメインエラー → HTTPステータスコードの変換はhandler層の責任**
  ```go
  // ✅ エラーの種類に応じてステータスコードを変える
  if errors.Is(err, domain.ErrNotFound) {
      c.JSON(http.StatusNotFound, ...)    // 404
  } else {
      c.JSON(http.StatusInternalServerError, ...) // 500
  }
  ```
- ✅ `errors.Is()` を使ってエラーを判定する（`==` は使わない）
- ❌ `database/sql` や直接DBを操作するコードを書かない
- ❌ ビジネスエラーを `errors.New()` で新たに作らない（domain層に定義済みのものを使う）

---

### module.go（DI組み立て）

**役割**：モジュール内の依存性注入（DI）の組み立てと、外部への公開APIの定義。

```go
func NewModule(db *sql.DB) *Module {
    // persistence → usecase → handler の順に組み立てる
    repo := persistence.NewPostRepository(db)
    repo.Migrate()
    uc := usecase.NewPostUsecase(repo)
    h  := handler.NewPostHandler(uc)
    return &Module{handler: h}
}

// 外部（router）に公開するのはこのメソッドだけ
func (m *Module) RegisterRoutes(r gin.IRouter) {
    r.GET("/posts", m.handler.GetPosts)
    // ...
}
```

**ベストプラクティス**

- ✅ `*gin.Engine` ではなく `gin.IRouter` を引数にする（RouterGroupにも対応できる）
- ✅ 外部に公開するメソッドを最小限にする（`RegisterRoutes` のみが基本）
- ✅ 他のモジュールが参照する場合は、このファイルに公開メソッドを追加する
  ```go
  // 他モジュールへ公開するメソッドの例（方法1）
  func (m *Module) PostExists(id int) bool {
      // postの存在確認だけを公開
  }
  ```

---

## 5. モジュール間の通信ルール

モジュール間で直接importするのは **禁止** です。

```go
// ❌ 絶対にやってはいけない
import "github.com/arakou0812/backend/modules/post/domain"  // commentモジュールからpostをimport
```

### 方法：公開インターフェース経由

```
modules/comment/usecase ──→ modules/post/module.go（公開メソッド）
```

```go
// Step 1: commentのusecase側でインターフェースを定義する
// modules/comment/usecase/add_comment.go
type PostChecker interface {
    PostExists(id int) bool
}

type AddCommentUsecase struct {
    commentRepo domain.CommentRepository
    postChecker PostChecker // postモジュールの型に依存しない
}

// Step 2: postのmodule.goに公開メソッドを追加する
// modules/post/module.go
func (m *Module) PostExists(id int) bool { ... }

// Step 3: router.goでpostModuleをcommentModuleに渡す
// router/router.go
postMod    := post.NewModule(db)
commentMod := comment.NewModule(db, postMod) // postModを渡す
```

---

## 6. 新しいモジュールの追加手順

`user` モジュールを追加する例：

```bash
# 1. ディレクトリ作成
mkdir -p backend/modules/user/{domain,usecase,handler,persistence}
```

```
# 2. 以下の順でファイルを作成する
modules/user/domain/user.go          # エンティティ・エラー定義
modules/user/domain/repository.go    # UserRepositoryインターフェース
modules/user/persistence/user_repository.go # DB実装
modules/user/usecase/user_usecase.go # ビジネスロジック
modules/user/handler/user_handler.go # HTTPハンドラー
modules/user/module.go               # DI組み立て・RegisterRoutes
```

```go
// 3. router/router.go に1行追加するだけで完了
import user "github.com/arakou0812/backend/modules/user"

func Setup(db *sql.DB) *gin.Engine {
    // ...
    post.NewModule(db).RegisterRoutes(r)
    user.NewModule(db).RegisterRoutes(r) // ← 追加
    return r
}
```

---

## まとめ：各層の責任早見表

| 層 | 役割 | 書いていいもの | 書いてはいけないもの |
|---|---|---|---|
| **domain** | エンティティ・エラー・インターフェース | 標準ライブラリのみ | 外部パッケージのimport |
| **persistence** | SQLによるDB操作 | SQL, `database/sql` | ビジネスロジック |
| **usecase** | ビジネスロジック | バリデーション, エラー変換 | Gin, HTTPステータスコード |
| **handler** | HTTPの受付・応答 | Gin, HTTPステータスコード | SQL, ビジネスロジック |
| **module.go** | DIの組み立て | NewXxx() の呼び出し | ビジネスロジック |
