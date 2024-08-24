package http

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	_ "http_server/docs"

	model "http_server/models"

	chi "github.com/go-chi/chi/v5"
	uuid "github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Storage interface {
	Get(key string) (*model.Task, error)
	Post(value model.Task) error
	Put(value model.Task) error
	Post_user(value model.User) error
	Get_user(value model.User) (*model.User, error)
	Post_session(value model.Session) error
	Get_session(key string) error
}

type Server struct {
	storage Storage
}

type regRequest struct {
	Login    string `json:"username"`
	Password string `json:"password"`
}

func newServer(curr_storage Storage) *Server {
	return &Server{storage: curr_storage}
}

func task_work(s *Server, id uuid.UUID) {
	time.Sleep(10 * time.Second)

	s.storage.Put(model.Task{ID: id.String(), Readiness: "ready", Result: "some rubish"})
}

func (s *Server) checkUnauthorization(w http.ResponseWriter, r *http.Request) bool {
	tmp := r.Header.Get("Authorization")

	if tmp == "" {
		http.Error(w, "401 Unauthorized", http.StatusUnauthorized)

		return false
	}

	bearerToken := strings.Split(tmp, " ")

	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		http.Error(w, "Invalid Authorization header", http.StatusBadRequest)

		return false
	}

	if err := s.storage.Get_session(bearerToken[1]); err != nil {
		http.Error(w, "Failed to get session", http.StatusBadRequest)

		return false
	}

	return true
}

// @Summary Post task task_id
// @Description Post task_id in DataBase and start task
// @Success 201 {json} json "task_id": "id value"
// @Failure 404 {string} string "Failed to store value"
// @Router /task [post]
func (s *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	s.checkUnauthorization(w, r)

	id, _ := uuid.NewUUID()

	if err := s.storage.Post(model.Task{ID: id.String(), Readiness: "in_progress", Result: ""}); err != nil {
		http.Error(w, "Failed to store value", http.StatusNotFound)

		return
	}

	go task_work(s, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"task_id": id.String()})
}

// @Summary Get task status
// @Description Post task_id in DataBase
// @Param task_id path string true "task_id"
// @Success 200 {json} json "status": "status state"
// @Failure 404 {string} string "Failed to get status"
// @Router /status/{task_id} [get]
func (s *Server) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	s.checkUnauthorization(w, r)

	url := strings.Split(r.URL.Path, "/")
	id := url[2]

	val, err := s.storage.Get(id)

	if err != nil {
		http.Error(w, "Failed to get status", http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": val.Readiness})
}

// @Summary Get task result
// @Description Post task_id in DataBase
// @Param task_id path string true "task_id"
// @Success 200 {json} json "result": "result state"
// @Failure 404 {string} string "Failed to get result"
// @Router /result/{task_id} [get]
func (s *Server) getResultHandler(w http.ResponseWriter, r *http.Request) {
	s.checkUnauthorization(w, r)

	url := strings.Split(r.URL.Path, "/")
	id := url[2]

	val, err := s.storage.Get(id)

	if err != nil {
		http.Error(w, "Failed to get result", http.StatusNotFound)

		return
	}

	if val.Result == "" {
		http.Error(w, "Failed to get result", http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"result": val.Result})
	}
}

func (s *Server) postRegHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.NewUUID()

	var req regRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := s.storage.Post_user(model.User{ID: id.String(), Login: req.Login, Password: req.Password}); err != nil {
		http.Error(w, "Failed to store value", http.StatusNotFound)

		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) postLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req regRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	user, err := s.storage.Get_user(model.User{ID: "", Login: req.Login, Password: req.Password})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	str_token := user.ID + user.Login + user.Password

	hash := sha256.Sum256([]byte(str_token))

	token := hex.EncodeToString(hash[:])

	if err := s.storage.Post_session(model.Session{User_id: user.ID, Session_id: token}); err != nil {
		http.Error(w, "Failed to store value", http.StatusNotFound)

		return
	}

	w.Header().Set("Authorization", "Bearer "+token)

	w.WriteHeader(http.StatusOK)
}

func CreateNewServer(storage Storage, addr string) error {
	server := newServer(storage)

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/", func(r chi.Router) {
		r.Post("/register, body: json{username,password}", server.postRegHandler)
		r.Post("/login, body: json{username,password}", server.postLoginHandler)
		r.Post("/task", server.postHandler)
		r.Get("/status/{task_id}", server.getStatusHandler)
		r.Get("/result/{task_id}", server.getResultHandler)
	})

	http_server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	return http_server.ListenAndServe()
}
