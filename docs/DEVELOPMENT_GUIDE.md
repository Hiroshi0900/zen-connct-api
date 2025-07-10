# 開発ガイドライン

## 概要

このドキュメントは、ZenConnect バックエンドの開発に関するガイドラインを提供します。

## 1. プロジェクト構造

### ディレクトリ構造の規則

#### 境界づけられたコンテキスト
```
internal/
├── {context}/           # 境界づけられたコンテキスト
│   ├── domain/         # ドメイン層
│   ├── application/    # アプリケーション層
│   │   ├── usecase/   # ユースケース
│   │   ├── service/   # サービス
│   │   └── dto/       # DTO
│   ├── interfaces/    # インターフェース層
│   └── infrastructure/ # インフラ層
└── shared/            # 共通基盤
```

#### ファイル命名規則
- **エンティティ**: `{entity}.go` (例: `user.go`)
- **値オブジェクト**: `{value_object}.go` (例: `email.go`)
- **リポジトリ**: `{entity}_repository.go` (例: `user_repository.go`)
- **ユースケース**: `{action}_{entity}.go` (例: `register_user.go`)
- **サービス**: `{entity}_service.go` (例: `user_service.go`)
- **ハンドラー**: `{entity}_handler.go` (例: `user_handler.go`)

## 2. 層の実装ガイドライン

### ドメイン層（domain/）

#### エンティティの実装
```go
// user.go
package domain

import (
    "time"
    "zen-connect/internal/shared/event"
)

// User ユーザーエンティティ（集約ルート）
type User struct {
    id       string
    email    Email
    profile  Profile
    events   []event.DomainEvent
    createdAt time.Time
    updatedAt time.Time
}

// NewUser コンストラクタ
func NewUser(email Email, profile Profile) *User {
    user := &User{
        id:        generateID(),
        email:     email,
        profile:   profile,
        createdAt: time.Now(),
        updatedAt: time.Now(),
    }
    
    // ドメインイベントの発行
    user.addEvent(NewUserRegistered(user.id, user.email.String()))
    
    return user
}

// ビジネスロジック
func (u *User) UpdateProfile(profile Profile) {
    u.profile = profile
    u.updatedAt = time.Now()
    u.addEvent(NewProfileUpdated(u.id, profile))
}

// プライベートメソッド
func (u *User) addEvent(event event.DomainEvent) {
    u.events = append(u.events, event)
}
```

#### 値オブジェクトの実装
```go
// email.go
package domain

import (
    "errors"
    "regexp"
    "strings"
)

// Email メールアドレス値オブジェクト
type Email struct {
    value string
}

// NewEmail コンストラクタ
func NewEmail(value string) (Email, error) {
    if value == "" {
        return Email{}, errors.New("email cannot be empty")
    }
    
    value = strings.ToLower(strings.TrimSpace(value))
    
    if !isValidEmail(value) {
        return Email{}, errors.New("invalid email format")
    }
    
    return Email{value: value}, nil
}

// String 文字列表現
func (e Email) String() string {
    return e.value
}

// プライベートメソッド
func isValidEmail(email string) bool {
    regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    return regex.MatchString(email)
}
```

#### リポジトリインターフェース
```go
// user_repository.go
package domain

import "context"

//go:generate mockgen -source=user_repository.go -destination=../../mocks/user_repository_mock.go -package=mocks

// UserRepository ユーザーリポジトリインターフェース
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
    FindByEmail(ctx context.Context, email Email) (*User, error)
    FindByAuth0UserID(ctx context.Context, auth0UserID string) (*User, error)
    Delete(ctx context.Context, id string) error
}
```

### アプリケーション層（application/）

#### ユースケースの実装
```go
// register_user.go
package usecase

import (
    "context"
    "zen-connect/internal/user/application/dto"
    "zen-connect/internal/user/application/service"
    "zen-connect/internal/user/domain"
)

// RegisterUserUseCase ユーザー登録ユースケース
type RegisterUserUseCase struct {
    userService service.UserService
}

// NewRegisterUserUseCase コンストラクタ
func NewRegisterUserUseCase(userService service.UserService) *RegisterUserUseCase {
    return &RegisterUserUseCase{
        userService: userService,
    }
}

// Execute ユースケース実行
func (uc *RegisterUserUseCase) Execute(ctx context.Context, req *dto.RegisterUserRequest) (*dto.RegisterUserResponse, error) {
    // 1. 入力検証
    if err := uc.validateInput(req); err != nil {
        return nil, err
    }
    
    // 2. ドメインオブジェクトの作成
    email, err := domain.NewEmail(req.Email)
    if err != nil {
        return nil, err
    }
    
    profile := domain.NewProfile(req.DisplayName, "", "")
    
    // 3. ビジネスロジックの実行
    user := domain.NewUser(email, profile)
    
    // 4. 永続化
    if err := uc.userService.Save(ctx, user); err != nil {
        return nil, err
    }
    
    // 5. レスポンスの作成
    return &dto.RegisterUserResponse{
        UserID:    user.ID(),
        Email:     user.Email().String(),
        CreatedAt: user.CreatedAt(),
    }, nil
}

// プライベートメソッド
func (uc *RegisterUserUseCase) validateInput(req *dto.RegisterUserRequest) error {
    if req.Email == "" {
        return errors.New("email is required")
    }
    return nil
}
```

