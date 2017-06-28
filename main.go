package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type (
	Server struct {
		logger *log.Logger
		mux    *http.ServeMux
		faker  *Faker
	}
	Faker struct {
		Members Members
	}
	Members []Member
	Member  struct {
		ID   string
		Name string
	}
)

var (
	lock sync.RWMutex
)

func NewServer(options ...func(*Server)) *Server {
	s := &Server{
		logger: log.New(os.Stdout, "", 0),
		mux:    http.NewServeMux(),
		faker:  new(Faker),
	}

	for _, f := range options {
		f(s)
	}

	s.mux.HandleFunc("/list", s.list)
	s.mux.HandleFunc("/add", s.add)

	return s
}

func (ms *Members) get(id string) (Member, error) {
	for _, m := range *ms {
		if m.ID == id {
			return m, nil
		}
	}
	return Member{}, errors.New("not found")
}

func (ms *Members) add(m Member) error {
	_, err := ms.get(m.ID)
	if err == nil {
		return errors.New("already exists")
	}
	*ms = append(*ms, m)
	return nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, "-", r.RequestURI)
	s.mux.ServeHTTP(w, r)
}

func (s *Server) add(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("MemberID")
	if id == "" {
		http.Error(w, "Empty MemberID", http.StatusBadRequest)
		return
	}
	name := r.FormValue("MemberName")
	if name == "" {
		http.Error(w, "Empty MemberName", http.StatusBadRequest)
		return
	}
	lock.Lock()
	defer lock.Unlock()
	err := s.faker.Members.add(Member{
		ID:   id,
		Name: name,
	})
	if err != nil {
		http.Error(w, "MemberID already exists", http.StatusBadRequest)
		return
	}
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	defer lock.RUnlock()
	for _, member := range s.faker.Members {
		fmt.Fprintf(w, "%#v\n", member)
	}
}

func main() {
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	logger := log.New(os.Stdout, "", 0)
	s := NewServer(func(s *Server) { s.logger = logger })
	addr := ":8080"
	h := &http.Server{Addr: addr, Handler: s}

	go func() {
		logger.Printf("Listening on http://0.0.0.0%s\n", addr)
		if err := h.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}()
	<-stop

	logger.Println("\nShutting down the server...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	h.Shutdown(ctx)
	logger.Println("Server gracefully stopped")
}
