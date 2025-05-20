package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"myproject/database"
	"myproject/middleware"
	"myproject/models"
)

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Error(w, "Для добавления комментария необходимо войти в систему", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	newsID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID новости", http.StatusBadRequest)
		return
	}

	content := r.FormValue("content")

	userID := middleware.GetCurrentUserID(r)
	username := middleware.GetCurrentUsername(r)

	comment := models.NewComment(newsID, userID, username, content)
	_, err = database.CreateComment(comment)
	if err != nil {
		log.Printf("Ошибка при создании комментария: %v", err)
		http.Error(w, "Ошибка при добавлении комментария", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/news/"+vars["id"], http.StatusSeeOther)
}

func EditCommentHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Error(w, "Для редактирования комментария необходимо войти в систему", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID комментария", http.StatusBadRequest)
		return
	}

	comment, err := database.GetCommentByID(commentID)
	if err != nil {
		http.Error(w, "Комментарий не найден", http.StatusNotFound)
		return
	}

	currentUserID := middleware.GetCurrentUserID(r)
	isAdmin := false
	session, _ := middleware.Store.Get(r, "session-name")
	if val, ok := session.Values["is_admin"].(bool); ok {
		isAdmin = val
	}

	if comment.UserID != currentUserID && !isAdmin {
		http.Error(w, "У вас нет прав для редактирования этого комментария", http.StatusForbidden)
		return
	}

	content := r.FormValue("content")

	comment.Content = content
	comment.UpdatedAt = time.Now()

	err = database.UpdateComment(comment)
	if err != nil {
		log.Printf("Ошибка при обновлении комментария: %v", err)
		http.Error(w, "Ошибка при редактировании комментария", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/news/"+strconv.Itoa(comment.NewsID), http.StatusSeeOther)
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Error(w, "Для удаления комментария необходимо войти в систему", http.StatusUnauthorized)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	commentID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID комментария", http.StatusBadRequest)
		return
	}

	comment, err := database.GetCommentByID(commentID)
	if err != nil {
		http.Error(w, "Комментарий не найден", http.StatusNotFound)
		return
	}

	currentUserID := middleware.GetCurrentUserID(r)
	isAdmin := false
	session, _ := middleware.Store.Get(r, "session-name")
	if val, ok := session.Values["is_admin"].(bool); ok {
		isAdmin = val
	}

	if comment.UserID != currentUserID && !isAdmin {
		http.Error(w, "У вас нет прав для удаления этого комментария", http.StatusForbidden)
		return
	}

	newsID := comment.NewsID

	err = database.DeleteComment(commentID)
	if err != nil {
		log.Printf("Ошибка при удалении комментария: %v", err)
		http.Error(w, "Ошибка при удалении комментария", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/news/"+strconv.Itoa(newsID), http.StatusSeeOther)
} 