# アーキテクチャ設計書

## 概要

ZenConnect バックエンドは DDD（ドメイン駆動設計）とクリーンアーキテクチャの原則に基づいて設計されています。

## 設計原則

### 1. DDD原則
- **ユビキタス言語**: ビジネスドメインの用語を統一
- **境界づけられたコンテキスト**: 責務の明確な分離
- **集約**: データの整合性境界
- **ドメインイベント**: ドメインの状態変化を表現

### 2. クリーンアーキテクチャ原則
- **依存性逆転**: 上位レイヤーは下位レイヤーに依存しない
- **関心の分離**: 各層は単一の責務を持つ
- **テスタビリティ**: 各層が独立してテスト可能
- **独立性**: フレームワークやDBに依存しない

## アーキテクチャ図

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Presentation Layer                          │
├─────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐      │
│  │  Auth Handler   │  │  User Handler   │  │ Experience      │      │
│  │                 │  │                 │  │ Handler         │      │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘      │
└─────────────────────────────────────────────────────────────────────┘
                                   │
                                   ▼
┌─────────────────────────────────────────────────────────────────────┐
│                        Application Layer                            │
├─────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐      │
│  │  Auth UseCase   │  │  User UseCase   │  │ Experience      │      │
│  │                 │  │                 │  │ UseCase         │      │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘      │
│                                   │                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐      │
│  │  Auth Service   │  │  User Service   │  │ Experience      │      │
│  │                 │  │                 │  │ Service         │      │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘      │
└─────────────────────────────────────────────────────────────────────┘
                                   │
                                   ▼
┌─────────────────────────────────────────────────────────────────────┐
│                          Domain Layer                               │
├─────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐      │
│  │ Authenticated   │  │      User       │  │   Experience    │      │
│  │     User        │  │                 │  │                 │      │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘      │
│                                   │                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐      │
│  │    Session      │  │    Profile      │  │ MeditationSession│      │
│  │                 │  │                 │  │                 │      │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘      │
│                                   │                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐      │
│  │  Domain Events  │  │  Repositories   │  │   Services      │      │
│  │                 │  │  (Interfaces)   │  │                 │      │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘      │
└─────────────────────────────────────────────────────────────────────┘
                                   │
                                   ▼
┌─────────────────────────────────────────────────────────────────────┐
│                       Infrastructure Layer                          │
├─────────────────────────────────────────────────────────────────────┤
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐      │
│  │  Auth0 Client   │  │  PostgreSQL     │  │   Event Bus     │      │
│  │                 │  │  Repository     │  │                 │      │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘      │
│                                   │                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐      │
│  │ Session Store   │  │ External APIs   │  │ Unit of Work    │      │
│  │                 │  │                 │  │                 │      │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘      │
└─────────────────────────────────────────────────────────────────────┘
```

## 境界づけられたコンテキスト

### 1. 認証コンテキスト（Auth）
```
┌─────────────────────────────────────────────────────────────────────┐
│                         Auth Context                                │
├─────────────────────────────────────────────────────────────────────┤
│  Domain:                                                            │
│  - AuthenticatedUser (エンティティ)                                  │
│  - Session (値オブジェクト)                                          │
│  - Token (値オブジェクト)                                            │
│                                                                     │
│  Application:                                                       │
│  - VerifyTokenUseCase                                              │
│  - CreateSessionUseCase                                            │
│  - LogoutUseCase                                                   │
│  - AuthService (インターフェース)                                     │
│                                                                     │
│  Infrastructure:                                                   │
│  - Auth0Service (Auth0実装)                                        │
│  - SessionStore (セッション管理)                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### 2. ユーザーコンテキスト（User）
```
┌─────────────────────────────────────────────────────────────────────┐
│                         User Context                                │
├─────────────────────────────────────────────────────────────────────┤
│  Domain:                                                            │
│  - User (エンティティ・集約ルート)                                    │
│  - Profile (値オブジェクト)                                          │
│  - Email (値オブジェクト)                                            │
│  - UserRepository (インターフェース)                                  │
│                                                                     │
│  Application:                                                       │
│  - RegisterUserUseCase                                             │
│  - UpdateProfileUseCase                                            │
│  - GetUserProfileUseCase                                           │
│  - UserService (インターフェース)                                     │
│                                                                     │
│  Infrastructure:                                                   │
│  - PostgresUserRepository                                          │
│  - InMemoryUserRepository (テスト用)                                │
└─────────────────────────────────────────────────────────────────────┘
```

