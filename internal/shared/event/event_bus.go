package event

import "context"

// DomainEvent ドメインイベントのインターフェース
type DomainEvent interface {
	EventName() string
	OccurredAt() string
	AggregateID() string
}

// EventHandler イベントハンドラーのインターフェース
type EventHandler interface {
	Handle(ctx context.Context, event DomainEvent) error
}

// EventBus イベントバスのインターフェース
type EventBus interface {
	// Register イベントハンドラーを登録
	Register(eventName string, handler EventHandler)
	
	// Publish イベントを発行
	Publish(ctx context.Context, events ...DomainEvent) error
	
	// PublishAsync イベントを非同期で発行
	PublishAsync(ctx context.Context, events ...DomainEvent) error
}

// inMemoryEventBus インメモリイベントバスの実装
type inMemoryEventBus struct {
	handlers map[string][]EventHandler
}

// NewInMemoryEventBus インメモリイベントバスのコンストラクタ
func NewInMemoryEventBus() EventBus {
	return &inMemoryEventBus{
		handlers: make(map[string][]EventHandler),
	}
}

// Register イベントハンドラーを登録
func (bus *inMemoryEventBus) Register(eventName string, handler EventHandler) {
	bus.handlers[eventName] = append(bus.handlers[eventName], handler)
}

// Publish イベントを発行
func (bus *inMemoryEventBus) Publish(ctx context.Context, events ...DomainEvent) error {
	for _, event := range events {
		handlers, exists := bus.handlers[event.EventName()]
		if !exists {
			continue
		}
		
		for _, handler := range handlers {
			if err := handler.Handle(ctx, event); err != nil {
				return err
			}
		}
	}
	return nil
}

// PublishAsync イベントを非同期で発行
func (bus *inMemoryEventBus) PublishAsync(ctx context.Context, events ...DomainEvent) error {
	go func() {
		bus.Publish(ctx, events...)
	}()
	return nil
}