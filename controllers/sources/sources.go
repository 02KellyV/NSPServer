package sources

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/hhrutter/pdfcpu"
	"github.com/juliotorresmoreno/neosmarthpen/util"
)

type Router struct {
	http.Handler
}

const filesPath = "files"

func NewRouter() http.Handler {
	route := mux.NewRouter()

	route.HandleFunc("/", util.ListFiles(path.Join(filesPath, "sources"))).Methods("GET")
	route.HandleFunc("/compile/{image}", Compile).Methods("GET")
	route.HandleFunc("/convert/{image}", Convert).Methods("GET")
	route.HandleFunc("/view/{image}", View).Methods("GET")
	route.HandleFunc("/details/{image}", Details).Methods("GET")

	return route
}

func Details(w http.ResponseWriter, r *http.Request) {
	image := mux.Vars(r)["image"]
	_image := image[0 : len(image)-4]
	output := path.Join(filesPath, "output", _image)
	files, _ := ioutil.ReadDir(output)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"details": map[string]interface{}{
			"length": len(files),
		},
	})
}

func View(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	image := mux.Vars(r)["image"]

	_image := image[0 : len(image)-4]
	file := path.Join(filesPath, "output", _image, _image+"_"+page+".jpg")

	http.ServeFile(w, r, file)
}

func convert(image, page string) error {
	file := path.Join(filesPath, "output", image, image+"_"+page+".pdf")
	imageName := path.Join(filesPath, "output", image, image+"_"+page+".jpg")

	if err := util.ConvertPdfToJpg(file, imageName); err != nil {
		return err
	}
	return nil
}

func Convert(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	image := mux.Vars(r)["image"]

	_image := image[0 : len(image)-4]

	convert(_image, page)
}

func Compile(w http.ResponseWriter, r *http.Request) {
	image := mux.Vars(r)["image"]
	_image := image[0 : len(image)-4]
	file := path.Join(filesPath, "sources", image)
	output := path.Join(filesPath, "output", _image)
	os.MkdirAll(output, 0755)
	if _, err := os.Stat(file); err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}
	err := pdfcpu.ExtractPages(file, output, nil, nil)
	if err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	files, err := ioutil.ReadDir(output)
	if err != nil {
		util.RenderError(w, http.StatusInternalServerError, err)
		return
	}

	for _, file := range files {
		name := file.Name()
		fileName := path.Join(filesPath, output, file.Name())
		imageName := path.Join(filesPath, output, name[0:len(name)-4]+".jpg")

		if err := util.ConvertPdfToJpg(fileName, imageName); err != nil {
			if err != nil {
				util.RenderError(w, http.StatusInternalServerError, err)
				return
			}
		}
		os.Remove(fileName)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}
