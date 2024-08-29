package middlewares

import (
	"context"
	"kara-bank/utils"
	"net/http"
)

type contextUserEmail string
type contextUserRole string

const ContextUserEmailKey contextUserEmail = "userEmail"
const ContextUserRoleKey contextUserRole = "userRole"

func AuthMiddleware(tokenMaker utils.TokenMaker, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestTarget := r.Method + " " + r.URL.Path

		endpointRoles, err := utils.IsProtectedRoute(requestTarget)

		if err != nil {
			http.Error(w, "Could not find roles for this endpoint", http.StatusInternalServerError)
			return
		}

		if endpointRoles == nil {
			http.Error(w, "Could not find roles for this endpoint", http.StatusInternalServerError)
			return
		}

		if endpointRoles[0] == "" {
			next.ServeHTTP(w, r)
			return
		}

		authToken, err := r.Cookie("access_token")

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		verifiedToken, err := tokenMaker.VerifyToken(authToken.Value)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		for _, v := range endpointRoles {
			if verifiedToken.Role == v {
				ctx := context.WithValue(r.Context(), ContextUserEmailKey, verifiedToken.Email)
				ctx = context.WithValue(ctx, ContextUserRoleKey, verifiedToken.Role)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		}

		http.Error(w, "You do not have the right role to do this", http.StatusUnauthorized)
	})
}
