package domain

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/rokafela/udemy-banking-auth/errs"
	"github.com/rokafela/udemy-banking-auth/logger"
)

type AuthToken struct {
	token *jwt.Token
}

const HMAC_SECRET = "test"
const ACCESS_TOKEN_DURATION = time.Hour

type AuthTokenClaims struct {
	Username   string   `json:"username"`
	Role       string   `json:"role"`
	CustomerId string   `json:"customer_id"`
	AccountId  []string `json:"account_id"`
	jwt.StandardClaims
}

func GetClaimsForRole(login *Login) AuthTokenClaims {
	if strings.ToLower(login.Role) == "user" {
		return claimsForUser(login)
	} else {
		return claimsForAdmin(login)
	}
}

func claimsForUser(login *Login) AuthTokenClaims {
	accounts := strings.Split(login.AccountId.String, ",")
	return AuthTokenClaims{
		CustomerId: login.CustomerId.String,
		AccountId:  accounts,
		Username:   login.Username,
		Role:       login.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
		},
	}
}

func claimsForAdmin(login *Login) AuthTokenClaims {
	return AuthTokenClaims{
		Username: login.Username,
		Role:     login.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ACCESS_TOKEN_DURATION).Unix(),
		},
	}
}

func GenerateAuthToken(claims *AuthTokenClaims) (string, *errs.AppError) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(HMAC_SECRET))
	if err != nil {
		logger.Error("Failed while signing auth token: " + err.Error())
		return "", errs.NewUnexpectedError("cannot generate auth token")
	}
	return signedString, nil
}
