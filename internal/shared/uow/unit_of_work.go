package uow

import (
	"context"
	"database/sql"
	"zen-connect/internal/shared/event"
)

// UnitOfWork Unit of Workパターンのインターフェース
type UnitOfWork interface {
	// WithinTx トランザクション内で処理を実行
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
	
	// CollectEvent ドメインイベントを収集
	CollectEvent(e event.DomainEvent)
	
	// GetTx 現在のトランザクションを取得
	GetTx(ctx context.Context) *sql.Tx
}

// unitOfWork Unit of Workの実装
type unitOfWork struct {
	db         *sql.DB
	eventBus   event.EventBus
	events     []event.DomainEvent
}

// NewUnitOfWork Unit of Workのコンストラクタ
func NewUnitOfWork(db *sql.DB, eventBus event.EventBus) UnitOfWork {
	return &unitOfWork{
		db:       db,
		eventBus: eventBus,
		events:   make([]event.DomainEvent, 0),
	}
}

// WithinTx トランザクション内で処理を実行
func (uow *unitOfWork) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := uow.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	
	// コンテキストにトランザクションを埋め込む
	ctx = context.WithValue(ctx, "tx", tx)
	
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	
	// 処理を実行
	if err := fn(ctx); err != nil {
		tx.Rollback()
		return err
	}
	
	// トランザクションをコミット
	if err := tx.Commit(); err != nil {
		return err
	}
	
	// イベントを発行
	if len(uow.events) > 0 {
		if err := uow.eventBus.Publish(ctx, uow.events...); err != nil {
			return err
		}
		uow.events = make([]event.DomainEvent, 0)
	}
	
	return nil
}

// CollectEvent ドメインイベントを収集
func (uow *unitOfWork) CollectEvent(e event.DomainEvent) {
	uow.events = append(uow.events, e)
}

// GetTx 現在のトランザクションを取得
func (uow *unitOfWork) GetTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value("tx").(*sql.Tx); ok {
		return tx
	}
	return nil
}