package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	_ "http_server/docs"

	chi "github.com/go-chi/chi/v5"
	uuid "github.com/google/uuid"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Storage interface {
	Get(key string) (*map[string]string, error)
	Post(key string, value map[string]string) error
	Put(key string, value map[string]string) error
}

type Server struct {
	storage Storage
}

func newServer(curr_storage Storage) *Server {
	return &Server{storage: curr_storage}
}

func task_work(s *Server, id uuid.UUID) {
	time.Sleep(10 * time.Second)

	s.storage.Put(id.String(), map[string]string{"status": "ready", "result": "some rubish"})
}

// @Summary Post task task_id
// @Description Post task_id in DataBase and start task
// @Success 201 {json} json "task_id": "id value"
// @Failure 404 {string} string "Failed to store value"
// @Router /task [post]
func (s *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.NewUUID()

	if err := s.storage.Post(id.String(), map[string]string{"status": "in_progress", "result": ""}); err != nil {
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
	url := strings.Split(r.URL.Path, "/")
	id := url[2]

	val, err := s.storage.Get(id)

	if err != nil {
		http.Error(w, "Failed to get status", http.StatusNotFound)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": (*val)["status"]})
}

// @Summary Get task result
// @Description Post task_id in DataBase
// @Param task_id path string true "task_id"
// @Success 200 {json} json "result": "result state"
// @Failure 404 {string} string "Failed to get result"
// @Router /result/{task_id} [get]
func (s *Server) getResultHandler(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	id := url[2]

	val, err := s.storage.Get(id)

	if err != nil {
		http.Error(w, "Failed to get result", http.StatusNotFound)

		return
	}

	if (*val)["result"] == "" {
		http.Error(w, "Failed to get result", http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"result": (*val)["result"]})
	}
}

func CreateNewServer(storage Storage, addr string) error {
	server := newServer(storage)

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/", func(r chi.Router) {
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
