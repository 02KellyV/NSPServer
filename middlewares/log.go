package middlewares

import (
	"fmt"
	"net/http"

	"github.com/onna-soft/onna-soft/helper"
)

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method, r.URL.Path)
		helper.Cors(w, r)
		if r.Method != "OPTIONS" {
			handler.ServeHTTP(w, r)
		}
	})
}
