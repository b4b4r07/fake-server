package main

import (
	"log"
	"net/http"
	"os"
)

type (
	Server struct {
		logger  *log.Logger
		mux     *http.ServeMux
		faker   *Faker
		common  service
		payment *paymentService
	}
	service struct {
		server *Server
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

func NewServer(options ...func(*Server)) *Server {
	s := &Server{
		logger: log.New(os.Stdout, "", 0),
		mux:    http.NewServeMux(),
		faker:  &Faker{},
	}

	for _, f := range options {
		f(s)
	}

	s.common.server = s
	s.payment = (*paymentService)(&s.common)

	s.mux.HandleFunc("/list", s.payment.list)
	s.mux.HandleFunc("/save", s.payment.SaveMember)
	// s.mux.HandleFunc("/payment/SaveMember.idPass", s.paymentSaveMember)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, "-", r.RequestURI)
	s.mux.ServeHTTP(w, r)
}
