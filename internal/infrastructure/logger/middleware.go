package logger

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// EchoLoggerConfig holds configuration for Echo logger middleware
type EchoLoggerConfig struct {
	Logger               *Logger
	RequestIDHeader      string
	CorrelationIDHeader  string
	SkipPaths           []string
	LogRequestBody      bool
	LogResponseBody     bool
	MaxBodySize         int64
	HideHealthChecks    bool
}

// DefaultEchoLoggerConfig returns default configuration
func DefaultEchoLoggerConfig() EchoLoggerConfig {
	return EchoLoggerConfig{
		Logger:              GetGlobalLogger(),
		RequestIDHeader:     "X-Request-ID",
		CorrelationIDHeader: "X-Correlation-ID",
		SkipPaths:          []string{},
		LogRequestBody:     false,
		LogResponseBody:    false,
		MaxBodySize:        1024 * 1024, // 1MB
		HideHealthChecks:   true,
	}
}

// RequestLoggerMiddleware returns Echo middleware for request logging
func RequestLoggerMiddleware(config ...EchoLoggerConfig) echo.MiddlewareFunc {
	cfg := DefaultEchoLoggerConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	return RequestLoggerWithConfig(cfg)
}

// RequestLoggerWithConfig returns Echo middleware with custom config
func RequestLoggerWithConfig(config EchoLoggerConfig) echo.MiddlewareFunc {
	if config.Logger == nil {
		config.Logger = GetGlobalLogger()
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip logging for certain paths
			if shouldSkipPath(c.Request().URL.Path, config.SkipPaths, config.HideHealthChecks) {
				return next(c)
			}

			start := time.Now()
			req := c.Request()
			res := c.Response()

			// Generate or extract request ID
			requestID := req.Header.Get(config.RequestIDHeader)
			if requestID == "" {
				requestID = generateRequestID()
				req.Header.Set(config.RequestIDHeader, requestID)
				res.Header().Set(config.RequestIDHeader, requestID)
			}

			// Extract correlation ID
			correlationID := req.Header.Get(config.CorrelationIDHeader)
			if correlationID == "" {
				correlationID = requestID // Use request ID as correlation ID if not provided
			}

			// Create log context
			logCtx := config.Logger.WithContext(req.Context()).
				WithRequestID(requestID).
				WithCorrelationID(correlationID).
				WithComponent(ComponentHTTP)

			// Add log context to request context
			ctx := WithLoggerContext(req.Context(), logCtx.Logger())
			c.SetRequest(req.WithContext(ctx))

			// Capture request body if configured
			var requestBody []byte
			if config.LogRequestBody && req.Body != nil {
				requestBody = captureRequestBody(req, config.MaxBodySize)
			}

			// Capture response body if configured
			var responseBody []byte
			if config.LogResponseBody {
				responseBody = captureResponseBody(res, config.MaxBodySize)
			}

			// Log request start
			logRequestStart(logCtx, req, requestBody)

			// Process request
			err := next(c)

			// Calculate duration
			duration := time.Since(start)

			// Log request completion
			logRequestComplete(logCtx, req, res, duration, requestBody, responseBody, err)

			return err
		}
	}
}

// SessionLoggerMiddleware adds session information to log context
func SessionLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			
			// Extract session information from context (set by session middleware)
			userID := GetUserID(ctx)
			auth0UserID := GetAuth0UserID(ctx)
			userEmail := GetUserEmail(ctx)
			userName := GetUserName(ctx)
			sessionID := GetSessionID(ctx)

			// If we have user information, add it to the log context
			if userID != "" || auth0UserID != "" {
				if logger := GetLoggerFromContext(ctx); logger != nil {
					logCtx := NewLogContext(ctx, logger, GetGlobalLogger().masker).
						WithUserInfo(userID, auth0UserID, userEmail, userName)
					
					if sessionID != "" {
						logCtx = logCtx.WithSessionID(sessionID)
					}

					// Update request context
					ctx = WithLoggerContext(ctx, logCtx.Logger())
					c.SetRequest(c.Request().WithContext(ctx))
				}
			}

			return next(c)
		}
	}
}

// ErrorLoggerMiddleware logs errors in a structured way
func ErrorLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			
			if err != nil {
				ctx := c.Request().Context()
				
				// Get logger from context or use global logger
				logger := GetLoggerFromContext(ctx)
				if logger == nil {
					logger = GetGlobalLogger().zap
				}

				// Extract all context fields
				fields := ExtractAllContextFields(ctx)
				
				// Add HTTP information
				fields = append(fields,
					zap.String(FieldHTTPMethod, c.Request().Method),
					zap.String(FieldHTTPPath, c.Request().URL.Path),
					zap.Int(FieldHTTPStatus, c.Response().Status),
				)

				// Add error information
				if httpError, ok := err.(*echo.HTTPError); ok {
					fields = append(fields,
						zap.Int(FieldErrorCode, httpError.Code),
						zap.String(FieldErrorType, "http_error"),
						zap.Any(FieldError, httpError.Message),
					)
				} else {
					fields = append(fields,
						zap.String(FieldErrorType, "internal_error"),
						zap.Error(err),
					)
				}

				logger.Error("Request error", fields...)
			}

			return err
		}
	}
}