#### サービスの実装
```go
// user_service.go
package service

import (
    "context"
    "zen-connect/internal/user/domain"
)

//go:generate mockgen -source=user_service.go -destination=../../../mocks/user_service_mock.go -package=mocks

// UserService ユーザーサービス
type UserService interface {
    Save(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id string) (*domain.User, error)
    FindByEmail(ctx context.Context, email domain.Email) (*domain.User, error)
}

// userServiceImpl 実装
type userServiceImpl struct {
    userRepo domain.UserRepository
}

// NewUserService コンストラクタ
func NewUserService(userRepo domain.UserRepository) UserService {
    return &userServiceImpl{
        userRepo: userRepo,
    }
}

// Save ユーザー保存
func (s *userServiceImpl) Save(ctx context.Context, user *domain.User) error {
    return s.userRepo.Save(ctx, user)
}
```

### インターフェース層（interfaces/）

#### HTTPハンドラーの実装
```go
// user_handler.go
package interfaces

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "zen-connect/internal/user/application/dto"
    "zen-connect/internal/user/application/usecase"
)

// UserHandler ユーザーハンドラー
type UserHandler struct {
    registerUserUseCase *usecase.RegisterUserUseCase
}

// NewUserHandler コンストラクタ
func NewUserHandler(registerUserUseCase *usecase.RegisterUserUseCase) *UserHandler {
    return &UserHandler{
        registerUserUseCase: registerUserUseCase,
    }
}

// RegisterUser ユーザー登録
func (h *UserHandler) RegisterUser(c echo.Context) error {
    // 1. リクエストバインディング
    var req dto.RegisterUserRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "invalid request format",
        })
    }
    
    // 2. バリデーション
    if err := c.Validate(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }
    
    // 3. ユースケース実行
    response, err := h.registerUserUseCase.Execute(c.Request().Context(), &req)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{
            "error": err.Error(),
        })
    }
    
    // 4. レスポンス返却
    return c.JSON(http.StatusCreated, response)
}
```

### インフラストラクチャ層（infrastructure/）

#### リポジトリの実装
```go
// postgres_user_repository.go
package postgres

import (
    "context"
    "database/sql"
    "zen-connect/internal/user/domain"
)

// postgresUserRepository PostgreSQL実装
type postgresUserRepository struct {
    db *sql.DB
}

// NewPostgresUserRepository コンストラクタ
func NewPostgresUserRepository(db *sql.DB) domain.UserRepository {
    return &postgresUserRepository{
        db: db,
    }
}

// Save ユーザー保存
func (r *postgresUserRepository) Save(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (id, email, display_name, bio, profile_image_url, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        ON CONFLICT (id) DO UPDATE SET
            email = EXCLUDED.email,
            display_name = EXCLUDED.display_name,
            bio = EXCLUDED.bio,
            profile_image_url = EXCLUDED.profile_image_url,
            updated_at = EXCLUDED.updated_at
    `
    
    _, err := r.db.ExecContext(ctx, query,
        user.ID(),
        user.Email().String(),
        user.Profile().DisplayName(),
        user.Profile().Bio(),
        user.Profile().ProfileImageURL(),
        user.CreatedAt(),
        user.UpdatedAt(),
    )
    
    return err
}
```

## 3. テストの実装ガイドライン

### Given/When/Then パターン
```go
// user_test.go
package domain_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "zen-connect/internal/user/domain"
)

func TestUser_UpdateProfile(t *testing.T) {
    t.Run("正常なプロフィール更新", func(t *testing.T) {
        // Given
        email, _ := domain.NewEmail("test@example.com")
        profile := domain.NewProfile("Old Name", "", "")
        user := domain.NewUser(email, profile)
        
        newProfile := domain.NewProfile("New Name", "new-bio", "new-image.jpg")
        
        // When
        user.UpdateProfile(newProfile)
        
        // Then
        assert.Equal(t, "New Name", user.Profile().DisplayName())
        assert.Equal(t, "new-bio", user.Profile().Bio())
        assert.Equal(t, "new-image.jpg", user.Profile().ProfileImageURL())
        assert.NotEmpty(t, user.Events())
    })
}
```

### gomockを使用したテスト
```go
// user_service_test.go
package service_test

import (
    "context"
    "testing"
    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"
    "zen-connect/internal/mocks"
    "zen-connect/internal/user/application/service"
    "zen-connect/internal/user/domain"
)

