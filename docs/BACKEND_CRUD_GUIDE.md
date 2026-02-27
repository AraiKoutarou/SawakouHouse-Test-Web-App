# Go (バックエンド) CRUD 実装・詳細解説ガイド

このガイドでは、実際のコード例を使いながら、掲示板アプリの「投稿作成（Create）」を例に、データが 4 つの層（Model -> Repository -> Service -> Handler）をどう流れるか詳しく解説します。

---

## 1. Model（モデル）：データの形を決める
`backend/model/post.go`
すべての層で共有される「データの設計図」です。

```go
type Post struct {
    ID        int       `json:"id"`         // データベースの自動採番ID
    Title     string    `json:"title"`      // 投稿タイトル
    Content   string    `json:"content"`    // 投稿内容
    Author    string    `json:"author"`     // 投稿者名
    CreatedAt time.Time `json:"created_at"` // 作成日時
}
```
*   **構造体タグ**: `` `json:"id"` `` の部分は、フロントエンドにデータを送る際の名前を指定しています。

---

## 2. Repository（リポジトリ）：DB操作の専門家
`backend/repository/post.go`
SQLを書いてデータベースと直接やり取りします。

```go
// --- 【C】Create (作成) ---
func (r *PostRepository) Create(title, content, author string) (model.Post, error) {
    var p model.Post
    err := r.db.QueryRow(
        "INSERT INTO posts (title, content, author) VALUES ($1, $2, $3) RETURNING id, title, content, author, created_at",
        title, content, author,
    ).Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.CreatedAt)
    return p, err
}

// --- 【U】Update (更新) ---
func (r *PostRepository) Update(id int, title, content string) error {
    // 値を返さない操作は Exec を使います
    _, err := r.db.Exec("UPDATE posts SET title = $1, content = $2 WHERE id = $3", title, content, id)
    return err
}

// --- 【D】Delete (削除) ---
func (r *PostRepository) Delete(id int) (bool, error) {
    result, err := r.db.Exec("DELETE FROM posts WHERE id = $1", id)
    rows, _ := result.RowsAffected() // 何行消せたか確認
    return rows > 0, err
}
```

---

## 3. Service（サービス）：ビジネスルールの守護神
`backend/service/post.go`
「保存していいデータか？」を判断する重要な場所です。

```go
func (s *PostService) Create(title, content, author string) (model.Post, error) {
    if title == "" { return model.Post{}, errors.New("タイトルは必須です") }
    return s.repo.Create(title, content, author)
}

func (s *PostService) Update(id int, title, content string) error {
    // 更新前にバリデーションを行う
    if title == "" { return errors.New("タイトルが空です") }
    return s.repo.Update(id, title, content)
}

func (s *PostService) Delete(id int) error {
    success, err := s.repo.Delete(id)
    if !success { return errors.New("見つかりませんでした") }
    return err
}
```

---

## 4. Handler（ハンドラー）：外の世界との窓口
`backend/handler/post.go`
ブラウザからのリクエストを受け取り、結果を返します。

```go
// --- 【C】Create (作成) ---
func (h *PostHandler) CreatePost(c *gin.Context) {
    var input struct {
        Title   string `json:"title" binding:"required"`
        Content string `json:"content" binding:"required"`
        Author  string `json:"author" binding:"required"`
    }
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "入力が正しくありません"})
        return
    }
    post, _ := h.svc.Create(input.Title, input.Content, input.Author)
    c.JSON(http.StatusCreated, post)
}

// --- 【U】Update (更新) ---
func (h *PostHandler) UpdatePost(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id")) // URLからIDを取得
    var input struct {
        Title   string `json:"title"`
        Content string `json:"content"`
    }
    c.ShouldBindJSON(&input)
    h.svc.Update(id, input.Title, input.Content)
    c.JSON(http.StatusOK, gin.H{"message": "更新完了"})
}

// --- 【D】Delete (削除) ---
func (h *PostHandler) DeletePost(c *gin.Context) {
    id, _ := strconv.Atoi(c.Param("id"))
    h.svc.Delete(id)
    c.JSON(http.StatusOK, gin.H{"message": "削除完了"})
}
```

---

## 5. main.go：すべての部品を組み立てる
`backend/main.go`
各層を「下から順に」組み立てて接続します。

```go
// 1. 下の層から順に作る
postRepo := repository.NewPostRepository(db.DB)
postSvc  := service.NewPostService(postRepo)
postHandler := handler.NewPostHandler(postSvc)

// 2. ルーティング（道案内）の設定
r := gin.Default()
r.POST("/posts", postHandler.CreatePost)
r.PUT("/posts/:id", postHandler.UpdatePost)    // 更新用 URL
r.DELETE("/posts/:id", postHandler.DeletePost) // 削除用 URL
```

---

## まとめ：データのバケツリレー
1.  **ブラウザ**：「このIDの投稿を消して！」(DELETE /posts/1)
2.  **Handler**：「URLからID 1を受け取ったよ。Serviceに削除を頼むね」
3.  **Service**：「RepoにID 1の削除を依頼。結果をチェックするよ」
4.  **Repository**：「SQLのDELETEを実行。消せたかどうかを返すね」
5.  **ブラウザ**：「200 OK。画面から消去するね」
