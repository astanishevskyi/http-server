package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/astanishevskyi/http-server/internal/apiserver/adapters"
	"github.com/astanishevskyi/http-server/internal/apiserver/configs"
	"github.com/astanishevskyi/http-server/internal/apiserver/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Server struct {
	config  *configs.Config
	router  *mux.Router
	storage models.UserService
}

func New(config *configs.Config) *Server {
	return &Server{
		config: config,
		router: mux.NewRouter(),
	}
}

func (s *Server) Run() error {
	log.Println("Server is running on " + s.config.BindAddr)
	srv := &http.Server{
		Addr:    s.config.BindAddr,
		Handler: s.router,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Println("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Print("Server Exited Properly")
	return nil
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
	s.router.HandleFunc("/user", s.GetUsers).Methods("GET")
	s.router.HandleFunc("/user", s.CreateUser).Methods("POST")
	s.router.HandleFunc("/user/{id:[0-9]+}", s.GetUser).Methods("GET")
	s.router.HandleFunc("/user/{id:[0-9]+}", s.UpdateUser).Methods("PUT")
	s.router.HandleFunc("/user/{id:[0-9]+}", s.DeleteUser).Methods("DELETE")
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 0, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("GET /user/%d", id)
	resp, err := s.storage.Retrieve(uint32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(respJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /user/")
	resp := s.storage.GetAll()
	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(respJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
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
	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(respJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 0, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
	respJSON, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(respJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 0, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
