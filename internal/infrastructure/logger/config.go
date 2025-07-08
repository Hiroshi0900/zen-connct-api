package logger

import (
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogLevel represents logging levels
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	FatalLevel LogLevel = "fatal"
	PanicLevel LogLevel = "panic"
)

// LogFormat represents logging output formats
type LogFormat string

const (
	JSONFormat    LogFormat = "json"
	ConsoleFormat LogFormat = "console"
)

// LogOutput represents logging output destinations
type LogOutput string

const (
	StdoutOutput LogOutput = "stdout"
	StderrOutput LogOutput = "stderr"
	FileOutput   LogOutput = "file"
)

// MaskingLevel represents levels of information masking
type MaskingLevel string

const (
	MaskingNone    MaskingLevel = "none"
	MaskingPartial MaskingLevel = "partial"
	MaskingFull    MaskingLevel = "full"
)

// Config holds all logger configuration
type Config struct {
	Level          LogLevel     `json:"level"`
	Format         LogFormat    `json:"format"`
	Output         LogOutput    `json:"output"`
	FilePath       string       `json:"file_path"`
	MaxSize        int          `json:"max_size"`        // MB
	MaxAge         int          `json:"max_age"`         // days
	MaxBackups     int          `json:"max_backups"`
	SamplingRate   float64      `json:"sampling_rate"`   // 0.0-1.0
	MaskPasswords  bool         `json:"mask_passwords"`
	MaskTokens     bool         `json:"mask_tokens"`
	MaskEmails     MaskingLevel `json:"mask_emails"`
	ServiceName    string       `json:"service_name"`
	ServiceVersion string       `json:"service_version"`
	Environment    string       `json:"environment"`
}

// NewConfig creates a new logger configuration from environment variables
func NewConfig() *Config {
	return &Config{
		Level:          parseLogLevel(getEnvOrDefault("LOG_LEVEL", "info")),
		Format:         parseLogFormat(getEnvOrDefault("LOG_FORMAT", "console")),
		Output:         parseLogOutput(getEnvOrDefault("LOG_OUTPUT", "stdout")),
		FilePath:       getEnvOrDefault("LOG_FILE_PATH", "/var/log/zenconnect.log"),
		MaxSize:        getEnvIntOrDefault("LOG_MAX_SIZE", 100),
		MaxAge:         getEnvIntOrDefault("LOG_MAX_AGE", 30),
		MaxBackups:     getEnvIntOrDefault("LOG_MAX_BACKUPS", 10),
		SamplingRate:   getEnvFloatOrDefault("LOG_SAMPLING_RATE", 1.0),
		MaskPasswords:  getEnvBoolOrDefault("LOG_MASK_PASSWORDS", true),
		MaskTokens:     getEnvBoolOrDefault("LOG_MASK_TOKENS", true),
		MaskEmails:     parseMaskingLevel(getEnvOrDefault("LOG_MASK_EMAILS", "partial")),
		ServiceName:    getEnvOrDefault("SERVICE_NAME", "zen-connect"),
		ServiceVersion: getEnvOrDefault("SERVICE_VERSION", "1.0.0"),
		Environment:    getEnvOrDefault("ENVIRONMENT", "development"),
	}
}

// ToZapConfig converts Config to zap.Config
func (c *Config) ToZapConfig() zap.Config {
	var config zap.Config

	// Base configuration based on environment
	if c.Environment == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Set log level
	config.Level = zap.NewAtomicLevelAt(c.toZapLevel())

	// Set encoding format
	if c.Format == JSONFormat {
		config.Encoding = "json"
		config.EncoderConfig = zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	} else {
		config.Encoding = "console"
		config.EncoderConfig = zapcore.EncoderConfig{
			TimeKey:        "T",
			LevelKey:       "L",
			NameKey:        "N",
			CallerKey:      "C",
			MessageKey:     "M",
			StacktraceKey:  "S",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	}

	// Set output paths
	switch c.Output {
	case StdoutOutput:
		config.OutputPaths = []string{"stdout"}
	case StderrOutput:
		config.OutputPaths = []string{"stderr"}
	case FileOutput:
		config.OutputPaths = []string{c.FilePath}
	default:
		config.OutputPaths = []string{"stdout"}
	}

	// Error output always goes to stderr
	config.ErrorOutputPaths = []string{"stderr"}

	// Sampling configuration for high-volume environments
	if c.SamplingRate < 1.0 {
		config.Sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: int(100 * c.SamplingRate),
		}
	}

	return config
}

// toZapLevel converts LogLevel to zapcore.Level
func (c *Config) toZapLevel() zapcore.Level {
	switch c.Level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	case PanicLevel:
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

// IsProductionMode returns true if running in production environment
func (c *Config) IsProductionMode() bool {
	return c.Environment == "production"
}

// Helper functions for environment variable parsing

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvFloatOrDefault(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func getEnvBoolOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

func parseLogLevel(level string) LogLevel {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	case "panic":
		return PanicLevel
	default:
		return InfoLevel
	}
}

func parseLogFormat(format string) LogFormat {
	switch strings.ToLower(format) {
	case "json":
		return JSONFormat
	case "console":
		return ConsoleFormat
	default:
		return ConsoleFormat
	}
}

func parseLogOutput(output string) LogOutput {
	switch strings.ToLower(output) {
	case "stdout":
		return StdoutOutput
	case "stderr":
		return StderrOutput
	case "file":
		return FileOutput
	default:
		return StdoutOutput
	}
}

func parseMaskingLevel(level string) MaskingLevel {
	switch strings.ToLower(level) {
	case "none":
		return MaskingNone
	case "partial":
		return MaskingPartial
	case "full":
		return MaskingFull
	default:
		return MaskingPartial
	}
}