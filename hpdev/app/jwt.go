package app

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// JWT -
type JWT struct {
	HMACSecret []byte
}

// Sign -
func (j *JWT) Sign(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.HMACSecret)
}

// StandardClaimsWithEmail -
func (j *JWT) StandardClaimsWithEmail(email string) *jwt.StandardClaims {
	now := time.Now()
	claims := &jwt.StandardClaims{
		// Id: "id",
		Audience:  "happy-developing.io",
		ExpiresAt: now.Add(time.Hour * 3).Unix(),
		IssuedAt:  now.Unix(),
		Issuer:    "happy-developing.io",
		NotBefore: now.Unix(),
		Subject:   email,
	}
	return claims
}

type contextKey string

const (
	idTokenContextKey contextKey = "idToken"
)

// SetIDToken -
func SetIDToken(ctx context.Context, token *jwt.Token) context.Context {
	return context.WithValue(ctx, idTokenContextKey, token)
}

// GetIDToken -
func GetIDToken(ctx context.Context) (*jwt.Token, bool) {
	maybeToken := ctx.Value(idTokenContextKey)
	token, ok := maybeToken.(*jwt.Token)
	return token, ok
}
