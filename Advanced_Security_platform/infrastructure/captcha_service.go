package infrastructure

import domain "security/domain"

// CaptchaValidator is a placeholder; integrate real provider (e.g., reCAPTCHA) later.
type CaptchaValidator struct{}

func NewCaptchaValidator() domain.ICaptchaValidator {
	return &CaptchaValidator{}
}

func (c *CaptchaValidator) Validate(token string) bool {
	// TODO: call external CAPTCHA service. For now, accept non-empty token.
	return token != ""
}

