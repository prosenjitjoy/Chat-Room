package main

import (
	"context"
	"log"
	"main/database"
	"main/internal/user"
	"main/internal/ws"
	"main/middlewares"
	"main/routes"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx := context.Background()
	conn, err := database.NewDatabase(ctx)
	if err != nil {
		log.Fatal("Could not initialize database connection")
	}
	defer conn.Close(ctx)

	userAdapter := user.NewAdapter(conn.GetDB())
	userService := user.NewService(userAdapter)
	userHandler := user.NewHandler(userService)

	chatServer := ws.NewChatServer()
	go chatServer.Run()
	wsHandler := ws.NewHandler(chatServer)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middlewares.Cors)
	routes.Use(router, userHandler, wsHandler)

	errc := make(chan error, 1)
	go func() {
		errc <- http.ListenAndServe(":5000", router)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	select {
	case err := <-errc:
		log.Printf("failed to serve: %v", err)
	case sig := <-sigs:
		log.Printf("terminating: %v", sig)
	}
}
