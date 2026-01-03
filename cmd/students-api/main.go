package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/student-ankitpandit/rest-api/http/handlers/student"
	"github.com/student-ankitpandit/rest-api/internal/config"
	"github.com/student-ankitpandit/rest-api/internal/storage/sqlite"
)

func main() {
	//load config
	cfg := config.MustLoad()

	//db setup
	storage, err := sqlite.New(cfg)
	if(err != nil) {
		log.Fatal(err)
	}

	slog.Info("storgae initialized", slog.String("current env is ", cfg.Env), slog.String("version", "1.0.0"))
	//setup router
	router := http.NewServeMux() //this is basically a router

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetStudentsById(storage))
	router.HandleFunc("GET /api/students", student.GetListm(storage))
	//server server

	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

	fmt.Printf("Server started %s", cfg.HTTPServer.Addr)


	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	//new go routine
	go func () {
		err := server.ListenAndServe()
		if(err != nil) {
			log.Fatal("Failed to start server")
		}
	}()

	<-done

	//logic to close/kill the server
	slog.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	// err := server.Shutdown(ctx)
	// if(err != nil) {
	// 	slog.Info("Failed to shutdown server", slog.String("error", err.Error()))
	// }
	if err := server.Shutdown(ctx); (err != nil) {
		slog.Info("Failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("Server shutdowned safely")
}

