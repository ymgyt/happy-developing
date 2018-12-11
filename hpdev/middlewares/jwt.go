package middlewares

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"

	"github.com/ymgyt/happy-developing/hpdev/app"
)

// JWTVerifier -
type JWTVerifier struct {
	Env *app.Env
	JWT *app.JWT
}

// ServeHTTP -
func (m *JWTVerifier) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rawToken := m.readToken(r)
	if rawToken == "" {
		// m.next.ServeHTTP(w, r)
		m.Env.Log.Debug("jwt token nof found")
		next(w, r)
		return
	}

	token, err := jwt.Parse(rawToken, m.keyFunc)
	if err != nil {
		// m.next.ServeHTTP(w, r)
		m.Env.Log.Debug("jwt parse failed", zap.String("err", err.Error()))
		next(w, r)
		return
	}
	m.Env.Log.Debug("jwt parse success", zap.Any("claims", token.Claims))

	rr := r.WithContext(app.SetIDToken(r.Context(), token))
	// m.next.ServeHTTP(w, rr)
	next(w, rr)
}

func (m *JWTVerifier) keyFunc(token *jwt.Token) (interface{}, error) {
	// TODO more verify
	return m.JWT.HMACSecret, nil
}

func (m *JWTVerifier) readToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token != "" {
		return token
	}

	return r.URL.Query().Get("id_token")
}
