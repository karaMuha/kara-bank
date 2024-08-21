package middlewares

import (
	"context"
	"kara-bank/utils"
	"net/http"
	"strings"
)

type contextUserEmail string

const contextUserEmailKey contextUserEmail = "email"

func AuthMiddleware(tokenMaker utils.TokenMaker, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestTarget := r.Method + " " + strings.Split(r.URL.Path, "/")[1]

		if !utils.IsProtectedRoute(requestTarget) {
			next.ServeHTTP(w, r)
			return
		}

		authToken, err := r.Cookie("paseto")

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		verifiedToken, err := tokenMaker.VerifyToken(authToken.Value)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), contextUserEmailKey, verifiedToken.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
