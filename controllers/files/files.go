package files

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gorilla/mux"
	"github.com/juliotorresmoreno/neosmarthpen/db"
	"github.com/juliotorresmoreno/neosmarthpen/util"
	"github.com/juliotorresmoreno/unravel-server/helper"
)

type Router struct {
	http.Handler
}

func NewRouter() http.Handler {
	route := mux.NewRouter()

	route.HandleFunc("/", util.ListFiles("files")).Methods("GET")
	route.HandleFunc("/upload", Upload).Methods("POST")

	return route
}

func View(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	image := mux.Vars(r)["image"]

	_image := image[0 : len(image)-4]
	file := path.Join("files", _image, _image+"_"+page+".jpg")

	http.ServeFile(w, r, file)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	client, err := db.NewCache()
	if err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}
	defer client.Close()
	c, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}
	ramdon := helper.GenerateRandomString(20)

	name := fmt.Sprintf("files/%v.pdf", ramdon)
	if err := ioutil.WriteFile(name, c, 0644); err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
