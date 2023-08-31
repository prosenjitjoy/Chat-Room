package middlewares

import (
	"net/http"

	"github.com/go-chi/cors"
)

func Cors(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Content-Type"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 300,
	}).Handler(next)
}
