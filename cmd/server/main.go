package main

import (
	"log"
	"net/http"

	"github.com/boldnguyen/friend-management/internal/handler"
	"github.com/boldnguyen/friend-management/internal/pkg/db"
	"github.com/boldnguyen/friend-management/internal/repository"
	"github.com/boldnguyen/friend-management/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// Open DB connection
	db, err := db.ConnectDB("postgres://friend-management:1234@db:5432/friend-management?sslmode=disable")
	if err != nil {
		log.Fatal("Faild to connect to the database", err)
	}
	defer db.Close()

	// Initialize friend repository and service with the database connection
	friendRepository := repository.NewFriendRepository(db)
	friendService := service.NewFriendService(friendRepository)

	// Init Router
	r := initRouter(friendService)

	// Start server
	log.Println("Starting app at port: 5000")
	if err := http.ListenAndServe(":5000", r); err != nil {
		log.Println("Server error", err)
	}

}

// The initRouter function is used to set up and configure the routes for my web application using the Chi router.
func initRouter(friendService service.FriendService) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.Post("/friend/create", handler.NewHandler(friendService))
	r.Post("/friend/list", handler.FriendListHandler(friendService))

	return r

}
