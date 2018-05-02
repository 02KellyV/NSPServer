package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
)

func RenderError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   err.Error(),
	})
}

func RenderErrors(w http.ResponseWriter, status int, errors []error) {
	w.WriteHeader(status)
	errs := map[string]interface{}{
		"_others": []string{},
	}
	for _, err := range errors {
		switch err.(type) {
		case govalidator.Error:
			name := err.(govalidator.Error).Name
			errs[strings.ToLower(name)] = err.Error()[len(name)+2:]
		default:
			errs["_others"] = append(errs["all"].([]string), err.Error())
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(map[string]interface{}{
		"success": false,
		"error":   errs,
	})
}

func ListFiles(path string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sources, err := ioutil.ReadDir(path)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
		}
		data := make([]string, len(sources))
		for i, v := range sources {
			data[i] = v.Name()
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    data,
		})
	}
}