func TestUserService_Save(t *testing.T) {
    t.Run("ユーザー保存成功", func(t *testing.T) {
        // Given
        ctrl := gomock.NewController(t)
        defer ctrl.Finish()
        
        mockRepo := mocks.NewMockUserRepository(ctrl)
        userService := service.NewUserService(mockRepo)
        
        email, _ := domain.NewEmail("test@example.com")
        profile := domain.NewProfile("Test User", "", "")
        user := domain.NewUser(email, profile)
        
        mockRepo.EXPECT().
            Save(gomock.Any(), user).
            Return(nil).
            Times(1)
        
        // When
        err := userService.Save(context.Background(), user)
        
        // Then
        assert.NoError(t, err)
    })
}
```

## 4. コーディング規則

### 命名規則
- **パッケージ**: 小文字、短く、意味のある名前
- **エクスポート**: PascalCase
- **非エクスポート**: camelCase
- **定数**: UPPER_SNAKE_CASE
- **インターフェース**: 動詞 + "er" または 名詞

### エラーハンドリング
```go
// カスタムエラーの定義
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidEmail = errors.New("invalid email format")
)

// エラーラッピング
func (s *userService) FindByID(ctx context.Context, id string) (*domain.User, error) {
    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to find user: %w", err)
    }
    return user, nil
}
```

### ログ出力
```go
// 構造化ログ
log.WithFields(log.Fields{
    "user_id": user.ID(),
    "action":  "register",
}).Info("User registered successfully")

// エラーログ
log.WithError(err).WithFields(log.Fields{
    "user_id": userID,
}).Error("Failed to save user")
```

## 5. 依存関係の管理

### DIコンテナの使用
```go
// main.go
func main() {
    // インフラストラクチャ
    db := initDB()
    
    // リポジトリ
    userRepo := postgres.NewPostgresUserRepository(db)
    
    // サービス
    userService := service.NewUserService(userRepo)
    
    // ユースケース
    registerUserUseCase := usecase.NewRegisterUserUseCase(userService)
    
    // ハンドラー
    userHandler := interfaces.NewUserHandler(registerUserUseCase)
    
    // ルーターセットアップ
    e := echo.New()
    userHandler.SetupRoutes(e)
    
    e.Start(":8080")
}
```

## 6. パフォーマンスのベストプラクティス

### データベースアクセス
```go
// バッチ処理の実装
func (r *postgresUserRepository) SaveBatch(ctx context.Context, users []*domain.User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    stmt, err := tx.PrepareContext(ctx, "INSERT INTO users ...")
    if err != nil {
        return err
    }
    defer stmt.Close()
    
    for _, user := range users {
        if _, err := stmt.ExecContext(ctx, user.ID(), user.Email().String()); err != nil {
            return err
        }
    }
    
    return tx.Commit()
}
```

### メモリ使用量の最適化
```go
// スライスの事前確保
users := make([]*domain.User, 0, expectedSize)

// 大きなオブジェクトの適切な解放
defer func() {
    users = nil
}()
```

## 7. セキュリティガイドライン

### 入力検証
```go
// バリデーション
func (h *UserHandler) RegisterUser(c echo.Context) error {
    var req dto.RegisterUserRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "invalid request format",
        })
    }
    
    // カスタムバリデーション
    if err := validateUserInput(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }
    
    // ...
}
```

### ログでの機密情報マスキング
```go
// 機密情報のマスキング
type MaskedUser struct {
    ID    string `json:"id"`
    Email string `json:"email,omitempty"`
}

func (u *User) MarshalJSON() ([]byte, error) {
    masked := MaskedUser{
        ID:    u.ID(),
        Email: maskEmail(u.Email().String()),
    }
    return json.Marshal(masked)
}
```

## 8. 継続的インテグレーション

### テストの実行
```bash
# 全テストの実行
go test ./...

# カバレッジの計測
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# ベンチマーク
go test -bench=. ./...
```

### 静的解析
```bash
# go vet
go vet ./...

# golint
golint ./...

# staticcheck
staticcheck ./...
```

## 9. ドキュメンテーション

### GoDoc コメント
```go
// User はユーザーを表すエンティティです。
// これは集約ルートとしても機能します。
type User struct {
    // ...
}

// NewUser は新しいユーザーを作成します。
// 
// パラメータ:
//   - email: 有効なメールアドレス
//   - profile: ユーザープロフィール
// 
// 戻り値:
//   - *User: 作成されたユーザー
func NewUser(email Email, profile Profile) *User {
    // ...
}
```

## 10. 今後の課題

### 技術的改善
- [ ] パフォーマンステストの充実
- [ ] メトリクス収集の実装
- [ ] 分散トレーシングの導入
- [ ] 障害テストの実装

### 開発体験の向上
- [ ] 開発環境の自動化
- [ ] APIドキュメントの自動生成
- [ ] コードジェネレーターの導入
- [ ] 開発者向けツールの整備