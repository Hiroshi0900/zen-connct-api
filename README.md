# Zen Connect Backend API

Auth0認証とクッキーベースのセッション管理を使用したGo製のバックエンドAPI。

## 🏗️ アーキテクチャ

- **認証**: Auth0 OAuth2 + OIDC
- **セッション管理**: 暗号化クッキー（gorilla/securecookie）
- **データベース**: PostgreSQL（ローカル開発用Docker、本番用Neon/Supabase）
- **ドメイン駆動設計**: CleanArchitecture + DDD

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