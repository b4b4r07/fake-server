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
		common  handler
		payment *paymentHandler
	}
	handler struct {
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
	Error struct {
		Code string
		Info string
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
	s.payment = (*paymentHandler)(&s.common)

	s.mux.HandleFunc("/list", s.payment.list)
	s.mux.HandleFunc("/payment/SaveMember.idPass", s.payment.SaveMember)
	s.mux.HandleFunc("/payment/UpdateMember.idPass", s.payment.UpdateMember)
	s.mux.HandleFunc("/payment/DeleteMember.idPass", s.payment.DeleteMember)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, "-", r.RequestURI)
	s.mux.ServeHTTP(w, r)
}
