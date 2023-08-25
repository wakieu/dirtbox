package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/wakieu/drtbox/database"
	"github.com/wakieu/drtbox/entity"
)

type ApiHandler struct {
	boxRepo *database.BoxRepository
}

type GetBoxResponse struct {
	Box      entity.Box `json:"box"`
	Children []string   `json:"children"`
}

func NewHandler(boxRepository *database.BoxRepository) *ApiHandler {
	return &ApiHandler{
		boxRepo: boxRepository,
	}
}

func (h *ApiHandler) GetBox(w http.ResponseWriter, r *http.Request) {
	path := cleanPath(r.URL.Path)
	box, err := h.boxRepo.GetContent(path)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}

	children, err := h.boxRepo.GetChildren(path)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}

	getBoxResponse := GetBoxResponse{
		box,
		children,
	}

	res, err := json.Marshal(getBoxResponse)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}

	w.Write(res)
}

func (h *ApiHandler) WriteBox(w http.ResponseWriter, r *http.Request) {
	path := cleanPath(r.URL.Path)
	var box entity.Box
	err := json.NewDecoder(r.Body).Decode(&box)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, "Bad request...", http.StatusBadRequest)
		return
	}
	box.BoxPath = cleanPath(box.BoxPath)

	ok, err := h.boxRepo.Exists(path)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}

	if ok {
		if box.IsEmpty() {
			err = h.boxRepo.Delete(box.BoxPath)
			if err != nil {
				log.Printf("%+v", err)
				http.Error(w, "Something went wrong...", http.StatusInternalServerError)
				return
			}
		} else {
			err = h.boxRepo.Write(box.BoxPath, box.Text)
			if err != nil {
				log.Printf("%+v", err)
				http.Error(w, "Something went wrong...", http.StatusInternalServerError)
				return
			}
		}
	} else {
		err = h.boxRepo.Save(&box)
		if err != nil {
			log.Printf("%+v", err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
			return
		}
	}

	json.NewEncoder(w).Encode(EasyResponse("OK"))
}

func (h *ApiHandler) DeleteBox(w http.ResponseWriter, r *http.Request) {
	path := cleanPath(r.URL.Path)
	err := h.boxRepo.Delete(path)
	if err != nil {
		log.Printf("%+v", err)
		http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(EasyResponse("OK"))
}

func EasyResponse(msg string) string {
	m := make(map[string]string)
	m["msg"] = msg
	return fmt.Sprintf("%+v", m)
}

func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	switch r.Method {
	case "GET":
		h.GetBox(w, r)
	case "POST":
		h.WriteBox(w, r)
	case "DELETE":
		h.DeleteBox(w, r)
	case "OPTIONS":
		w.Write([]byte("cors OKOK"))
	default:
		http.Error(w, "Method not alllowed.", http.StatusMethodNotAllowed)
	}
}

func cleanPath(s string) string {
	if s[len(s)-1] == '/' {
		s = s[:len(s)-1]
	}
	return s
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:3030")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
