package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"myproject/config"
	"myproject/database"
	"myproject/handlers"
	"myproject/middleware"
	"myproject/utils"
)

func main() {
	config.LoadConfig()

	err := database.InitDB(config.AppConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer database.CloseDB()

	utils.LoadTemplates()

	middleware.Store.MaxAge(86400 * 7) // 1 неделя

	r := mux.NewRouter()

	r.Use(middleware.LoggerMiddleware)

	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")

	r.HandleFunc("/register", handlers.RegisterHandler).Methods("GET", "POST")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET")

	r.HandleFunc("/news/{id:[0-9]+}", handlers.ViewNewsHandler).Methods("GET")
	r.HandleFunc("/add-news", handlers.AddNewsHandler).Methods("GET", "POST")
	r.HandleFunc("/edit-news/{id:[0-9]+}", handlers.EditNewsHandler).Methods("GET", "POST")
	r.HandleFunc("/delete-news/{id:[0-9]+}", handlers.DeleteNewsHandler).Methods("POST")

	r.HandleFunc("/news/{id:[0-9]+}/comment", handlers.AddCommentHandler).Methods("POST")
	r.HandleFunc("/comment/{id:[0-9]+}/edit", handlers.EditCommentHandler).Methods("POST")
	r.HandleFunc("/comment/{id:[0-9]+}/delete", handlers.DeleteCommentHandler).Methods("POST")

	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.AdminRequired)
	adminRouter.HandleFunc("", handlers.AdminHandler).Methods("GET")
	adminRouter.HandleFunc("/delete-user/{id:[0-9]+}", handlers.DeleteUserHandler).Methods("POST")
	adminRouter.HandleFunc("/block-user/{id:[0-9]+}", handlers.BlockUserHandler).Methods("POST")
	adminRouter.HandleFunc("/unblock-user/{id:[0-9]+}", handlers.UnblockUserHandler).Methods("POST")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Сервер запущен на порту %s", config.AppConfig.ServerPort)
	log.Fatal(http.ListenAndServe(":"+config.AppConfig.ServerPort, r))
} 