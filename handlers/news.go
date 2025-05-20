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
	"myproject/utils"
)

func ViewNewsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID новости", http.StatusBadRequest)
		return
	}

	news, err := database.GetNewsByID(id)
	if err != nil {
		http.Error(w, "Новость не найдена", http.StatusNotFound)
		return
	}

	comments, err := database.GetCommentsByNewsID(id)
	if err != nil {
		log.Printf("Ошибка при получении комментариев: %v", err)
	}

	data := map[string]interface{}{
		"News":     news,
		"Comments": comments,
	}

	utils.RenderTemplate(w, r, "view_news.html", data)
}

func AddNewsHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method == "POST" {
		title := r.FormValue("title")
		content := r.FormValue("content")

		userID := middleware.GetCurrentUserID(r)
		username := middleware.GetCurrentUsername(r)

		news := models.NewNews(title, content, userID, username)
		_, err := database.CreateNews(news)
		if err != nil {
			log.Printf("Ошибка при создании новости: %v", err)
			http.Error(w, "Ошибка при добавлении новости", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	utils.RenderTemplate(w, r, "add_news.html", nil)
}

func EditNewsHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID новости", http.StatusBadRequest)
		return
	}

	news, err := database.GetNewsByID(id)
	if err != nil {
		http.Error(w, "Новость не найдена", http.StatusNotFound)
		return
	}

	currentUserID := middleware.GetCurrentUserID(r)
	isAdmin := false
	session, _ := middleware.Store.Get(r, "session-name")
	if val, ok := session.Values["is_admin"].(bool); ok {
		isAdmin = val
	}

	if news.AuthorID != currentUserID && !isAdmin {
		http.Error(w, "У вас нет прав для редактирования этой новости", http.StatusForbidden)
		return
	}

	if r.Method == "POST" {
		title := r.FormValue("title")
		content := r.FormValue("content")

		news.Title = title
		news.Content = content
		news.UpdatedAt = time.Now()

		err = database.UpdateNews(news)
		if err != nil {
			log.Printf("Ошибка при обновлении новости: %v", err)
			http.Error(w, "Ошибка при редактировании новости", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/news/"+vars["id"], http.StatusSeeOther)
		return
	}

	data := map[string]interface{}{
		"News": news,
	}

	utils.RenderTemplate(w, r, "edit_news.html", data)
}

func DeleteNewsHandler(w http.ResponseWriter, r *http.Request) {
	if !middleware.IsAuthenticated(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID новости", http.StatusBadRequest)
		return
	}

	news, err := database.GetNewsByID(id)
	if err != nil {
		http.Error(w, "Новость не найдена", http.StatusNotFound)
		return
	}

	currentUserID := middleware.GetCurrentUserID(r)
	isAdmin := false
	session, _ := middleware.Store.Get(r, "session-name")
	if val, ok := session.Values["is_admin"].(bool); ok {
		isAdmin = val
	}

	if news.AuthorID != currentUserID && !isAdmin {
		http.Error(w, "У вас нет прав для удаления этой новости", http.StatusForbidden)
		return
	}

	err = database.DeleteNews(id)
	if err != nil {
		log.Printf("Ошибка при удалении новости: %v", err)
		http.Error(w, "Ошибка при удалении новости", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
} 