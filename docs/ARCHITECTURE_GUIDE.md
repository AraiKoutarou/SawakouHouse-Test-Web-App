# バックエンドアーキテクチャガイド (Go)

このプロジェクトのバックエンドは、Go言語の慣習に基づいた **Standard Project Layout** と **Clean Architecture** を採用しています。
この設計により、保守性、テスト容易性、および関心の分離を実現しています。

---

## 1. 全体像：設計思想

### クリーンアーキテクチャ
「依存関係は常に内側（ドメイン層）に向かう」という原則に従い、ビジネスロジックを外部ライブラリやデータベースの詳細から分離します。

### Standard Project Layout
Goコミュニティで広く認知されている構成を採用し、プロジェクトの規模が大きくなっても見通しを良くします。特に `internal/` パッケージを活用することで、外部からの不正なインポートを防ぎ、カプセル化を強制しています。

---

## 2. ディレクトリ構成

```text
backend/
├── cmd/
│   └── api/
│       └── main.go                # エントリーポイント。依存性の注入(DI)と起動のみを行う。
├── internal/                      # 外部（他プロジェクト等）からインポート不能なコアロジック
│   ├── domain/                    # 【中心】ドメイン層（ビジネスルール・エンティティ）
│   │   └── [domain]/              # 例: post
│   │       ├── entity.go          # ドメインモデル、ドメインエラー
│   │       └── repository.go      # リポジトリインターフェース（契約）
│   ├── usecase/                   # アプリケーション層（ユースケース実装）
│   │   └── [domain]/
│   │       └── usecase.go         # 具体的なビジネス手順の実装
│   ├── infrastructure/            # インフラストラクチャ層（外部依存の実装）
│   │   ├── persistence/           # データベース永続化の実装（リポジトリ実装）
│   │   └── db/                    # DB接続の初期化
│   └── delivery/                  # プレゼンテーション層（外部との接点）
│       └── http/
│           ├── handler/           # HTTPハンドラー（リクエスト/レスポンス制御）
│           └── router.go          # ルーティング定義
├── go.mod
└── go.sum
```

---

## 3. 依存関係のルール

依存関係は常に **Delivery → Usecase → Domain ← Infrastructure** の方向になります。

- **Domain層** は他の一切の内部パッケージに依存してはいけません。
- **Usecase層** は Domain層のインターフェースに依存します。
- **Infrastructure層** は Domain層で定義されたインターフェースを実装（具体化）します。
- **Delivery層** は Usecase層を呼び出します。

### なぜインターフェースを使うのか
インターフェース（`domain/repository.go`）を使うことで、テスト時に本物のデータベースを使わずに「モック」に差し替えることが容易になります。これにより、高速で安定したユニットテストが可能になります。

---

## 4. 各層の役割

### 4.1 Domain層 (`internal/domain`)
- **役割**: システムの核心となるデータ構造とインターフェースを定義します。
- **entity.go**: `Post` などの構造体。
- **repository.go**: 「DBから取得する」「保存する」といった抽象的な契約。
- **ルール**: 標準ライブラリ以外の外部パッケージ（Gin, GORMなど）を import してはいけません。

### 4.2 Usecase層 (`internal/usecase`)
- **役割**: 具体的なビジネスロジック（例：NGワードのチェック、バリデーション）を実装します。
- **特徴**: リポジトリの具体的な実装（SQLなど）は知らず、Domain層のインターフェース経由で操作します。
- **ルール**: HTTPの概念（ステータスコードやJSONタグなど）を持ち込んではいけません。

### 4.3 Infrastructure層 (`internal/infrastructure`)
- **役割**: データベースや外部APIなど、外部技術との接続を担当します。
- **persistence**: SQLを発行し、Domain層の `Repository` インターフェースを実装します。
- **特徴**: 技術的な詳細（PostgreSQL, Redis等）をこの層に閉じ込めます。

### 4.4 Delivery層 (`internal/delivery`)
- **役割**: 外部（ユーザー、フロントエンド）とのインターフェースです。
- **handler**: HTTPリクエストを受け取り、JSONをパースし、Usecaseを呼び出して結果を返します。
- **特徴**: HTTPステータスコードの管理やレスポンス形式の定義を行います。

---

## 5. データの流れ (Data Flow)

1.  **Request**: フロントエンドが `POST /posts` を送信。
2.  **Delivery (Handler)**: リクエストをパースし、`Usecase.Create()` を呼び出す。
3.  **Usecase**: ビジネスルールのチェック（NGワード等）を行い、`Repository.Create()` を呼び出す。
4.  **Infrastructure (Persistence)**: 実際に `INSERT INTO...` というSQLを実行し、DBに保存。
5.  **Response**: 結果が逆の順序で戻り、Handlerが `201 Created` を返す。

---

## 6. 新しい機能を追加する際の手順

例：`User` 機能を追加する場合

1.  **Domain作成**: `internal/domain/user/` に `entity.go` と `repository.go` を作成。
2.  **Infrastructure実装**: `internal/infrastructure/persistence/user/` にリポジトリの実装を作成。
3.  **Usecase実装**: `internal/usecase/user/` にビジネスロジックを実装。
4.  **Delivery作成**: `internal/delivery/http/handler/` にハンドラを作成し、`router.go` に登録。
5.  **DI（依存性注入）**: `cmd/api/main.go` で各層をインスタンス化し、組み立てる。
