package routes

import (
	"main/internal/user"
	"main/internal/ws"
	"main/middlewares"

	"github.com/go-chi/chi/v5"
)

func Use(router *chi.Mux, user *user.Handler, ws *ws.Handler) {
	router.Group(func(r chi.Router) {
		r.Post("/signup", user.Register)
		r.Post("/signin", user.Login)
	})

	router.Group(func(r chi.Router) {
		r.Use(middlewares.Authenticator)
		r.Post("/signout", user.Logout)
		r.Post("/createRoom", ws.CreateRoom)
		r.Get("/getRooms", ws.GetRooms)
		r.Get("/joinRoom/{roomId}", ws.JoinRoom)
		r.Get("/getSubscriber/{roomId}", ws.GetSubscriber)
	})
}
