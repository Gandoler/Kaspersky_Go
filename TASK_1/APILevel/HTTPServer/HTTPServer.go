package HttpR

import (
	"Kaspersky_Go/APILevel/HTTPServer/Templates"
	"Kaspersky_Go/APILevel/IAdapters"
	"Kaspersky_Go/ModeLevel/Structures"
	"context"
	"encoding/json"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
	queue  IAdapters.Queue
	state  IAdapters.StateStore
}

func NewHTTPServer(addr string, queue IAdapters.Queue, state IAdapters.StateStore) *HTTPServer {
	mux := http.NewServeMux()

	s := &HTTPServer{
		server: &http.Server{Addr: addr, Handler: mux},
		queue:  queue,
		state:  state,
	}

	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/enqueue", s.handleEnqueue)

	return s
}

func (s *HTTPServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func (s *HTTPServer) handleEnqueue(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req Templates.EnqueueRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	job := Structures.Job{ID: req.ID, Payload: req.Payload, MaxRetries: req.MaxRetries}
	if err := s.queue.Enqueue(job); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.state.Set(job.ID, Structures.JobStatus{State: Structures.StateQueued})
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "queued"})

}
