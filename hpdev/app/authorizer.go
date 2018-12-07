package app

// Authorizer -
type Authorizer interface {
	Authorize(*AuthorizeRequest) (*AuthorizeResponse, error)
}

// AuthorizeRequest -
type AuthorizeRequest struct {
	Realm    Realm
	Email    string
	Password string
}

// AuthorizeResponse -
type AuthorizeResponse struct {
	OK      bool
	Message string
}

// Realm -
type Realm string

const (
	// InvalidRealm -
	InvalidRealm = ""
	// PasswordRealm -
	PasswordRealm = "password"
)
