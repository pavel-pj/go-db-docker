package router

import (
	"db200/handlers"
	repository "db200/repositories"
	"net/http"
)

type Router struct {
	userHandler *handlers.UserHandler
}

func NewRouter(db *repository.UserRepository) *Router {
	return &Router{
		userHandler: handlers.NewUserHandler(db),
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := r.URL.Path

	switch {
	case path == "/api/users" && r.Method == http.MethodPost:
		rt.userHandler.CreateUser(w, r)
	/*case path == "/api/users" && r.Method == http.MethodGet:
	      rt.userHandler.GetAllUsers(w, r)
	  case strings.HasPrefix(path, "/api/users/") && r.Method == http.MethodGet:
	      rt.userHandler.GetUserByID(w, r)
	*/
	case path == "/health":
		rt.healthCheck(w, r)
	case path == "/":
		rt.homeHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (rt *Router) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Welcome to Users API"}`))
}

func (rt *Router) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}
