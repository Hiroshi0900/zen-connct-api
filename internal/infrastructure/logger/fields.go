package logger

import (
	"time"

	"go.uber.org/zap"
)

// Field constants for consistent logging
const (
	// Service identification
	FieldService        = "service"
	FieldServiceVersion = "service_version"
	FieldEnvironment    = "environment"
	FieldComponent      = "component"

	// Request tracking
	FieldRequestID      = "request_id"
	FieldCorrelationID  = "correlation_id"
	FieldSessionID      = "session_id"
	FieldTraceID        = "trace_id"
	FieldSpanID         = "span_id"

	// User information
	FieldUserID         = "user_id"
	FieldAuth0UserID    = "auth0_user_id"
	FieldUserEmail      = "user_email"
	FieldUserName       = "user_name"
	FieldUserRole       = "user_role"

	// HTTP request/response
	FieldHTTPMethod     = "http_method"
	FieldHTTPPath       = "http_path"
	FieldHTTPQuery      = "http_query"
	FieldHTTPStatus     = "http_status"
	FieldHTTPUserAgent  = "http_user_agent"
	FieldHTTPReferer    = "http_referer"
	FieldHTTPRemoteAddr = "http_remote_addr"
	FieldHTTPHeaders    = "http_headers"

	// Performance metrics
	FieldDuration       = "duration_ms"
	FieldResponseSize   = "response_size_bytes"
	FieldRequestSize    = "request_size_bytes"

	// Database operations
	FieldDBOperation    = "db_operation"
	FieldDBTable        = "db_table"
	FieldDBQuery        = "db_query"
	FieldDBDuration     = "db_duration_ms"
	FieldDBRowsAffected = "db_rows_affected"

	// Authentication & authorization
	FieldAuthAction     = "auth_action"
	FieldAuthResult     = "auth_result"
	FieldAuthProvider   = "auth_provider"
	FieldAuthScope      = "auth_scope"

	// Error handling
	FieldError          = "error"
	FieldErrorCode      = "error_code"
	FieldErrorType      = "error_type"
	FieldStackTrace     = "stack_trace"

	// Business logic
	FieldAction         = "action"
	FieldResource       = "resource"
	FieldResult         = "result"
	FieldReason         = "reason"

	// External services
	FieldExternalService = "external_service"
	FieldExternalURL     = "external_url"
	FieldExternalStatus  = "external_status"
)

// Component names for service identification
const (
	ComponentAuth       = "auth"
	ComponentUser       = "user"
	ComponentSession    = "session"
	ComponentDatabase   = "database"
	ComponentHTTP       = "http"
	ComponentMiddleware = "middleware"
	ComponentHandler    = "handler"
	ComponentRepository = "repository"
	ComponentService    = "service"
	ComponentExternal   = "external"
)

// ServiceInfo creates service identification fields
func ServiceInfo(serviceName, version, environment string) []zap.Field {
	return []zap.Field{
		zap.String(FieldService, serviceName),
		zap.String(FieldServiceVersion, version),
		zap.String(FieldEnvironment, environment),
	}
}

// RequestInfo creates request tracking fields
func RequestInfo(requestID, correlationID, sessionID string) []zap.Field {
	fields := []zap.Field{}
	
	if requestID != "" {
		fields = append(fields, zap.String(FieldRequestID, requestID))
	}
	if correlationID != "" {
		fields = append(fields, zap.String(FieldCorrelationID, correlationID))
	}
	if sessionID != "" {
		fields = append(fields, zap.String(FieldSessionID, sessionID))
	}
	
	return fields
}

// UserInfo creates user information fields
func UserInfo(userID, auth0UserID, email, name string) []zap.Field {
	fields := []zap.Field{}
	
	if userID != "" {
		fields = append(fields, zap.String(FieldUserID, userID))
	}
	if auth0UserID != "" {
		fields = append(fields, zap.String(FieldAuth0UserID, auth0UserID))
	}
	if email != "" {
		fields = append(fields, zap.String(FieldUserEmail, email))
	}
	if name != "" {
		fields = append(fields, zap.String(FieldUserName, name))
	}
	
	return fields
}

// HTTPInfo creates HTTP request/response fields
func HTTPInfo(method, path, query string, status int, userAgent, remoteAddr string) []zap.Field {
	fields := []zap.Field{
		zap.String(FieldHTTPMethod, method),
		zap.String(FieldHTTPPath, path),
		zap.Int(FieldHTTPStatus, status),
	}
	
	if query != "" {
		fields = append(fields, zap.String(FieldHTTPQuery, query))
	}
	if userAgent != "" {
		fields = append(fields, zap.String(FieldHTTPUserAgent, userAgent))
	}
	if remoteAddr != "" {
		fields = append(fields, zap.String(FieldHTTPRemoteAddr, remoteAddr))
	}
	
	return fields
}

