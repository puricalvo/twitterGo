package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/puricalvo/twitterGo/models"
)

func TestRegistrationRouteDoesNotRequireToken(t *testing.T) {
	t.Setenv("ALLOW_REGISTRATION", "true")
	t.Setenv("DEMO_READ_ONLY", "false")
	ctx := context.Background()
	ctx = context.WithValue(ctx, models.Key("path"), "/registro")
	ctx = context.WithValue(ctx, models.Key("method"), "POST")
	ctx = context.WithValue(ctx, models.Key("body"), "")
	ctx = context.WithValue(ctx, models.Key("jwtSign"), "test-secret")

	response := Manejadores(ctx, models.APIRequest{})

	if response.Message == "Token requerido" {
		t.Fatal("registration route must be public")
	}
}

func TestRegistrationRouteCanBeDisabled(t *testing.T) {
	t.Setenv("ALLOW_REGISTRATION", "false")
	ctx := context.Background()
	ctx = context.WithValue(ctx, models.Key("path"), "/registro")
	ctx = context.WithValue(ctx, models.Key("method"), "POST")
	ctx = context.WithValue(ctx, models.Key("body"), "")
	ctx = context.WithValue(ctx, models.Key("jwtSign"), "test-secret")

	response := Manejadores(ctx, models.APIRequest{})

	if response.Status != 403 {
		t.Fatalf("status = %d, want %d", response.Status, 403)
	}
	if response.Message != "El registro público está desactivado" {
		t.Fatalf("message = %q, want registration-disabled message", response.Message)
	}
}

func TestRegistrationRouteRejectsInvalidConfiguration(t *testing.T) {
	t.Setenv("ALLOW_REGISTRATION", "not-a-boolean")
	t.Setenv("DEMO_READ_ONLY", "false")
	ctx := context.Background()
	ctx = context.WithValue(ctx, models.Key("path"), "/registro")
	ctx = context.WithValue(ctx, models.Key("method"), "POST")
	ctx = context.WithValue(ctx, models.Key("body"), "")
	ctx = context.WithValue(ctx, models.Key("jwtSign"), "test-secret")

	response := Manejadores(ctx, models.APIRequest{})

	if response.Status != 500 {
		t.Fatalf("status = %d, want %d", response.Status, 500)
	}
}

func TestReadOnlyDemoRejectsWrites(t *testing.T) {
	t.Setenv("ALLOW_REGISTRATION", "true")
	t.Setenv("DEMO_READ_ONLY", "true")
	ctx := context.Background()
	ctx = context.WithValue(ctx, models.Key("path"), "/registro")
	ctx = context.WithValue(ctx, models.Key("method"), "POST")
	ctx = context.WithValue(ctx, models.Key("body"), "")
	ctx = context.WithValue(ctx, models.Key("jwtSign"), "test-secret")

	response := Manejadores(ctx, models.APIRequest{})

	if response.Status != 403 {
		t.Fatalf("status = %d, want %d", response.Status, 403)
	}
	if response.Message != "La demo está en modo de solo lectura" {
		t.Fatalf("message = %q, want read-only message", response.Message)
	}
}

func TestReadOnlyWriteRequestClassification(t *testing.T) {
	for _, test := range []struct {
		method string
		path   string
		want   bool
	}{
		{method: http.MethodPost, path: "login", want: false},
		{method: http.MethodGet, path: "leoTweets", want: false},
		{method: http.MethodPost, path: "tweet", want: true},
		{method: http.MethodPut, path: "modificarPerfil", want: true},
		{method: http.MethodDelete, path: "eliminarTweet", want: true},
	} {
		if got := isWriteRequest(test.method, test.path); got != test.want {
			t.Errorf("isWriteRequest(%q, %q) = %t, want %t", test.method, test.path, got, test.want)
		}
	}
}

func TestReadOnlyRejectsInvalidConfiguration(t *testing.T) {
	t.Setenv("DEMO_READ_ONLY", "not-a-boolean")
	ctx := context.Background()
	ctx = context.WithValue(ctx, models.Key("path"), "/registro")
	ctx = context.WithValue(ctx, models.Key("method"), "POST")
	ctx = context.WithValue(ctx, models.Key("body"), "")
	ctx = context.WithValue(ctx, models.Key("jwtSign"), "test-secret")

	response := Manejadores(ctx, models.APIRequest{})

	if response.Status != 500 {
		t.Fatalf("status = %d, want %d", response.Status, 500)
	}
}
