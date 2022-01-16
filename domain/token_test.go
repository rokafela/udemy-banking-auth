package domain

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_get_claims_user(t *testing.T) {
	// Arrange
	var test_login = Login{
		Username: "2000",
		Role:     "user",
		CustomerId: sql.NullString{
			String: "2000",
			Valid:  true,
		},
		AccountId: sql.NullString{
			String: "95470",
			Valid:  true,
		},
	}

	// Act
	var claims = GetClaimsForRole(&test_login)

	// Assert
	if claims.Username != test_login.Username {
		t.Error("Invalid username while testing user auth token claims")
	}
	if claims.Role != test_login.Role {
		t.Error("Invalid role while testing user auth token claims")
	}
	if claims.CustomerId != test_login.CustomerId.String {
		t.Error("Invalid customer id while testing user auth token claims")
	}
	accounts := strings.Split(test_login.AccountId.String, ",")
	if !cmp.Equal(claims.AccountId, accounts) {
		t.Error("Invalid account id while testing user auth token claims")
	}
	if claims.StandardClaims.ExpiresAt == 0 {
		t.Error("Invalid standard claim expires at while testing user auth token claims")
	}
}

func Test_get_claims_admin(t *testing.T) {
	// Arrange
	var test_login = Login{
		Username: "admin",
		Role:     "admin",
	}

	// Act
	var claims = GetClaimsForRole(&test_login)

	// Assert
	if claims.Username != test_login.Username {
		t.Error("Invalid username while testing admin auth token claims")
	}
	if claims.Role != test_login.Role {
		t.Error("Invalid role while testing admin auth token claims")
	}
	if claims.StandardClaims.ExpiresAt == 0 {
		t.Error("Invalid standard claim expires at while testing admin auth token claims")
	}
}
