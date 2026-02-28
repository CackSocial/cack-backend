package auth

import (
	"testing"
)

func TestGenerateToken(t *testing.T) {
	token, err := GenerateToken("user-123", "mysecret", 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestValidateToken_Valid(t *testing.T) {
	secret := "mysecret"
	token, err := GenerateToken("user-123", secret, 1)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	parsed, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !parsed.Valid {
		t.Fatal("expected token to be valid")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	secret := "mysecret"
	// Using -1 hours expiry to create an already-expired token.
	token, err := GenerateToken("user-123", secret, -1)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = ValidateToken(token, secret)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	token, err := GenerateToken("user-123", "secret1", 1)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	_, err = ValidateToken(token, "wrongsecret")
	if err == nil {
		t.Fatal("expected error for wrong secret, got nil")
	}
}

func TestExtractUserID(t *testing.T) {
	secret := "mysecret"
	userID := "user-456"

	tokenStr, err := GenerateToken(userID, secret, 1)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	parsed, err := ValidateToken(tokenStr, secret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	extracted, err := ExtractUserID(parsed)
	if err != nil {
		t.Fatalf("ExtractUserID failed: %v", err)
	}
	if extracted != userID {
		t.Fatalf("expected %q, got %q", userID, extracted)
	}
}
