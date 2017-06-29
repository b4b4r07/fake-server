package main

import (
	"fmt"
	"net/http"
	"sync"
)

type (
	paymentService service
)

var (
	lock sync.RWMutex
)

func (p *paymentService) SaveMember(w http.ResponseWriter, r *http.Request) {
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
	err := p.server.faker.Members.add(Member{
		ID:   id,
		Name: name,
	})
	if err != nil {
		http.Error(w, "MemberID already exists", http.StatusBadRequest)
		return
	}
}

func (p *paymentService) list(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	defer lock.RUnlock()
	for _, member := range p.server.faker.Members {
		fmt.Fprintf(w, "%#v\n", member)
	}
}