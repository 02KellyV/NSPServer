package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/juliotorresmoreno/neosmarthpen/config"
	"github.com/juliotorresmoreno/neosmarthpen/controllers/auth"
	"github.com/juliotorresmoreno/neosmarthpen/controllers/files"
	"github.com/juliotorresmoreno/neosmarthpen/controllers/sources"
	"github.com/juliotorresmoreno/neosmarthpen/middlewares"
	"github.com/juliotorresmoreno/unravel-server/helper"
)

func NewRouter(config config.Config) http.Handler {
	router := mux.NewRouter()

	router.Use(middlewares.Log)

	router.PathPrefix("/api/sources").Handler(
		helper.StripPrefix(
			"/api/sources",
			sources.NewRouter(),
		),
	)
	router.PathPrefix("/api/files").Handler(
		helper.StripPrefix(
			"/api/files",
			files.NewRouter(),
		),
	)
	router.PathPrefix("/api/auth").Handler(
		helper.StripPrefix(
			"/api/auth",
			auth.NewRouter(config),
		),
	)

	return router
}
