package logger

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents the enterprise-grade logger
type Logger struct {
	zap    *zap.Logger
	config *Config
	masker *Masker
	mutex  sync.RWMutex
}

// Global logger instance
var (
	globalLogger *Logger
	once         sync.Once
)

// Initialize initializes the global logger with configuration
func Initialize(config *Config) error {
	var err error
	once.Do(func() {
		globalLogger, err = New(config)
	})
	return err
}

// New creates a new logger instance
func New(config *Config) (*Logger, error) {
	if config == nil {
		config = NewConfig()
	}

	// Build zap logger
	zapConfig := config.ToZapConfig()
	zapLogger, err := zapConfig.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	// Add service information as default fields
	zapLogger = zapLogger.With(ServiceInfo(
		config.ServiceName,
		config.ServiceVersion,
		config.Environment,
	)...)

	// Create masker
	masker := NewMasker(config)

	return &Logger{
		zap:    zapLogger,
		config: config,
		masker: masker,
	}, nil
}

// GetGlobalLogger returns the global logger instance
func GetGlobalLogger() *Logger {
	if globalLogger == nil {
		// Initialize with default config if not already initialized
		_ = Initialize(NewConfig())
	}
	return globalLogger
}

// Close flushes any pending log entries
func (l *Logger) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	
	if l.zap != nil {
		return l.zap.Sync()
	}
	return nil
}

// WithContext creates a log context for request-scoped logging
func (l *Logger) WithContext(ctx context.Context) *LogContext {
	return NewLogContext(ctx, l.zap, l.masker)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	l.zap.Debug(l.masker.SanitizeLogMessage(msg), l.sanitizeFields(fields...)...)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	l.zap.Info(l.masker.SanitizeLogMessage(msg), l.sanitizeFields(fields...)...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	l.zap.Warn(l.masker.SanitizeLogMessage(msg), l.sanitizeFields(fields...)...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	l.zap.Error(l.masker.SanitizeLogMessage(msg), l.sanitizeFields(fields...)...)
}

// Fatal logs a fatal message and calls os.Exit(1)
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	l.zap.Fatal(l.masker.SanitizeLogMessage(msg), l.sanitizeFields(fields...)...)
}

// Panic logs a panic message and panics
func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	l.zap.Panic(l.masker.SanitizeLogMessage(msg), l.sanitizeFields(fields...)...)
}

// WithFields adds fields to the logger
func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	
	return &Logger{
		zap:    l.zap.With(l.sanitizeFields(fields...)...),
		config: l.config,
		masker: l.masker,
	}
}

// WithComponent adds component information to the logger
func (l *Logger) WithComponent(component string) *Logger {
	return l.WithFields(Component(component))
}

// sanitizeFields sanitizes field values using the masker
func (l *Logger) sanitizeFields(fields ...zap.Field) []zap.Field {
	if l.masker == nil {
		return fields
	}

	sanitized := make([]zap.Field, len(fields))
	for i, field := range fields {
		sanitized[i] = l.sanitizeField(field)
	}
	return sanitized
}

// sanitizeField sanitizes a single field
func (l *Logger) sanitizeField(field zap.Field) zap.Field {
	switch field.Type {
	case zapcore.StringType:
		return zap.String(field.Key, l.masker.MaskValue(field.Key, field.String))
	case zapcore.ErrorType:
		if field.Interface != nil {
			if err, ok := field.Interface.(error); ok {
				return zap.String(field.Key, l.masker.SanitizeLogMessage(err.Error()))
			}
		}
		return field
	default:
		return field
	}
}

// Convenience methods for structured logging

// LogHTTPRequest logs an HTTP request with all relevant information
func (l *Logger) LogHTTPRequest(ctx context.Context, method, path, query string, status int, duration time.Duration, requestSize, responseSize int64, userAgent, remoteAddr string) {
	fields := CombineFields(
		ExtractAllContextFields(ctx),
		HTTPInfo(method, path, query, status, userAgent, remoteAddr),
		PerformanceInfo(duration, requestSize, responseSize),
		[]zap.Field{Component(ComponentHTTP)},
	)
	
	l.Info("HTTP request completed", fields...)
}

