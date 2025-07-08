package logger

import (
	"context"

	"go.uber.org/zap"
)

// Context keys for storing logging-related data
type contextKey string

const (
	contextKeyLogger        contextKey = "logger"
	contextKeyRequestID     contextKey = "request_id"
	contextKeyCorrelationID contextKey = "correlation_id"
	contextKeySessionID     contextKey = "session_id"
	contextKeyUserID        contextKey = "user_id"
	contextKeyAuth0UserID   contextKey = "auth0_user_id"
	contextKeyUserEmail     contextKey = "user_email"
	contextKeyUserName      contextKey = "user_name"
	contextKeyTraceID       contextKey = "trace_id"
	contextKeySpanID        contextKey = "span_id"
	contextKeyComponent     contextKey = "component"
)

// LogContext wraps context with logging functionality
type LogContext struct {
	ctx    context.Context
	logger *zap.Logger
	masker *Masker
	fields []zap.Field
}

// NewLogContext creates a new log context
func NewLogContext(ctx context.Context, logger *zap.Logger, masker *Masker) *LogContext {
	return &LogContext{
		ctx:    ctx,
		logger: logger,
		masker: masker,
		fields: []zap.Field{},
	}
}

// Context returns the underlying context
func (lc *LogContext) Context() context.Context {
	return lc.ctx
}

// Logger returns the logger with accumulated fields
func (lc *LogContext) Logger() *zap.Logger {
	if len(lc.fields) > 0 {
		return lc.logger.With(lc.fields...)
	}
	return lc.logger
}

// WithFields adds fields to the log context
func (lc *LogContext) WithFields(fields ...zap.Field) *LogContext {
	newContext := &LogContext{
		ctx:    lc.ctx,
		logger: lc.logger,
		masker: lc.masker,
		fields: make([]zap.Field, len(lc.fields)+len(fields)),
	}
	copy(newContext.fields, lc.fields)
	copy(newContext.fields[len(lc.fields):], fields)
	return newContext
}

// WithRequestID adds request ID to the context
func (lc *LogContext) WithRequestID(requestID string) *LogContext {
	ctx := context.WithValue(lc.ctx, contextKeyRequestID, requestID)
	return &LogContext{
		ctx:    ctx,
		logger: lc.logger,
		masker: lc.masker,
		fields: append(lc.fields, zap.String(FieldRequestID, requestID)),
	}
}

// WithCorrelationID adds correlation ID to the context
func (lc *LogContext) WithCorrelationID(correlationID string) *LogContext {
	ctx := context.WithValue(lc.ctx, contextKeyCorrelationID, correlationID)
	return &LogContext{
		ctx:    ctx,
		logger: lc.logger,
		masker: lc.masker,
		fields: append(lc.fields, zap.String(FieldCorrelationID, correlationID)),
	}
}

// WithSessionID adds session ID to the context
func (lc *LogContext) WithSessionID(sessionID string) *LogContext {
	ctx := context.WithValue(lc.ctx, contextKeySessionID, sessionID)
	return &LogContext{
		ctx:    ctx,
		logger: lc.logger,
		masker: lc.masker,
		fields: append(lc.fields, zap.String(FieldSessionID, sessionID)),
	}
}

// WithUserInfo adds user information to the context
func (lc *LogContext) WithUserInfo(userID, auth0UserID, email, name string) *LogContext {
	ctx := lc.ctx
	fields := make([]zap.Field, 0, 4)
	
	if userID != "" {
		ctx = context.WithValue(ctx, contextKeyUserID, userID)
		fields = append(fields, zap.String(FieldUserID, userID))
	}
	if auth0UserID != "" {
		ctx = context.WithValue(ctx, contextKeyAuth0UserID, auth0UserID)
		fields = append(fields, zap.String(FieldAuth0UserID, auth0UserID))
	}
	if email != "" {
		ctx = context.WithValue(ctx, contextKeyUserEmail, email)
		maskedEmail := email
		if lc.masker != nil {
			maskedEmail = lc.masker.MaskValue("email", email)
		}
		fields = append(fields, zap.String(FieldUserEmail, maskedEmail))
	}
	if name != "" {
		ctx = context.WithValue(ctx, contextKeyUserName, name)
		fields = append(fields, zap.String(FieldUserName, name))
	}
	
	return &LogContext{
		ctx:    ctx,
		logger: lc.logger,
		masker: lc.masker,
		fields: append(lc.fields, fields...),
	}
}

// WithComponent adds component name to the context
func (lc *LogContext) WithComponent(component string) *LogContext {
	ctx := context.WithValue(lc.ctx, contextKeyComponent, component)
	return &LogContext{
		ctx:    ctx,
		logger: lc.logger,
		masker: lc.masker,
		fields: append(lc.fields, Component(component)),
	}
}

