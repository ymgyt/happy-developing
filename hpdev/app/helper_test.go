package app_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ymgyt/happy-developing/hpdev/errors"
)

func CmpErr(t *testing.T, got, want error, opts ...cmp.Option) {
	t.Helper()
	if errors.GetCode(got) != errors.GetCode(want) {
		t.Fatalf("(-got +want)\n%s", cmp.Diff(got, want, opts...))
	}
}

func Cmp(t *testing.T, got, want interface{}, opts ...cmp.Option) {
	t.Helper()
	if diff := cmp.Diff(got, want, opts...); diff != "" {
		t.Fatalf("(-got +want)\n%s", diff)
	}
}
