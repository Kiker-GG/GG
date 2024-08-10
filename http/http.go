package http

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	time.Sleep(10000 * time.Second)

	s.storage.Put(id.String(), map[string]string{"status": "ready", "result": "some rubish"})
}

// @Summary Post task task_id
// @Description Post task_id in DataBase and start task
// @Success 200 {string} string "Id"
// @Failure 500 {string} string "Failed to store value"
// @Router /task [post]
func (s *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.NewUUID()

	if err := s.storage.Post(id.String(), map[string]string{"status": "in_progress", "result": ""}); err != nil {
		http.Error(w, "Failed to store value", http.StatusInternalServerError)

		return
	}

	go task_work(s, id)

	fmt.Fprint(w, id.String())
}

// @Summary Get task status
// @Description Post task_id in DataBase
// @Param task_id query string true "task_id"
// @Success 200 {json} json "status": "status state"
// @Failure 500 {string} string "Failed to get status"
// @Router /status/ [get]
func (s *Server) getStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("task_id")
	val, err := s.storage.Get(id)

	if err != nil {
		http.Error(w, "Failed to get status", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": (*val)["status"]})

	//w.WriteHeader(http.StatusOK)
}

// @Summary Get task result
// @Description Post task_id in DataBase
// @Param task_id query string true "task_id"
// @Success 200 {json} json "result": "result state"
// @Failure 500 {string} string "Failed to get result"
// @Router /result/ [get]
func (s *Server) getResultHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("task_id")
	val, err := s.storage.Get(id)

	if err != nil {
		http.Error(w, "Failed to get result", http.StatusInternalServerError)

		return
	}

	if (*val)["result"] == "" {
		fmt.Fprint(w, "Task hasn't finished yet!")
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"result": (*val)["result"]})
	}

	//w.WriteHeader(http.StatusOK)
}

func CreateNewServer(storage Storage, addr string) error {
	server := newServer(storage)

	r := chi.NewRouter()

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/", func(r chi.Router) {
		r.Post("/task", server.postHandler)
		r.Get("/status/", server.getStatusHandler)
		r.Get("/result/", server.getResultHandler)
	})

	http_server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	return http_server.ListenAndServe()
}
