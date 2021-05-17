package apiserver

import (
	"encoding/json"
	"errors"
	"github.com/astanishevskyi/http-server/internal/apiserver/adapters"
	"github.com/astanishevskyi/http-server/internal/apiserver/configs"
	"github.com/astanishevskyi/http-server/internal/apiserver/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Server struct {
	config  *configs.Config
	storage models.UserService
}

func New(config *configs.Config) *Server {
	return &Server{
		config: config,
	}
}

func (s *Server) Run() error {
	log.Println("Server is running on " + s.config.BindAddr)
	return http.ListenAndServe(s.config.BindAddr, nil)
}

func (s *Server) ConfigStorage() error {
	if s.config.Storage == "in-memory" {
		log.Println("Storage is in-memory")
		inMemoryStorage := adapters.NewInMemoryUserStorage()
		s.storage = inMemoryStorage
		return nil

	}
	return errors.New("no storage is set")
}

func (s *Server) ConfigRouter() {
	http.HandleFunc("/user/", s.User)
}

func (s *Server) User(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:

		id, _ := strconv.ParseUint(strings.TrimPrefix(r.URL.Path, "/user/"), 0, 32)

		if id == 0 {
			log.Println("GET /user/")
			resp := s.storage.GetAll()
			respJson, _ := json.Marshal(resp)
			_, err := w.Write(respJson)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			log.Printf("GET /user/%d", id)
			resp, err := s.storage.Retrieve(uint32(id))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			respJson, _ := json.Marshal(resp)
			_, err = w.Write(respJson)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	case http.MethodPost:
		log.Println("POST /user/")
		if err := r.ParseMultipartForm(10000); err != nil {
			http.Error(w, "wrong type of request, use multipart", http.StatusBadRequest)
			return
		}
		name := r.PostFormValue("name")
		email := r.PostFormValue("email")
		ageString := r.PostFormValue("age")
		age, err := strconv.ParseUint(ageString, 10, 8)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if (name == "") || (email == "") || (age == 0) {
			http.Error(w, "no name, email or age in body request", http.StatusBadRequest)
			return
		}

		resp := s.storage.Add(name, email, uint8(age))
		respJson, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(respJson)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case http.MethodPut:
		id, _ := strconv.ParseUint(strings.TrimPrefix(r.URL.Path, "/user/"), 0, 32)
		log.Printf("PUT /user/%d", id)

		if id == 0 {
			http.Error(w, "no user id", http.StatusBadRequest)
			return
		}
		name := r.PostFormValue("name")
		email := r.PostFormValue("email")
		ageString := r.PostFormValue("age")
		age, err := strconv.ParseUint(ageString, 10, 8)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp, err := s.storage.Update(uint32(id), name, email, uint8(age))
		if err != nil {
			http.Error(w, "no user id", http.StatusBadRequest)
			return
		}
		respJson, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(respJson)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodDelete:
		id, _ := strconv.ParseUint(strings.TrimPrefix(r.URL.Path, "/user/"), 0, 32)
		log.Printf("GET /user/%d", id)

		if id == 0 {
			http.Error(w, "no user id", http.StatusBadRequest)
			return
		}
		res, err := s.storage.Remove(uint32(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		marshal, err := json.Marshal(map[string]uint32{"id": res})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(marshal)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
