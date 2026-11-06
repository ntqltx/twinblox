package keychain

import (
	"bytes"
	"errors"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	const account = "test"
	if err := Set(account, []byte("secret")); err != nil {
		t.Fatal(err)
	}

	got, err := Get(account)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, []byte("secret")) {
		t.Fatalf("got %q", got)
	}

	if err := Delete(account); err != nil {
		t.Fatal(err)
	}

	if _, err := Get(account); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
