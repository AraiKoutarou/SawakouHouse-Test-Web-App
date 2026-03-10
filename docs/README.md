# SawakouHouse-Test-Web-App

Go（バックエンド）+ Next.js（フロントエンド）の掲示板アプリ（練習用リポジトリ）

## 技術スタック

| 項目 | 内容 |
|------|------|
| バックエンド | Go + Gin |
| DB | PostgreSQL (Docker) |
| フロントエンド | Next.js (App Router) + Tailwind CSS |
| 通信 | REST API (axios) |

## ディレクトリ構成

```
SawakouHouse-Test-Web-App/
├── backend/                   # Go API サーバー
│   ├── main.go                # エントリポイント
│   ├── go.mod
│   ├── shared/db/             # DB接続（共通）
│   ├── router/                # ルーティング統合
│   └── modules/               # 機能モジュール群
│       └── post/              # 投稿モジュール
│           ├── module.go      # DI組み立て・公開API
│           ├── domain/        # エンティティ・インターフェース
│           ├── usecase/       # ビジネスロジック
│           ├── handler/       # HTTPハンドラー
│           └── persistence/   # DB操作実装
├── frontend/                  # Next.js アプリ
│   └── src/app/
│       ├── page.tsx           # 投稿一覧
│       └── posts/
│           ├── [id]/page.tsx  # 投稿詳細・削除
│           └── new/page.tsx   # 新規投稿フォーム
├── docs/                      # ドキュメント
├── docker-compose.yml
└── README.md
```

> アーキテクチャの詳細は [docs/ARCHITECTURE_GUIDE.md](./ARCHITECTURE_GUIDE.md) を参照してください。

## 起動方法

### 1. PostgreSQL 起動

```bash
docker compose up -d
```

### 2. バックエンド起動

```bash
cd backend
go run main.go
# → http://localhost:8080
```

### 3. フロントエンド起動

```bash
cd frontend
npm run dev
# → http://localhost:3000
```

## API エンドポイント

| メソッド | パス | 説明 |
|---------|------|------|
| GET | /posts | 投稿一覧取得 |
| POST | /posts | 新規投稿 |
| GET | /posts/:id | 投稿詳細 |
| DELETE | /posts/:id | 投稿削除 |
