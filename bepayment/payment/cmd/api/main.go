package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chonlatee/payment/internal/routes"
)

func main() {

	r := routes.Route()
	port := "8080"

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}

	}()

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	log.Println("Shuting down server...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("server shutdown err: ", err)
	}

	log.Println("Server exiting")
}
