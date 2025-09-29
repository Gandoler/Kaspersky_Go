package HandlersTESTS

import (
	"Kaspersky_Go/APILevel/Adapters"
	"Kaspersky_Go/APILevel/HTTPServer"
	"Kaspersky_Go/ModeLevel/Structures"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHealth(t *testing.T) {
	queue := Adapters.NewMemoryQueue(1)
	state := Adapters.NewMemoryStateStore()
	srv := HTTPServer.NewHTTPServer(":0", queue, state)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	srv.HandleHealth(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestHandleEnqueue(t *testing.T) {
	queue := Adapters.NewMemoryQueue(1)
	state := Adapters.NewMemoryStateStore()
	srv := HTTPServer.NewHTTPServer(":0", queue, state)

	body, _ := json.Marshal(map[string]interface{}{
		"id": "job1", "payload": "test", "max_retries": 3,
	})
	req := httptest.NewRequest(http.MethodPost, "/enqueue", bytes.NewReader(body))
	w := httptest.NewRecorder()
	srv.HandleEnqueue(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusAccepted)
	}

	// Проверяем, что статус в state установлен
	st, ok := state.Get("job1")
	if !ok {
		t.Fatal("expected job status in state store")
	}
	if st.State != Structures.StateQueued {
		t.Errorf("got %v, want %v", st.State, Structures.StateQueued)
	}
}

func TestHandleEnqueue_MethodNotAllowed(t *testing.T) {
	queue := Adapters.NewMemoryQueue(1)
	state := Adapters.NewMemoryStateStore()
	srv := HTTPServer.NewHTTPServer(":0", queue, state)

	req := httptest.NewRequest(http.MethodGet, "/enqueue", nil)
	w := httptest.NewRecorder()
	srv.HandleEnqueue(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestHandleEnqueue_BadJSON(t *testing.T) {
	queue := Adapters.NewMemoryQueue(1)
	state := Adapters.NewMemoryStateStore()
	srv := HTTPServer.NewHTTPServer(":0", queue, state)

	badBody := bytes.NewBufferString("{bad json")
	req := httptest.NewRequest(http.MethodPost, "/enqueue", badBody)
	w := httptest.NewRecorder()
	srv.HandleEnqueue(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}