// WithTraceInfo adds distributed tracing information
func (lc *LogContext) WithTraceInfo(traceID, spanID string) *LogContext {
	ctx := lc.ctx
	fields := make([]zap.Field, 0, 2)
	
	if traceID != "" {
		ctx = context.WithValue(ctx, contextKeyTraceID, traceID)
		fields = append(fields, zap.String(FieldTraceID, traceID))
	}
	if spanID != "" {
		ctx = context.WithValue(ctx, contextKeySpanID, spanID)
		fields = append(fields, zap.String(FieldSpanID, spanID))
	}
	
	return &LogContext{
		ctx:    ctx,
		logger: lc.logger,
		masker: lc.masker,
		fields: append(lc.fields, fields...),
	}
}

// Debug logs a debug message
func (lc *LogContext) Debug(msg string, fields ...zap.Field) {
	lc.Logger().Debug(msg, fields...)
}

// Info logs an info message
func (lc *LogContext) Info(msg string, fields ...zap.Field) {
	lc.Logger().Info(msg, fields...)
}

// Warn logs a warning message
func (lc *LogContext) Warn(msg string, fields ...zap.Field) {
	lc.Logger().Warn(msg, fields...)
}

// Error logs an error message
func (lc *LogContext) Error(msg string, fields ...zap.Field) {
	lc.Logger().Error(msg, fields...)
}

// Fatal logs a fatal message and calls os.Exit(1)
func (lc *LogContext) Fatal(msg string, fields ...zap.Field) {
	lc.Logger().Fatal(msg, fields...)
}

// Panic logs a panic message and panics
func (lc *LogContext) Panic(msg string, fields ...zap.Field) {
	lc.Logger().Panic(msg, fields...)
}

// Context helper functions

// GetRequestID extracts request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(contextKeyRequestID).(string); ok {
		return requestID
	}
	return ""
}

// GetCorrelationID extracts correlation ID from context
func GetCorrelationID(ctx context.Context) string {
	if correlationID, ok := ctx.Value(contextKeyCorrelationID).(string); ok {
		return correlationID
	}
	return ""
}

// GetSessionID extracts session ID from context
func GetSessionID(ctx context.Context) string {
	if sessionID, ok := ctx.Value(contextKeySessionID).(string); ok {
		return sessionID
	}
	return ""
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(contextKeyUserID).(string); ok {
		return userID
	}
	return ""
}

// GetAuth0UserID extracts Auth0 user ID from context
func GetAuth0UserID(ctx context.Context) string {
	if auth0UserID, ok := ctx.Value(contextKeyAuth0UserID).(string); ok {
		return auth0UserID
	}
	return ""
}

// GetUserEmail extracts user email from context
func GetUserEmail(ctx context.Context) string {
	if email, ok := ctx.Value(contextKeyUserEmail).(string); ok {
		return email
	}
	return ""
}

// GetUserName extracts user name from context
func GetUserName(ctx context.Context) string {
	if name, ok := ctx.Value(contextKeyUserName).(string); ok {
		return name
	}
	return ""
}

// GetTraceID extracts trace ID from context
func GetTraceID(ctx context.Context) string {
	if traceID, ok := ctx.Value(contextKeyTraceID).(string); ok {
		return traceID
	}
	return ""
}

// GetSpanID extracts span ID from context
func GetSpanID(ctx context.Context) string {
	if spanID, ok := ctx.Value(contextKeySpanID).(string); ok {
		return spanID
	}
	return ""
}

// GetComponent extracts component name from context
func GetComponent(ctx context.Context) string {
	if component, ok := ctx.Value(contextKeyComponent).(string); ok {
		return component
	}
	return ""
}

// GetLoggerFromContext extracts logger from context
func GetLoggerFromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(contextKeyLogger).(*zap.Logger); ok {
		return logger
	}
	return nil
}

// WithLoggerContext adds logger to context
func WithLoggerContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKeyLogger, logger)
}

// ExtractAllContextFields extracts all logging fields from context
func ExtractAllContextFields(ctx context.Context) []zap.Field {
	var fields []zap.Field

	if requestID := GetRequestID(ctx); requestID != "" {
		fields = append(fields, zap.String(FieldRequestID, requestID))
	}
	if correlationID := GetCorrelationID(ctx); correlationID != "" {
		fields = append(fields, zap.String(FieldCorrelationID, correlationID))
	}
	if sessionID := GetSessionID(ctx); sessionID != "" {
		fields = append(fields, zap.String(FieldSessionID, sessionID))
	}
	if userID := GetUserID(ctx); userID != "" {
		fields = append(fields, zap.String(FieldUserID, userID))
	}
	if auth0UserID := GetAuth0UserID(ctx); auth0UserID != "" {
		fields = append(fields, zap.String(FieldAuth0UserID, auth0UserID))
	}
	if email := GetUserEmail(ctx); email != "" {
		fields = append(fields, zap.String(FieldUserEmail, email))
	}
	if name := GetUserName(ctx); name != "" {
		fields = append(fields, zap.String(FieldUserName, name))
	}
	if traceID := GetTraceID(ctx); traceID != "" {
		fields = append(fields, zap.String(FieldTraceID, traceID))
	}
	if spanID := GetSpanID(ctx); spanID != "" {
		fields = append(fields, zap.String(FieldSpanID, spanID))
	}
	if component := GetComponent(ctx); component != "" {
		fields = append(fields, zap.String(FieldComponent, component))
	}

	return fields
}