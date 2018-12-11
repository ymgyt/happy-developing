package middlewares

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/ymgyt/happy-developing/hpdev/app"
	"go.uber.org/zap"
)

// Authorizer -
type Authorizer struct {
	Env   *app.Env
	Email string
}

func (a *Authorizer) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	token, found := app.GetIDToken(r.Context())
	if !found {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		panic("only map claims is supported now")
	}

	got, want := mapClaims["sub"].(string), a.Email
	if got != want {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	a.Env.Log.Debug("authorize success", zap.String("claims_subject", got), zap.String("want", a.Email))
	next(w, r)
}