// LogDatabaseOperation logs a database operation
func (l *Logger) LogDatabaseOperation(ctx context.Context, operation, table, query string, duration time.Duration, rowsAffected int64, err error) {
	fields := CombineFields(
		ExtractAllContextFields(ctx),
		DatabaseInfo(operation, table, query, duration, rowsAffected),
		[]zap.Field{Component(ComponentDatabase)},
	)
	
	if err != nil {
		fields = append(fields, ErrorInfo(err, "", "database_error")...)
		l.Error("Database operation failed", fields...)
	} else {
		l.Info("Database operation completed", fields...)
	}
}

// LogAuthenticationAttempt logs an authentication attempt
func (l *Logger) LogAuthenticationAttempt(ctx context.Context, action, result, provider, email string, err error) {
	fields := CombineFields(
		ExtractAllContextFields(ctx),
		AuthInfo(action, result, provider, ""),
		[]zap.Field{
			Component(ComponentAuth),
			SafeString(l.masker, "email", email),
		},
	)
	
	if err != nil {
		fields = append(fields, ErrorInfo(err, "", "auth_error")...)
		l.Warn("Authentication attempt failed", fields...)
	} else {
		l.Info("Authentication attempt successful", fields...)
	}
}

// LogExternalServiceCall logs an external service call
func (l *Logger) LogExternalServiceCall(ctx context.Context, service, url string, status int, duration time.Duration, err error) {
	fields := CombineFields(
		ExtractAllContextFields(ctx),
		ExternalServiceInfo(service, url, status, duration),
		[]zap.Field{Component(ComponentExternal)},
	)
	
	if err != nil {
		fields = append(fields, ErrorInfo(err, "", "external_service_error")...)
		l.Error("External service call failed", fields...)
	} else {
		l.Info("External service call completed", fields...)
	}
}

// LogBusinessAction logs a business logic action
func (l *Logger) LogBusinessAction(ctx context.Context, action, resource, result, reason string, duration time.Duration) {
	fields := CombineFields(
		ExtractAllContextFields(ctx),
		BusinessInfo(action, resource, result, reason),
		PerformanceInfo(duration, 0, 0),
		[]zap.Field{Component(ComponentService)},
	)
	
	if result == "success" {
		l.Info("Business action completed", fields...)
	} else {
		l.Warn("Business action failed", fields...)
	}
}

// LogSecurityEvent logs a security-related event
func (l *Logger) LogSecurityEvent(ctx context.Context, event, severity, description string, additionalFields ...zap.Field) {
	fields := CombineFields(
		ExtractAllContextFields(ctx),
		[]zap.Field{
			Component("security"),
			zap.String("security_event", event),
			zap.String("severity", severity),
			zap.String("description", l.masker.SanitizeLogMessage(description)),
		},
		additionalFields,
	)
	
	switch severity {
	case "critical", "high":
		l.Error("Security event detected", fields...)
	case "medium":
		l.Warn("Security event detected", fields...)
	default:
		l.Info("Security event detected", fields...)
	}
}

// Global convenience functions

// Debug logs a debug message using the global logger
func Debug(msg string, fields ...zap.Field) {
	GetGlobalLogger().Debug(msg, fields...)
}

// Info logs an info message using the global logger
func Info(msg string, fields ...zap.Field) {
	GetGlobalLogger().Info(msg, fields...)
}

// Warn logs a warning message using the global logger
func Warn(msg string, fields ...zap.Field) {
	GetGlobalLogger().Warn(msg, fields...)
}

// Error logs an error message using the global logger
func Error(msg string, fields ...zap.Field) {
	GetGlobalLogger().Error(msg, fields...)
}

// Fatal logs a fatal message using the global logger
func Fatal(msg string, fields ...zap.Field) {
	GetGlobalLogger().Fatal(msg, fields...)
}

// Panic logs a panic message using the global logger
func Panic(msg string, fields ...zap.Field) {
	GetGlobalLogger().Panic(msg, fields...)
}

// WithContext creates a log context using the global logger
func WithContext(ctx context.Context) *LogContext {
	return GetGlobalLogger().WithContext(ctx)
}

// WithComponent adds component information using the global logger
func WithComponent(component string) *Logger {
	return GetGlobalLogger().WithComponent(component)
}

// Close flushes any pending log entries from the global logger
func Close() error {
	if globalLogger != nil {
		return globalLogger.Close()
	}
	return nil
}