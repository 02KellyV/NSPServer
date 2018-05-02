package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-xorm/xorm"
	"github.com/gorilla/mux"
	"github.com/juliotorresmoreno/neosmarthpen/config"
	"github.com/juliotorresmoreno/neosmarthpen/db"
	"github.com/juliotorresmoreno/neosmarthpen/models"
	"github.com/juliotorresmoreno/neosmarthpen/util"
	"github.com/juliotorresmoreno/unravel-server/helper"
)

type Router struct {
	*mux.Router
	config.Config
}

func NewRouter(config config.Config) http.Handler {
	route := Router{
		Router: mux.NewRouter(),
		Config: config,
	}

	route.HandleFunc("/", util.ListFiles("files")).Methods("GET")
	route.HandleFunc("/sign-in", route.signIn).Methods("POST")
	route.HandleFunc("/sign-up", route.signUp).Methods("POST")

	sessionFunc := util.Protect(util.HandlerFunc(route.session))
	route.Handle("/session", sessionFunc).Methods("GET")

	logoutFunc := util.Protect(util.HandlerFunc(route.logout))
	route.Handle("/logout", logoutFunc).Methods("GET")

	return route
}

func (that Router) logout(w http.ResponseWriter, r *http.Request, session util.Session) {
	token := helper.GetCookie(r, "token")
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	c, err := db.NewCache()
	if err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}
	defer c.Close()

	c.Set(token, "", time.Millisecond*1)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func (that Router) signUp(w http.ResponseWriter, r *http.Request) {
	data := helper.GetPostParams(r)
	email := data.Get("email")
	password := data.Get("password")
	firstName := data.Get("first_name")
	lastName := data.Get("last_name")

	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
	}

	if ok, err := helper.ValidateStruct(&user); !ok {
		var errors govalidator.Errors
		switch err.(type) {
		case govalidator.Errors:
			errors = err.(govalidator.Errors)
		default:
			errors = govalidator.Errors{err}
		}
		util.RenderErrors(w, http.StatusNotAcceptable, errors)
		return
	}
	user.Password = helper.Encript(user.Password)

	conn, err := db.NewConn(that.Config)
	if err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()

	if _, err = conn.Insert(user); err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}
	signIn(w, conn, email, password)
}

func (that Router) signIn(w http.ResponseWriter, r *http.Request) {
	data := helper.GetPostParams(r)
	email := data.Get("email")
	password := data.Get("password")

	conn, err := db.NewConn(that.Config)
	if err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()

	signIn(w, conn, email, password)
}

func (that Router) session(w http.ResponseWriter, r *http.Request, session util.Session) {
	token := helper.GetCookie(r, "token")
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	c, err := db.NewCache()
	if err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}
	defer c.Close()

	userID, err := c.Get(token).Result()
	if err != nil {
		util.RenderError(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	conn, err := db.NewConn(that.Config)
	if err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()

	user := models.User{}
	_, err = conn.Where("id = ?", userID).Get(&user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    user,
		"token":   token,
	})
}

func signIn(w http.ResponseWriter, conn *xorm.Engine, email, password string) {
	user := models.User{}
	_, err := conn.Where("email = ?", email).Get(&user)
	if err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	if user.ID == 0 || !helper.IsValid(user.Password, password) {
		util.RenderError(w, http.StatusUnauthorized, errors.New("Usuario o contraseña no valido"))
		return
	}
	cache, err := db.NewCache()
	defer cache.Close()
	if err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	token := helper.GenerateRandomString(40)

	cache.Set(token, string(user.ID), 2*time.Hour)

	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{
		MaxAge:   7200,
		Secure:   false,
		HttpOnly: true,
		Name:     "token",
		Value:    token,
		Path:     "/",
	})
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    user,
		"token":   token,
	})
}
