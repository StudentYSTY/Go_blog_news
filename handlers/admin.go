package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"myproject/database"
	"myproject/middleware"
	"myproject/utils"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	isAdmin, adminOk := session.Values["is_admin"].(bool)

	if !ok || !auth || !adminOk || !isAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	users, err := database.GetAllUsers()
	if err != nil {
		log.Printf("Ошибка при получении пользователей: %v", err)
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Users": users,
	}

	utils.RenderTemplate(w, r, "admin.html", data)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	isAdmin, adminOk := session.Values["is_admin"].(bool)

	if !ok || !auth || !adminOk || !isAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	currentUserID := middleware.GetCurrentUserID(r)
	if user.ID == currentUserID {
		http.Error(w, "Вы не можете удалить свою учетную запись", http.StatusForbidden)
		return
	}

	err = database.DeleteUser(userID)
	if err != nil {
		log.Printf("Ошибка при удалении пользователя: %v", err)
		http.Error(w, "Ошибка при удалении пользователя", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func BlockUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	isAdmin, adminOk := session.Values["is_admin"].(bool)

	if !ok || !auth || !adminOk || !isAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	currentUserID := middleware.GetCurrentUserID(r)
	if user.ID == currentUserID {
		http.Error(w, "Вы не можете заблокировать свою учетную запись", http.StatusForbidden)
		return
	}

	err = database.BlockUser(userID)
	if err != nil {
		log.Printf("Ошибка при блокировке пользователя: %v", err)
		http.Error(w, "Ошибка при блокировке пользователя", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func UnblockUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := middleware.Store.Get(r, "session-name")
	auth, ok := session.Values["authenticated"].(bool)
	isAdmin, adminOk := session.Values["is_admin"].(bool)

	if !ok || !auth || !adminOk || !isAdmin {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID пользователя", http.StatusBadRequest)
		return
	}

	err = database.UnblockUser(userID)
	if err != nil {
		log.Printf("Ошибка при разблокировке пользователя: %v", err)
		http.Error(w, "Ошибка при разблокировке пользователя", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
} 