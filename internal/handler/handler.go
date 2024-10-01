package handler

import (
	"YandexLearnMiddle/internal/service"
	"github.com/gorilla/mux"
)

// Server представляет основной сервер приложения с настройками роутов.
type Server struct {
	Router *mux.Router
}

// NewServer создает новый сервер с заданными роутами.
func NewServer() *Server {
	s := &Server{
		Router: mux.NewRouter(),
	}
	s.routes() // Настраиваем маршруты
	return s
}

// routes настраивает маршруты сервера.
func (s *Server) routes() {
	linkService := service.NewLinkService()    // Инициализируем сервис ссылок
	linkHandler := NewLinkHandler(linkService) // Инициализируем хендлер для работы с ссылками

	// Публичные маршруты
	s.Router.HandleFunc("/register", SignUp).Methods("POST")
	s.Router.HandleFunc("/auth", SignIn).Methods("POST")

	// Группа защищенных маршрутов с MiddlewareJWT
	protected := s.Router.PathPrefix("/api").Subrouter()
	protected.Use(MiddlewareJWT) // Подключаем Middleware для защиты маршрутов
	protected.HandleFunc("/shorten", linkHandler.ShortenLink).Methods("POST")
	protected.HandleFunc("/original", linkHandler.GetOriginalLink).Methods("GET")
	protected.HandleFunc("/links", linkHandler.GetAllLinks).Methods("GET")
	protected.HandleFunc("/user/urls", linkHandler.DeleteLinks).Methods("DELETE")
}
