package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ismailtsdln/DevTestrider/internal/config"
	"github.com/ismailtsdln/DevTestrider/internal/engine"
)

type Server struct {
	Router     *chi.Mux
	Config     config.ServerConfig
	clients    map[chan string]bool
	mu         sync.Mutex
	LastResult *engine.TestResult
}

func NewServer(cfg config.ServerConfig) *Server {
	s := &Server{
		Router:  chi.NewRouter(),
		Config:  cfg,
		clients: make(map[chan string]bool),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)

	s.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// API Routes
	s.Router.Route("/api", func(r chi.Router) {
		r.Get("/events", s.handleEvents)
		r.Get("/results/latest", s.handleLatestResult)
	})

	// Serve Static Files (Frontend)
	// In development, we might just rely on Vite dev server,
	// but for the final binary we'd serve 'web/dist'.
	// For now, let's assume we serve from ./web/dist
	fs := http.FileServer(http.Dir("./web/dist"))
	s.Router.Handle("/*", http.StripPrefix("/", fs))
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.Config.Port)
	if s.Config.Port == 0 {
		addr = ":8080"
	}
	fmt.Printf("Starting server on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, s.Router)
}

func (s *Server) handleLatestResult(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if s.LastResult == nil {
		w.Write([]byte("null"))
		return
	}
	json.NewEncoder(w).Encode(s.LastResult)
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	messageChan := make(chan string)
	s.mu.Lock()
	s.clients[messageChan] = true
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.clients, messageChan)
		s.mu.Unlock()
		close(messageChan)
	}()

	notify := r.Context().Done()

	for {
		select {
		case <-notify:
			return
		case msg := <-messageChan:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-time.After(15 * time.Second):
			// Keep-alive
			fmt.Fprintf(w, ": keepalive\n\n")
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		}
	}
}

func (s *Server) Broadcast(result *engine.TestResult) {
	s.mu.Lock()
	s.LastResult = result

	data, err := json.Marshal(result)
	if err != nil {
		s.mu.Unlock()
		return
	}
	msg := string(data)

	for client := range s.clients {
		select {
		case client <- msg:
		default:
			// Client blocked, skip
		}
	}
	s.mu.Unlock()
}
