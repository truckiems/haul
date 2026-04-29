package env

import (
	"testing"
)

func TestValidate_AllValid(t *testing.T) {
	env := map[string]string{
		"APP_ENV":       "production",
		"DATABASE_URL":  "postgres://localhost/db",
		"_PRIVATE":      "secret",
		"MixedCase123":  "value",
	}
	result := Validate(env)
	if result.HasErrors() {
		t.Errorf("expected no errors, got: %s", result.Summary())
	}
}

func TestValidate_StartsWithDigit(t *testing.T) {
	env := map[string]string{
		"1BAD_KEY": "value",
	}
	result := Validate(env)
	if !result.HasErrors() {
		t.Fatal("expected validation error for key starting with digit")
	}
	if result.Errors[0].Key != "1BAD_KEY" {
		t.Errorf("unexpected key in error: %s", result.Errors[0].Key)
	}
}

func TestValidate_InvalidCharacter(t *testing.T) {
	env := map[string]string{
		"BAD-KEY": "value",
	}
	result := Validate(env)
	if !result.HasErrors() {
		t.Fatal("expected validation error for key with hyphen")
	}
}

func TestValidate_EmptyKey(t *testing.T) {
	env := map[string]string{
		"": "value",
	}
	result := Validate(env)
	if !result.HasErrors() {
		t.Fatal("expected validation error for empty key")
	}
	if result.Errors[0].Message != "key must not be empty" {
		t.Errorf("unexpected message: %s", result.Errors[0].Message)
	}
}

func TestValidate_MultipleErrors(t *testing.T) {
	env := map[string]string{
		"GOOD_KEY": "ok",
		"BAD KEY":  "has space",
		"2NDKEY":   "starts with digit",
	}
	result := Validate(env)
	if len(result.Errors) < 2 {
		t.Errorf("expected at least 2 errors, got %d", len(result.Errors))
	}
}

func TestValidationResult_Summary_NoErrors(t *testing.T) {
	result := ValidationResult{}
	if result.Summary() != "all keys valid" {
		t.Errorf("unexpected summary: %s", result.Summary())
	}
}

func TestValidationResult_Summary_WithErrors(t *testing.T) {
	result := ValidationResult{
		Errors: []ValidationError{
			{Key: "BAD KEY", Message: "invalid character ' ' in key"},
		},
	}
	summary := result.Summary()
	if summary == "all keys valid" {
		t.Error("expected non-empty summary for errors")
	}
}
