package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"

	"github.com/ymgyt/happy-developing/hpdev/oauth2"
	"github.com/ymgyt/happy-developing/hpdev/view"
)

const (
	// defaultのlogin後の遷移先
	defaultRedirect = "/author/posts"
)

// OAuth2 -
type OAuth2 struct {
	Config     *oauth2.Config
	HTTPClient *http.Client
	*base
	ts *templateSet
}

// RenderLogin -
func (o *OAuth2) RenderLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	err := o.ts.ExecuteTemplate(w, "login", &view.Login{Config: o.Config})
	o.handleRenderError(err)
}

// GithubCallback -
func (o *OAuth2) GithubCallback(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log := o.Env.Log.With(zap.String("sp", "github"))

	// check state
	state := r.URL.Query().Get("state")
	if strings.Compare(state, o.Config.CSRFToken) != 0 {
		log.Error("oauth2", zap.Error(fmt.Errorf("state does not match. got=%s, want=%s", state, o.Config.CSRFToken)))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		panic("empty code")
	}
	log.Debug("oauth2", zap.String("code", code))

	// fetch access token
	c := o.Config.Github
	tokenEndpoint := fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s&state=%s", c.TokenURL, c.ClientID, c.ClientSecret, code, state)

	tokenReq, err := http.NewRequest(http.MethodPost, tokenEndpoint, nil)
	if err != nil {
		log.Error("oauth2", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenReq.Header.Set("Accept", "application/json")

	tokenRes, err := o.HTTPClient.Do(tokenReq)
	if err != nil {
		log.Error("oauth2", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer tokenRes.Body.Close()

	var accessTokenPayload oauth2.AccessTokenResponse
	if err := json.NewDecoder(tokenRes.Body).Decode(&accessTokenPayload); err != nil {
		log.Error("oauth2", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Debug("oauth2", zap.Reflect("token_response", accessTokenPayload))

	// now we can use access token like this
	// GET https://api.github.com/user?access_token=...
	// curl -H "Authorization: token OAUTH-TOKEN" https://api.github.com/user

	// redirect user with access token
	w.Header().Set("Location", "/authenticate?realm=github&access_token="+accessTokenPayload.AccessToken)
	w.WriteHeader(http.StatusFound)
}