// PerformanceInfo creates performance measurement fields
func PerformanceInfo(duration time.Duration, requestSize, responseSize int64) []zap.Field {
	fields := []zap.Field{
		zap.Int64(FieldDuration, duration.Milliseconds()),
	}
	
	if requestSize > 0 {
		fields = append(fields, zap.Int64(FieldRequestSize, requestSize))
	}
	if responseSize > 0 {
		fields = append(fields, zap.Int64(FieldResponseSize, responseSize))
	}
	
	return fields
}

// DatabaseInfo creates database operation fields
func DatabaseInfo(operation, table, query string, duration time.Duration, rowsAffected int64) []zap.Field {
	fields := []zap.Field{
		zap.String(FieldDBOperation, operation),
		zap.Int64(FieldDBDuration, duration.Milliseconds()),
	}
	
	if table != "" {
		fields = append(fields, zap.String(FieldDBTable, table))
	}
	if query != "" {
		fields = append(fields, zap.String(FieldDBQuery, query))
	}
	if rowsAffected >= 0 {
		fields = append(fields, zap.Int64(FieldDBRowsAffected, rowsAffected))
	}
	
	return fields
}

// AuthInfo creates authentication fields
func AuthInfo(action, result, provider, scope string) []zap.Field {
	fields := []zap.Field{}
	
	if action != "" {
		fields = append(fields, zap.String(FieldAuthAction, action))
	}
	if result != "" {
		fields = append(fields, zap.String(FieldAuthResult, result))
	}
	if provider != "" {
		fields = append(fields, zap.String(FieldAuthProvider, provider))
	}
	if scope != "" {
		fields = append(fields, zap.String(FieldAuthScope, scope))
	}
	
	return fields
}

// ErrorInfo creates error information fields
func ErrorInfo(err error, errorCode, errorType string) []zap.Field {
	fields := []zap.Field{}
	
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	if errorCode != "" {
		fields = append(fields, zap.String(FieldErrorCode, errorCode))
	}
	if errorType != "" {
		fields = append(fields, zap.String(FieldErrorType, errorType))
	}
	
	return fields
}

// BusinessInfo creates business logic fields
func BusinessInfo(action, resource, result, reason string) []zap.Field {
	fields := []zap.Field{}
	
	if action != "" {
		fields = append(fields, zap.String(FieldAction, action))
	}
	if resource != "" {
		fields = append(fields, zap.String(FieldResource, resource))
	}
	if result != "" {
		fields = append(fields, zap.String(FieldResult, result))
	}
	if reason != "" {
		fields = append(fields, zap.String(FieldReason, reason))
	}
	
	return fields
}

// ExternalServiceInfo creates external service call fields
func ExternalServiceInfo(service, url string, status int, duration time.Duration) []zap.Field {
	fields := []zap.Field{
		zap.String(FieldExternalService, service),
		zap.Int64(FieldDuration, duration.Milliseconds()),
	}
	
	if url != "" {
		fields = append(fields, zap.String(FieldExternalURL, url))
	}
	if status > 0 {
		fields = append(fields, zap.Int(FieldExternalStatus, status))
	}
	
	return fields
}

// Component creates component identification field
func Component(name string) zap.Field {
	return zap.String(FieldComponent, name)
}

// TraceInfo creates distributed tracing fields
func TraceInfo(traceID, spanID string) []zap.Field {
	fields := []zap.Field{}
	
	if traceID != "" {
		fields = append(fields, zap.String(FieldTraceID, traceID))
	}
	if spanID != "" {
		fields = append(fields, zap.String(FieldSpanID, spanID))
	}
	
	return fields
}

// CombineFields safely combines multiple field slices
func CombineFields(fieldSlices ...[]zap.Field) []zap.Field {
	var totalLength int
	for _, slice := range fieldSlices {
		totalLength += len(slice)
	}
	
	combined := make([]zap.Field, 0, totalLength)
	for _, slice := range fieldSlices {
		combined = append(combined, slice...)
	}
	
	return combined
}

// SafeString creates a string field with masking if needed
func SafeString(masker *Masker, key, value string) zap.Field {
	if masker != nil {
		value = masker.MaskValue(key, value)
	}
	return zap.String(key, value)
}

// SafeStringMap creates a string map field with masking if needed
func SafeStringMap(masker *Masker, key string, values map[string]string) zap.Field {
	if masker != nil {
		safeValues := make(map[string]string)
		for k, v := range values {
			safeValues[k] = masker.MaskValue(k, v)
		}
		return zap.Any(key, safeValues)
	}
	return zap.Any(key, values)
}