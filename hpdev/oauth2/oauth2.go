package oauth2

// Config -
type Config struct {
	Github    *Entry
	CSRFToken string // handler側で管理すべき..?
}

// Entry -
type Entry struct {
	*Endpoint
	*Credential
	CallbackURL string
}

// Endpoint -
type Endpoint struct {
	AuthorizeURL string
	TokenURL     string
}

// Credential -
type Credential struct {
	ClientID     string
	ClientSecret string
}

// AccessTokenResponse -
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}
