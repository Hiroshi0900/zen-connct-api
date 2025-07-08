package logger

import (
	"regexp"
	"strings"
)

// Masker handles sensitive information masking
type Masker struct {
	config *Config
}

// NewMasker creates a new masker instance
func NewMasker(config *Config) *Masker {
	return &Masker{
		config: config,
	}
}

// MaskSensitiveData masks sensitive information in text
func (m *Masker) MaskSensitiveData(text string) string {
	if text == "" {
		return text
	}

	result := text

	// Mask passwords
	if m.config.MaskPasswords {
		result = m.maskPasswords(result)
	}

	// Mask tokens
	if m.config.MaskTokens {
		result = m.maskTokens(result)
	}

	// Mask emails
	if m.config.MaskEmails != MaskingNone {
		result = m.maskEmails(result)
	}

	// Mask credit card numbers
	result = m.maskCreditCards(result)

	return result
}

// MaskValue masks a specific value based on its key
func (m *Masker) MaskValue(key, value string) string {
	if value == "" {
		return value
	}

	keyLower := strings.ToLower(key)

	// Password fields
	if m.config.MaskPasswords && m.isPasswordField(keyLower) {
		return m.maskString(value, 0)
	}

	// Token fields
	if m.config.MaskTokens && m.isTokenField(keyLower) {
		return m.maskToken(value)
	}

	// Email fields
	if m.config.MaskEmails != MaskingNone && m.isEmailField(keyLower) {
		return m.maskEmail(value)
	}

	// Credit card fields
	if m.isCreditCardField(keyLower) {
		return m.maskString(value, 0)
	}

	// Apply general masking to the value
	return m.MaskSensitiveData(value)
}

// isPasswordField checks if a field name indicates it contains a password
func (m *Masker) isPasswordField(fieldName string) bool {
	passwordPatterns := []string{
		"password", "passwd", "pwd", "secret", "auth", "token", "key",
		"credentials", "credential", "pass", "passphrase",
	}

	for _, pattern := range passwordPatterns {
		if strings.Contains(fieldName, pattern) {
			return true
		}
	}

	return false
}

// isTokenField checks if a field name indicates it contains a token
func (m *Masker) isTokenField(fieldName string) bool {
	tokenPatterns := []string{
		"token", "jwt", "bearer", "authorization", "access_token",
		"refresh_token", "id_token", "api_key", "apikey",
	}

	for _, pattern := range tokenPatterns {
		if strings.Contains(fieldName, pattern) {
			return true
		}
	}

	return false
}

// isEmailField checks if a field name indicates it contains an email
func (m *Masker) isEmailField(fieldName string) bool {
	emailPatterns := []string{
		"email", "mail", "e-mail", "emailaddress", "email_address",
	}

	for _, pattern := range emailPatterns {
		if strings.Contains(fieldName, pattern) {
			return true
		}
	}

	return false
}

// isCreditCardField checks if a field name indicates it contains credit card info
func (m *Masker) isCreditCardField(fieldName string) bool {
	ccPatterns := []string{
		"card", "credit", "debit", "cc", "cvv", "cvc", "card_number",
		"cardnumber", "creditcard", "credit_card",
	}

	for _, pattern := range ccPatterns {
		if strings.Contains(fieldName, pattern) {
			return true
		}
	}

	return false
}

// maskPasswords masks password-like patterns in text
func (m *Masker) maskPasswords(text string) string {
	// Pattern for password in JSON or form data
	passwordRegex := regexp.MustCompile(`(?i)(password|passwd|pwd|secret)["']?\s*[:=]\s*["']?([^"',\s}]+)`)
	return passwordRegex.ReplaceAllStringFunc(text, func(match string) string {
		parts := passwordRegex.FindStringSubmatch(match)
		if len(parts) >= 3 {
			return strings.Replace(match, parts[2], "***MASKED***", 1)
		}
		return match
	})
}

