package utils

import (
	"context"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		token := r.Header.Get("Authorization")
		if token != "" {
			ctx = context.WithValue(ctx,"authToken", token)
		}
		next.ServeHTTP(w, r.WithContext(ctx))

		//response, err := Redis.Do("GET", token)
		//if err != nil || response == nil {
		//	if err != nil {
		//		w.WriteHeader(http.StatusInternalServerError)
		//	}
		//	if response == nil {
		//		w.WriteHeader(http.StatusUnauthorized)
		//	}
		//} else {
		//	_, err := Redis.Do("EXPIRE", token, 3600)
		//	if err != nil {
		//		w.WriteHeader(http.StatusInternalServerError)
		//	}
		//	next.ServeHTTP(w, r)
		//}
	})
}