### 3. 体験記録コンテキスト（Experience）
```
┌─────────────────────────────────────────────────────────────────────┐
│                      Experience Context                             │
├─────────────────────────────────────────────────────────────────────┤
│  Domain:                                                            │
│  - Experience (エンティティ・集約ルート)                              │
│  - MeditationSession (値オブジェクト)                                 │
│  - EmotionalState (値オブジェクト)                                   │
│  - ExperienceRepository (インターフェース)                           │
│                                                                     │
│  Application:                                                       │
│  - CreateExperienceUseCase                                         │
│  - UpdateExperienceUseCase                                         │
│  - GetExperienceUseCase                                            │
│  - ExperienceService (インターフェース)                              │
│                                                                     │
│  Infrastructure:                                                   │
│  - PostgresExperienceRepository                                    │
└─────────────────────────────────────────────────────────────────────┘
```

## データフロー

### 1. 認証フロー
```
[Client] → [Auth Handler] → [VerifyTokenUseCase] → [AuthService] → [Auth0]
                                      ↓
[Session Store] ← [CreateSessionUseCase] ← [SessionService]
```

### 2. ユーザー登録フロー
```
[Client] → [User Handler] → [RegisterUserUseCase] → [UserService] → [UserRepository]
                                      ↓
[Event Bus] ← [Domain Events] ← [User Entity]
```

### 3. 体験記録フロー
```
[Client] → [Experience Handler] → [CreateExperienceUseCase] → [ExperienceService]
                                            ↓
[ExperienceRepository] ← [Experience Entity] → [Domain Events] → [Event Bus]
```

## イベント駆動アーキテクチャ

### ドメインイベント
```
┌─────────────────────────────────────────────────────────────────────┐
│                       Domain Events                                 │
├─────────────────────────────────────────────────────────────────────┤
│  User Context:                                                      │
│  - UserRegistered                                                   │
│  - ProfileUpdated                                                   │
│  - EmailVerified                                                    │
│                                                                     │
│  Experience Context:                                                │
│  - ExperienceCreated                                                │
│  - ExperienceUpdated                                                │
│  - ExperiencePublished                                              │
│                                                                     │
│  Auth Context:                                                      │
│  - UserLoggedIn                                                     │
│  - SessionExpired                                                   │
└─────────────────────────────────────────────────────────────────────┘
```

### イベントハンドラー
```
UserRegistered → SendWelcomeEmailHandler
ProfileUpdated → UpdateSearchIndexHandler
ExperienceCreated → NotifyFollowersHandler
```

## トランザクション管理

### Unit of Work パターン
```
┌─────────────────────────────────────────────────────────────────────┐
│                     Unit of Work Pattern                            │
├─────────────────────────────────────────────────────────────────────┤
│  UseCase.Execute() {                                                │
│    uow.WithinTx(ctx, func(ctx context.Context) error {             │
│      // 1. ドメインオブジェクトの操作                                │
│      user := domain.NewUser(...)                                    │
│      user.UpdateProfile(...)                                        │
│                                                                     │
│      // 2. リポジトリでの永続化                                      │
│      repo.Save(ctx, user)                                           │
│                                                                     │
│      // 3. イベントの収集                                           │
│      uow.CollectEvent(user.Events()...)                            │
│                                                                     │
│      return nil                                                     │
│    })                                                               │
│    // 4. トランザクションコミット後にイベント発行                     │
│  }                                                                  │
└─────────────────────────────────────────────────────────────────────┘
```

## テストアーキテクチャ

### テスト戦略
```
┌─────────────────────────────────────────────────────────────────────┐
│                        Test Strategy                                │
├─────────────────────────────────────────────────────────────────────┤
│  Unit Tests:                                                        │
│  - Domain Layer (ビジネスロジック)                                    │
│  - Application Layer (ユースケース)                                  │
│  - Infrastructure Layer (リポジトリ実装)                              │
│                                                                     │
│  Integration Tests:                                                 │
│  - API エンドポイント                                                │
│  - データベース統合                                                   │
│  - 外部サービス統合                                                   │
│                                                                     │
│  Test Doubles:                                                      │
│  - Mock (gomock生成)                                                │
│  - Stub (インメモリ実装)                                              │
│  - Fake (テスト用実装)                                               │
└─────────────────────────────────────────────────────────────────────┘
```

## 今後の拡張性

### 1. 新しいコンテキストの追加
- Notification Context (通知機能)
- Analytics Context (分析機能)
- Social Context (ソーシャル機能)

### 2. 外部サービス統合
- メール送信サービス
- プッシュ通知サービス
- 分析サービス

### 3. パフォーマンス最適化
- キャッシュ層の追加
- 読み取り専用レプリカの活用
- CDN統合