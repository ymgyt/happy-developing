package app

// Services -
type Services struct {
	PostService PostService
	TagService  TagService
	JWTService  *JWT
}
