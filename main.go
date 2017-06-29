package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

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
