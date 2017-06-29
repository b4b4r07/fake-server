package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	port = flag.Int("port", 8000, "port number")
)

func main() {
	flag.Parse()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)

	logger := log.New(os.Stdout, "", 0)
	s := NewServer(func(s *Server) { s.logger = logger })
	addr := fmt.Sprintf(":%d", *port)
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
