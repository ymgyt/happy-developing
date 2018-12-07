package app

import (
	"strings"

	"github.com/ymgyt/happy-developing/hpdev/errors"
)

// EnvVarAuthorizer -
type EnvVarAuthorizer struct {
	Email    string
	Password string
}

// Authorize -
func (a *EnvVarAuthorizer) Authorize(req *AuthorizeRequest) (*AuthorizeResponse, error) {
	if req.Realm != PasswordRealm {
		return nil, errors.InvalidInput("want password realm", nil)
	}

	if strings.Compare(req.Email, a.Email) != 0 || strings.Compare(req.Password, a.Password) != 0 {
		return &AuthorizeResponse{OK: false, Message: "wrong email/password"}, nil
	}

	return &AuthorizeResponse{OK: true}, nil
}
