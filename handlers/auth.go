package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"myproject/database"
	"myproject/middleware"
	"myproject/models"
	"myproject/utils"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")
		email := r.FormValue("email")

		_, err := database.GetUserByUsername(username)
		if err == nil {
			http.Error(w, "Пользователь с таким именем уже существует", http.StatusBadRequest)
			return
		} else if err != sql.ErrNoRows {
			log.Printf("Ошибка при проверке пользователя: %v", err)
			http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Ошибка при хешировании пароля: %v", err)
			http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
			return
		}

		user := models.NewUser(username, string(hashedPassword), email)
		_, err = database.CreateUser(user)
		if err != nil {
			log.Printf("Ошибка при создании пользователя: %v", err)
			http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	utils.RenderTemplate(w, r, "register.html", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := database.GetUserByUsername(username)
		if err != nil {
			http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
			return
		}

		if user.IsBlocked {
			http.Error(w, "Ваш аккаунт заблокирован", http.StatusForbidden)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
			return
		}

		session, _ := middleware.Store.Get(r, "session-name")
		session.Values["authenticated"] = true
		session.Values["username"] = user.Username
		session.Values["user_id"] = user.ID
		session.Values["is_admin"] = (user.Username == "admin") // Простая проверка на админа
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	utils.RenderTemplate(w, r, "login.html", nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Values["user_id"] = 0
	session.Values["is_admin"] = false
	session.Options.MaxAge = -1
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
} 