package handlers

import (
	"context"
	"testing"

	"github.com/puricalvo/twitterGo/models"
)

func TestRegistrationRouteDoesNotRequireToken(t *testing.T) {
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