// RecoveryLoggerMiddleware logs panic recoveries
func RecoveryLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					ctx := c.Request().Context()
					
					// Get logger from context or use global logger
					logger := GetLoggerFromContext(ctx)
					if logger == nil {
						logger = GetGlobalLogger().zap
					}

					// Extract all context fields
					fields := ExtractAllContextFields(ctx)
					
					// Add HTTP information
					fields = append(fields,
						zap.String(FieldHTTPMethod, c.Request().Method),
						zap.String(FieldHTTPPath, c.Request().URL.Path),
						zap.String(FieldErrorType, "panic"),
						zap.Any("panic_value", r),
						Component("recovery"),
					)

					logger.Error("Request panic recovered", fields...)

					// Return internal server error
					err = echo.NewHTTPError(500, "Internal Server Error")
				}
			}()

			return next(c)
		}
	}
}

// Helper functions

func shouldSkipPath(path string, skipPaths []string, hideHealthChecks bool) bool {
	// Skip health check endpoints if configured
	if hideHealthChecks && (path == "/health" || strings.HasPrefix(path, "/health/")) {
		return true
	}

	// Skip explicitly configured paths
	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	return false
}

func generateRequestID() string {
	return uuid.New().String()
}

func captureRequestBody(req *http.Request, maxSize int64) []byte {
	if req.Body == nil {
		return nil
	}

	// Read body
	bodyBytes, err := io.ReadAll(io.LimitReader(req.Body, maxSize))
	if err != nil {
		return nil
	}

	// Restore body for further processing
	req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return bodyBytes
}

func captureResponseBody(res *echo.Response, maxSize int64) []byte {
	// This is more complex in Echo as we need to wrap the writer
	// For now, we'll skip response body capture to avoid complexity
	// In production, you might want to use a custom response writer
	return nil
}

func logRequestStart(logCtx *LogContext, req *http.Request, requestBody []byte) {
	fields := []zap.Field{
		zap.String(FieldHTTPMethod, req.Method),
		zap.String(FieldHTTPPath, req.URL.Path),
		zap.String(FieldHTTPUserAgent, req.UserAgent()),
		zap.String(FieldHTTPRemoteAddr, req.RemoteAddr),
	}

	if req.URL.RawQuery != "" {
		fields = append(fields, zap.String(FieldHTTPQuery, req.URL.RawQuery))
	}

	if len(requestBody) > 0 {
		fields = append(fields, zap.String("request_body", string(requestBody)))
	}

	if contentLength := req.Header.Get("Content-Length"); contentLength != "" {
		if size, err := strconv.ParseInt(contentLength, 10, 64); err == nil {
			fields = append(fields, zap.Int64(FieldRequestSize, size))
		}
	}

	// Lower log level for frequently accessed endpoints
	isFrequentEndpoint := req.URL.Path == "/auth/me" || req.URL.Path == "/health"
	
	if isFrequentEndpoint {
		logCtx.Debug("HTTP request started", fields...)
	} else {
		logCtx.Info("HTTP request started", fields...)
	}
}

func logRequestComplete(logCtx *LogContext, req *http.Request, res *echo.Response, duration time.Duration, requestBody, responseBody []byte, err error) {
	fields := []zap.Field{
		zap.String(FieldHTTPMethod, req.Method),
		zap.String(FieldHTTPPath, req.URL.Path),
		zap.Int(FieldHTTPStatus, res.Status),
		zap.Int64(FieldDuration, duration.Milliseconds()),
		zap.Int64(FieldResponseSize, res.Size),
	}

	if req.URL.RawQuery != "" {
		fields = append(fields, zap.String(FieldHTTPQuery, req.URL.RawQuery))
	}

	if len(responseBody) > 0 {
		fields = append(fields, zap.String("response_body", string(responseBody)))
	}

	// Log level based on status code and error
	message := "HTTP request completed"
	
	// Lower log level for frequently accessed endpoints
	isFrequentEndpoint := req.URL.Path == "/auth/me" || req.URL.Path == "/health"
	
	if err != nil {
		fields = append(fields, zap.Error(err))
		logCtx.Error(message, fields...)
	} else if res.Status >= 500 {
		logCtx.Error(message, fields...)
	} else if res.Status >= 400 {
		logCtx.Warn(message, fields...)
	} else if isFrequentEndpoint {
		logCtx.Debug(message, fields...)
	} else {
		logCtx.Info(message, fields...)
	}
}