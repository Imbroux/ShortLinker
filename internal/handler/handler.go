package handler

import (
	"YandexLearnMiddle/internal/service"
	"github.com/gorilla/mux"
	"net/http"
)

func InitRouters() http.Handler {
	r := mux.NewRouter()
	linkService := service.NewLinkService()
	linkHandler := NewLinkHandler(linkService)

	r.HandleFunc("/register", SignUp).Methods("POST")
	r.HandleFunc("/auth", SignIn).Methods("POST")

	// Группа защищенных маршрутов с MiddlewareJWT
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(MiddlewareJWT)
	protected.HandleFunc("/shorten", linkHandler.ShortenLink).Methods("POST")
	protected.HandleFunc("/original", linkHandler.GetOriginalLink).Methods("GET")
	protected.HandleFunc("/links", linkHandler.GetAllLinks).Methods("GET")

	return r
}
