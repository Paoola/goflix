package main

import (
	"log"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
)

func logRequestMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%v] %v", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	}
}

func (s *server) loggedOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jd := jwtmiddleware.New(jwtmiddleware.Options{
			ValidationKeyGetter: func(*jwt.Token) (interface{}, error) {
				return []byte(JWT_APP_KEY), nil
			},
			SigningMethod: jwt.SigningMethodHS256,
		})
		jd.HandlerWithNext(w, r, next)
	}
}
