package main

import (
	"context"
	//"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/callmeskyy111/golang-students-api/internal/config"
	"github.com/callmeskyy111/golang-students-api/internal/http/handlers/student"
	"github.com/callmeskyy111/golang-students-api/internal/storage/sqlite"
)

func main() {
	// load config
	cfg:= config.MustLoad()

	// db setup
	storage,err:=sqlite.New(cfg)
	if err!=nil{
		log.Fatal(err)
	}
	slog.Info("storage initialised",slog.String("env",cfg.Env), slog.String("version","1.0.0"))

	// setup router
	router:=http.NewServeMux()
	router.HandleFunc("POST /api/students",student.New(storage))
	router.HandleFunc("GET /api/students/{id}",student.GetById(storage))
	router.HandleFunc("GET /api/students",student.GetList(storage))
	


	// setup server
	server:=http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

	slog.Info("Server started..", slog.String("address",cfg.Addr))
	//fmt.Printf("Server started.. %s ✅", cfg.Addr)


	// Graceful Shutdown - goroutines
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func(){
	err:=server.ListenAndServe()
	if err!=nil{
		log.Fatal("Failed to start server 🔴")
	}
	}()

	<- done

	ctxt,cancel := context.WithTimeout(context.Background(),5 * time.Second)
	defer cancel()


	if err:=server.Shutdown(ctxt); err != nil{
		slog.Error("Failed to shutdown server", slog.String("error",err.Error()))
	}
	
	slog.Info("Server SHUTDOWN Successfully!")
}

// go run cmd/students-api/main.go -config config/local.yaml

// 02.04