// maskTokens masks token-like patterns in text
func (m *Masker) maskTokens(text string) string {
	// JWT token pattern
	jwtRegex := regexp.MustCompile(`eyJ[A-Za-z0-9+/=]*\.eyJ[A-Za-z0-9+/=]*\.[A-Za-z0-9+/=]*`)
	text = jwtRegex.ReplaceAllStringFunc(text, func(token string) string {
		return m.maskToken(token)
	})

	// Bearer token pattern
	bearerRegex := regexp.MustCompile(`(?i)bearer\s+([A-Za-z0-9+/=._-]+)`)
	text = bearerRegex.ReplaceAllStringFunc(text, func(match string) string {
		parts := bearerRegex.FindStringSubmatch(match)
		if len(parts) >= 2 {
			return strings.Replace(match, parts[1], m.maskToken(parts[1]), 1)
		}
		return match
	})

	// API key pattern
	apiKeyRegex := regexp.MustCompile(`(?i)(api[_-]?key|apikey)["']?\s*[:=]\s*["']?([A-Za-z0-9+/=._-]{16,})`)
	text = apiKeyRegex.ReplaceAllStringFunc(text, func(match string) string {
		parts := apiKeyRegex.FindStringSubmatch(match)
		if len(parts) >= 3 {
			return strings.Replace(match, parts[2], m.maskToken(parts[2]), 1)
		}
		return match
	})

	return text
}

// maskEmails masks email addresses in text
func (m *Masker) maskEmails(text string) string {
	emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	return emailRegex.ReplaceAllStringFunc(text, func(email string) string {
		return m.maskEmail(email)
	})
}

// maskCreditCards masks credit card numbers
func (m *Masker) maskCreditCards(text string) string {
	// Credit card number pattern (13-19 digits, may have spaces or dashes)
	ccRegex := regexp.MustCompile(`\b(?:\d{4}[\s-]?){3}\d{1,4}\b`)
	return ccRegex.ReplaceAllString(text, "****-****-****-****")
}

// maskToken masks a token showing only first 4 characters
func (m *Masker) maskToken(token string) string {
	if len(token) <= 8 {
		return "***MASKED***"
	}
	return token[:4] + "..." + strings.Repeat("*", len(token)-8) + "..." + token[len(token)-4:]
}

// maskEmail masks an email address based on masking level
func (m *Masker) maskEmail(email string) string {
	switch m.config.MaskEmails {
	case MaskingFull:
		return "***@***.***"
	case MaskingPartial:
		parts := strings.Split(email, "@")
		if len(parts) != 2 {
			return email
		}
		username := parts[0]
		domain := parts[1]

		// Mask username (show first and last character if length > 2)
		var maskedUsername string
		if len(username) <= 2 {
			maskedUsername = strings.Repeat("*", len(username))
		} else {
			maskedUsername = string(username[0]) + strings.Repeat("*", len(username)-2) + string(username[len(username)-1])
		}

		// Mask domain (show first character and TLD)
		domainParts := strings.Split(domain, ".")
		if len(domainParts) >= 2 {
			mainDomain := domainParts[0]
			tld := domainParts[len(domainParts)-1]
			var maskedDomain string
			if len(mainDomain) <= 2 {
				maskedDomain = strings.Repeat("*", len(mainDomain))
			} else {
				maskedDomain = string(mainDomain[0]) + strings.Repeat("*", len(mainDomain)-1)
			}
			return maskedUsername + "@" + maskedDomain + "." + tld
		}
		return maskedUsername + "@" + domain
	case MaskingNone:
		return email
	default:
		return email
	}
}

// maskString masks a string showing only the specified number of characters from the start
func (m *Masker) maskString(value string, showChars int) string {
	if len(value) <= showChars {
		return strings.Repeat("*", len(value))
	}
	if showChars == 0 {
		return "***MASKED***"
	}
	return value[:showChars] + strings.Repeat("*", len(value)-showChars)
}

// SanitizeLogMessage sanitizes a log message by removing or masking sensitive data
func (m *Masker) SanitizeLogMessage(message string) string {
	return m.MaskSensitiveData(message)
}

// SanitizeStructuredData sanitizes structured data (map) by masking sensitive values
func (m *Masker) SanitizeStructuredData(data map[string]interface{}) map[string]interface{} {
	sanitized := make(map[string]interface{})

	for key, value := range data {
		switch v := value.(type) {
		case string:
			sanitized[key] = m.MaskValue(key, v)
		case map[string]interface{}:
			sanitized[key] = m.SanitizeStructuredData(v)
		case []interface{}:
			sanitized[key] = m.sanitizeSlice(key, v)
		default:
			sanitized[key] = value
		}
	}

	return sanitized
}

// sanitizeSlice sanitizes slice data
func (m *Masker) sanitizeSlice(key string, slice []interface{}) []interface{} {
	sanitized := make([]interface{}, len(slice))

	for i, item := range slice {
		switch v := item.(type) {
		case string:
			sanitized[i] = m.MaskValue(key, v)
		case map[string]interface{}:
			sanitized[i] = m.SanitizeStructuredData(v)
		default:
			sanitized[i] = item
		}
	}

	return sanitized
}