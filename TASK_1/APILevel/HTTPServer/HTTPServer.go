package HTTPServer

import (
	"Kaspersky_Go/APILevel/HTTPServer/Templates"
	"Kaspersky_Go/ModeLevel/Structures"
	IAdapters2 "Kaspersky_Go/ServiceLevel/Interfaces/IAdapters"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type HTTPServer struct {
	server *http.Server
	queue  IAdapters2.Queue
	state  IAdapters2.StateStore
}

func NewHTTPServer(addr string, queue IAdapters2.Queue, state IAdapters2.StateStore) *HTTPServer {
	mux := http.NewServeMux()

	s := &HTTPServer{
		server: &http.Server{Addr: addr, Handler: mux},
		queue:  queue,
		state:  state,
	}

	mux.HandleFunc("/healthz", s.handleHealth)
	mux.HandleFunc("/enqueue", s.handleEnqueue)

	fmt.Printf("HTTP server initialized on %s\n", addr)
	return s
}

func (s *HTTPServer) Start() error {
	fmt.Println("Starting HTTP server…")
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	fmt.Println("Stopping HTTP server…")
	return s.server.Shutdown(ctx)
}

func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Health check called")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func (s *HTTPServer) handleEnqueue(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Incoming %s request to /enqueue\n", r.Method)

	if r.Method != "POST" {
		fmt.Println("Method not allowed")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req Templates.EnqueueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("Failed to decode request: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Decoded request: ID=%s, MaxRetries=%d\n", req.ID, req.MaxRetries)

	job := Structures.Job{
		ID:         req.ID,
		Payload:    req.Payload,
		MaxRetries: req.MaxRetries,
	}
	if err := s.queue.Enqueue(job); err != nil {
		fmt.Printf("Failed to enqueue job %s: %v\n", job.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Printf("Job %s enqueued successfully\n", job.ID)

	s.state.Set(job.ID, Structures.JobStatus{State: Structures.StateQueued})
	fmt.Printf("State for job %s set to queued\n", job.ID)

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "queued"})
	fmt.Printf("Response sent for job %s\n", job.ID)
}
