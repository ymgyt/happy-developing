package middlewares

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/ymgyt/happy-developing/hpdev/app"
)

// JWTVerifier -
type JWTVerifier struct {
	JWT  *app.JWT
	next http.Handler
}

// ServeHTTP -
func (m *JWTVerifier) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rawToken := m.readToken(r)
	if rawToken == "" {
		m.next.ServeHTTP(w, r)
		return
	}

	token, err := jwt.Parse(rawToken, m.keyFunc)
	if err != nil {
		// TODO logging
		m.next.ServeHTTP(w, r)
		return
	}

	rr := r.WithContext(app.SetIDToken(r.Context(), token))
	m.next.ServeHTTP(w, rr)
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

// SetNext -
func (m *JWTVerifier) SetNext(next http.Handler) {
	m.next = next
}
