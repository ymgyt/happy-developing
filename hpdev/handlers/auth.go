package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/ymgyt/happy-developing/hpdev/app"
	"github.com/ymgyt/happy-developing/hpdev/errors"
	"go.uber.org/zap"
)

// Auth -
type Auth struct {
	*base

	HTTPClient *http.Client
	JWT        *app.JWT
}

type authenticateRequest struct {
	Realm       string `validate:"required"`
	AccessToken string `validate:"required"`
}

// Authenticate -
func (a *Auth) Authenticate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	authReq, err := a.authenticateRequestFromQuery(r)
	if err != nil {
		a.json(&apiResponse{W: w, Err: err})
		return
	}

	var email string
	switch r := strings.ToLower(authReq.Realm); r {
	case "github":
		email, err = (&githubService{accessToken: authReq.AccessToken, a: a}).email()
	default:
		a.json(&apiResponse{W: w, Err: errors.InvalidInput("unexpected realm", nil)})
		return
	}

	// TODO handle error
	if err != nil {
		a.json(&apiResponse{W: w, Err: err})
		return
	}
	if email == "" {
		panic("empty email")
	}

	idToken, err := a.JWT.Sign(a.JWT.StandardClaimsWithEmail(email))
	if err != nil {
		panic(err)
	}

	// どうにかredirect先を渡せないものか.今のところstateの中にしこんでgithubから返してもらう方法しか思いつかない.
	w.Header().Set("Location", "/author/posts?id_token="+idToken)
	w.WriteHeader(http.StatusFound)
}

type githubService struct {
	accessToken string
	a           *Auth
}

type githubAuthenticatedUserResponse struct {
	Type  string `json:"type"`
	Login string `json:"login"`
	Email string `json:"email"`
}

func (gs *githubService) email() (string, error) {
	// access_tokenに対応したユーザ情報を返す
	const ep = "https://api.github.com/user"

	req, err := http.NewRequest(http.MethodGet, ep, nil)
	if err != nil {
		return "", err
	}

	// https://developer.github.com/v3/#current-version
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", "token "+gs.accessToken)

	res, err := gs.a.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}

	var payload githubAuthenticatedUserResponse
	if err = json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return "", err
	}
	gs.a.Env.Log.Debug("authenticate", zap.String("sp", "github"), zap.Reflect("user_api", payload))

	return payload.Email, nil
}

func (a *Auth) authenticateRequestFromQuery(r *http.Request) (*authenticateRequest, error) {
	authReq := &authenticateRequest{
		Realm:       r.URL.Query().Get("realm"),
		AccessToken: r.URL.Query().Get("access_token"),
	}

	return authReq, a.Env.Validator.Validate(authReq)
}
