package app_test

import (
	"testing"

	"github.com/ymgyt/happy-developing/hpdev/app"
	"github.com/ymgyt/happy-developing/hpdev/errors"
)

func TestEnvVarAuthorizer_Authorize(t *testing.T) {
	email := "gopher@example.com"
	password := "secret"

	tests := []struct {
		desc  string
		input *app.AuthorizeRequest
		want  *app.AuthorizeResponse
		err   error
	}{
		{
			desc: "authorize success",
			input: &app.AuthorizeRequest{
				Realm:    app.PasswordRealm,
				Email:    email,
				Password: password,
			},
			want: &app.AuthorizeResponse{
				OK: true,
			},
		},
		{
			desc: "authorize fail",
			input: &app.AuthorizeRequest{
				Realm:    app.PasswordRealm,
				Email:    email,
				Password: "invalid",
			},
			want: &app.AuthorizeResponse{
				OK:      false,
				Message: "wrong email/password",
			},
		},
		{
			desc: "wrong realm",
			input: &app.AuthorizeRequest{
				Realm: "",
			},
			err: errors.New(errors.InvalidInputErr, ""),
		},
	}

	auth := app.EnvVarAuthorizer{Email: email, Password: password}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := auth.Authorize(tc.input)
			CmpErr(t, err, tc.err)
			Cmp(t, got, tc.want)
		})
	}
}
