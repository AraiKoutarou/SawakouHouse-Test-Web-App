# SawakouHouse-Test-Web-App

Go（バックエンド）+ Next.js（フロントエンド）の掲示板アプリ（練習用リポジトリ）。

---

## 技術スタック

| 項目 | 内容 |
|------|------|
| バックエンド | Go 1.25 + Gin + github.com/lib/pq |
| DB | PostgreSQL 16（Docker） |
| フロントエンド | Next.js (App Router) + TypeScript + Tailwind CSS + axios |
| モジュール名 | `github.com/arakou0812/backend` |

---

## ディレクトリ構成

```
SawakouHouse-Test-Web-App/
├── CLAUDE.md
├── GEMINI.md
├── docker-compose.yml
├── docs/
│   ├── ARCHITECTURE_GUIDE.md   # アーキテクチャ全体の説明
│   ├── BEST_PRACTICES_GUIDE.md # ベストプラクティス（初心者向け）
│   ├── BACKEND_CRUD_GUIDE.md   # バックエンドCRUD実装ガイド
│   ├── FRONTEND_CRUD_GUIDE.md  # フロントエンドCRUD実装ガイド
│   ├── DOCKER_GUIDE.md         # Docker環境構築ガイド
│   ├── GIT_GUIDE.md            # Git操作ガイド
│   ├── GIT_CONVENTION_GUIDE.md # Gitコミット規約
│   ├── HTTP_GUIDE.md           # HTTP基礎知識
│   └── README.md               # プロジェクト概要
├── backend/
│   ├── main.go                 # エントリポイント（DB接続 + サーバー起動のみ）
│   ├── go.mod
│   ├── shared/db/db.go         # DB接続（*sql.DB を返す）
│   ├── router/router.go        # 全モジュールのルート統合 + CORS設定
│   └── modules/
│       └── post/               # 投稿モジュール
│           ├── module.go       # DI組み立て・RegisterRoutes
│           ├── domain/
│           │   ├── post.go     # Postエンティティ・ドメインエラー
│           │   └── repository.go # PostRepositoryインターフェース
│           ├── usecase/
│           │   └── post_usecase.go # ビジネスロジック
│           ├── handler/
│           │   └── post_handler.go # HTTPハンドラー
│           └── persistence/
│               └── post_repository.go # PostgreSQL実装・Migrate()
└── frontend/
    └── src/app/
        ├── page.tsx            # 投稿一覧（トップ）
        └── posts/
            ├── [id]/page.tsx   # 投稿詳細・削除
            └── new/page.tsx    # 新規投稿フォーム
```

---

## アーキテクチャ

**モジュラーモノリス + クリーンアーキテクチャ**を採用。

### 依存の方向

```
handler → usecase → domain ← persistence
```

- `domain` 層は何にも依存しない（標準ライブラリのみ）
- 各モジュールは `module.go` を入口として外部に公開する
- モジュール間の直接 import は禁止。公開インターフェース経由で通信する

### 層の責任

| 層 | 役割 | 書いてはいけないもの |
|----|------|-------------------|
| domain | エンティティ・エラー・インターフェース定義 | 外部パッケージのimport |
| persistence | SQL によるDB操作 | ビジネスロジック |
| usecase | バリデーション・エラー変換 | Gin・HTTPステータスコード |
| handler | HTTPの受付・JSONレスポンス | SQL・ビジネスロジック |

---

## API エンドポイント

| メソッド | パス | 説明 |
|---------|------|------|
| GET | /posts | 投稿一覧取得 |
| POST | /posts | 新規投稿（title, content, author が必須） |
| GET | /posts/:id | 投稿1件取得 |
| DELETE | /posts/:id | 投稿削除 |

CORS は `http://localhost:3000` のみ許可。

---

## 開発環境のセットアップ

```bash
# 1. PostgreSQL 起動
docker compose up -d

# 2. バックエンド起動（http://localhost:8080）
cd backend
go run main.go

# 3. フロントエンド起動（http://localhost:3000）
cd frontend
npm run dev
```

---

## コーディング規約

### Go

- エラーは必ず上の層に返す。`_` で無視しない
- `errors.Is()` でエラーを判定する（`==` は使わない）
- SQLにユーザー入力を直接埋め込まない。プレースホルダー（`$1`, `$2`）を必ず使う
- DB結果セット（`rows`）は `defer rows.Close()` で必ず閉じる
- 日本語の文字数カウントは `len([]rune(s))` を使う（`len(s)` はバイト数）
- 構造体はポインタで扱い、メソッドはポインタレシーバで定義する
- コンストラクタは `NewXxx` の形で定義する
- ビジネスエラーは `ErrXxx` の形で `domain` 層に定義する

### 新しいモジュールを追加するとき

1. `modules/<name>/` 以下に `domain/`, `usecase/`, `handler/`, `persistence/` を作成
2. `module.go` で DI を組み立て `RegisterRoutes` を定義
3. `router/router.go` に1行追加するだけで完了

---

## 注意事項

- DB接続情報は環境変数（`DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`）で変更可能。デフォルト値あり
- テーブルはサーバー起動時に `persistence.Migrate()` で自動作成される
- グローバル変数は使わない。`db.Connect()` で `*sql.DB` を返して値として受け渡す
