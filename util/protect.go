package util

import (
	"net/http"
)

type Session struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Token     string `json:"token"`
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request, Session)
}

type HandlerFunc func(http.ResponseWriter, *http.Request, Session)

func (that HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request, session Session) {
	that(w, r, session)
}

func Protect(handler Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := Session{}
		handler.ServeHTTP(w, r, session)
	})
}
