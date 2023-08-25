package pages

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/wakieu/drtbox/api"
	"github.com/wakieu/drtbox/database"
)

type PageHandler struct {
	boxRepo *database.BoxRepository
	landT   *template.Template
	boxT    *template.Template
}

func NewTemplate(filepath string) (*template.Template, error) {
	pageTemplateByte, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	pageTemplate := template.Must(template.New("").Parse(string(pageTemplateByte)))

	return pageTemplate, nil
}

func NewHandler(boxRepository *database.BoxRepository, lt, bt *template.Template) *PageHandler {
	return &PageHandler{
		boxRepo: boxRepository,
		landT:   lt,
		boxT:    bt,
	}
}

func (h *PageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not alllowed.", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "/" {
		h.landT.Execute(w, nil)
	} else {
		res, err := http.Get("http://localhost:3131/" + r.URL.Path)
		if err != nil {
			log.Printf("%v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
			return
		}
		var getBoxResponse api.GetBoxResponse
		err = json.NewDecoder(res.Body).Decode(&getBoxResponse)
		if err != nil {
			log.Printf("%+v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
			return
		}
		h.boxT.Execute(w, getBoxResponse)
	}

}
