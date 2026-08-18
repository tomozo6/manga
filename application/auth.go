package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"firebase.google.com/go/v4/auth"
)

type identity struct {
	UID   string
	Name  string
	Email string
}

type tokenVerifier interface {
	Verify(context.Context, string) (identity, error)
}

type firebaseVerifier struct{ client *auth.Client }

func (v firebaseVerifier) Verify(ctx context.Context, rawToken string) (identity, error) {
	token, err := v.client.VerifyIDToken(ctx, rawToken)
	if err != nil {
		return identity{}, err
	}
	email, _ := token.Claims["email"].(string)
	name, _ := token.Claims["name"].(string)
	if email == "" {
		return identity{}, errors.New("Firebase ID token does not include an email address")
	}
	return identity{UID: token.UID, Name: name, Email: email}, nil
}

func (a app) authenticate(w http.ResponseWriter, r *http.Request) (identity, bool) {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "authentication is required"})
		return identity{}, false
	}
	user, err := a.verifier.Verify(r.Context(), strings.TrimPrefix(header, "Bearer "))
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid Firebase ID token"})
		return identity{}, false
	}
	if _, ok := a.allowed[strings.ToLower(user.Email)]; !ok {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "this account is not allowed"})
		return identity{}, false
	}
	return user, true
}
