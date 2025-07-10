# ZenConnect Backend

瞑想アプリ「ZenConnect」のバックエンドAPIサーバーです。

## アーキテクチャ

このプロジェクトは **DDD（ドメイン駆動設計）** と **クリーンアーキテクチャ** の原則に基づいて設計されています。

### 全体構成

```
backend-app/
├── cmd/
│   └── zen-connect/         # アプリケーションエントリーポイント
├── internal/
│   ├── auth/               # 認証・認可コンテキスト
│   ├── user/               # ユーザー管理コンテキスト
│   ├── experience/         # 体験記録コンテキスト
│   └── shared/             # 共通基盤
└── migrations/             # データベースマイグレーション
```

### 境界づけられたコンテキスト（Bounded Context）

#### 1. 認証コンテキスト（`auth/`）
- **責務**: 認証・認可の管理
- **外部連携**: Auth0による認証プロバイダー
- **主な機能**:
  - トークン検証
  - セッション管理
  - ログイン/ログアウト処理

#### 2. ユーザーコンテキスト（`user/`）
- **責務**: ユーザー情報の管理
- **主な機能**:
  - ユーザー登録・更新
  - プロフィール管理
  - ユーザー情報の取得

#### 3. 体験記録コンテキスト（`experience/`）
- **責務**: 瞑想体験の記録と管理
- **主な機能**:
  - 体験記録の作成・更新・削除
  - 公開・非公開設定
  - 体験記録の検索・一覧

### 各コンテキストの内部構造

各境界づけられたコンテキストは以下の4層で構成されています：

```
context/
├── domain/                 # ドメイン層
│   ├── entities.go        # エンティティ
│   ├── value_objects.go   # 値オブジェクト
│   ├── repository.go      # リポジトリインターフェース
│   ├── events.go          # ドメインイベント
│   └── errors.go          # ドメインエラー
├── application/            # アプリケーション層
│   ├── usecase/           # ユースケース
│   ├── service/           # アプリケーションサービス
│   └── dto/               # データ転送オブジェクト
├── interfaces/             # インターフェース層
│   ├── routes.go          # ルーティング
│   ├── handlers.go        # HTTPハンドラー
│   └── middleware.go      # ミドルウェア
└── infrastructure/         # インフラストラクチャ層
    ├── postgres/          # PostgreSQL実装
    ├── auth0/             # Auth0実装
    └── external/          # 外部API実装
```

### 層の責務

#### ドメイン層（domain/）
- **責務**: ビジネスルールとドメインロジック
- **特徴**:
  - 外部依存なし
  - 純粋なビジネスロジック
  - エンティティ、値オブジェクト、ドメインサービス

#### アプリケーション層（application/）
- **usecase/**: 単一のユーザー操作を表現
- **service/**: 複数のユースケースで共通するビジネスロジック
- **dto/**: データ転送オブジェクト

#### インターフェース層（interfaces/）
- **責務**: 外部からのアクセス制御
- **特徴**:
  - HTTPハンドラー
  - ルーティング定義
  - 入力検証・出力変換

#### インフラストラクチャ層（infrastructure/）
- **責務**: 外部システムとの連携
- **特徴**:
  - データベース実装
  - 外部API実装
  - 具体的な技術への依存

### 依存関係の方向

```
interfaces → application → domain
    ↓
infrastructure → domain (インターフェース実装)
```

- **依存性逆転の原則**: 上位レイヤーは下位レイヤーに依存しない
- **インターフェース分離**: 抽象に依存し、具象に依存しない

### 共通基盤（shared/）

#### Event Bus（`shared/event/`）
- ドメインイベントの発行・購読
- コンテキスト間の疎結合な連携

#### Unit of Work（`shared/uow/`）
- トランザクション管理
- DB更新とイベント発行の整合性保証

## 🚀 セットアップ

### 1. 環境変数の設定

```bash
cp .env.example .env
# .envファイルを編集して適切な値を設定
```

### 2. PostgreSQLの起動（ローカル開発）

```bash
docker-compose up -d postgres
```

### 3. アプリケーションの実行

```bash
go run cmd/zen-connect/main.go
```

## 📋 API エンドポイント

### 認証フロー

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/auth/login` | Auth0のUniversal Loginにリダイレクト |
| GET | `/auth/callback` | Auth0からのコールバック処理 |
| GET | `/auth/logout` | ログアウト（セッションクリア + Auth0ログアウト） |
| GET | `/auth/me` | 現在のユーザー情報取得（未ログイン時は404） |

### ヘルスチェック

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | 基本的なヘルスチェック |
| GET | `/health/protected` | 認証が必要なヘルスチェック |
| GET | `/health/db` | データベース接続チェック |

### API情報

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/routes` | 利用可能なエンドポイント一覧 |
| GET | `/api/auth/login-url` | Auth0ログインURL取得（AJAX用） |

## 🔐 認証フロー

1. **ログイン**: `/auth/login` → Auth0 Universal Login
2. **コールバック**: Auth0 → `/auth/callback` → ユーザー作成/更新 → セッション作成
3. **API利用**: 暗号化クッキーでセッション検証
4. **ログアウト**: `/auth/logout` → セッションクリア → Auth0ログアウト

## 🗄️ データベーススキーマ

### usersテーブル

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    auth0_user_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(255),
    bio TEXT,
    profile_image_url TEXT,
    email_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    verified_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

## 🧪 テスト

```bash
# ユニットテスト
go test ./...

# カバレッジ付きテスト
go test -cover ./...
```

## 🚢 デプロイ

### Neon PostgreSQL

1. [Neon](https://neon.tech)でプロジェクト作成
2. CONNECTION_STRINGを`.env`に設定
3. マイグレーション実行

### Supabase PostgreSQL

1. [Supabase](https://supabase.com)でプロジェクト作成
2. PostgreSQL CONNECTION_STRINGを`.env`に設定
3. マイグレーション実行

## 🔧 開発

### 依存関係の追加

```bash
go mod tidy
```

### ログレベルの設定

```bash
export LOG_LEVEL=debug  # debug, info, warn, error
```

## 📝 環境変数

| Variable | Description | Example |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL接続文字列 | `postgres://user:pass@host/db` |
| `AUTH0_DOMAIN` | Auth0ドメイン | `your-domain.auth0.com` |
| `AUTH0_CLIENT_ID` | Auth0クライアントID | `abc123...` |
| `AUTH0_CLIENT_SECRET` | Auth0クライアントシークレット | `xyz789...` |
| `AUTH0_AUDIENCE` | Auth0 API識別子 | `https://api.zenconnect.com` |
| `SESSION_SECRET` | セッション暗号化キー（32バイト） | `your-32-byte-secret-key-here-1234` |
| `FRONTEND_URL` | フロントエンドURL | `http://localhost:3000` |

## 🏷️ バージョン

- Go: 1.21+
- PostgreSQL: 16+
- Auth0: Universal Login 2.0