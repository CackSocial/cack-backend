package hash

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "mysecretpassword"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if hashed == "" {
		t.Fatal("expected non-empty hash")
	}
	if hashed == password {
		t.Fatal("hash should differ from the original password")
	}
}

func TestCheckPassword_Correct(t *testing.T) {
	password := "mysecretpassword"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !CheckPassword(password, hashed) {
		t.Fatal("expected CheckPassword to return true for correct password")
	}
}

func TestCheckPassword_Wrong(t *testing.T) {
	password := "mysecretpassword"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if CheckPassword("wrongpassword", hashed) {
		t.Fatal("expected CheckPassword to return false for wrong password")
	}
}
