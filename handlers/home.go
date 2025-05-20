package handlers

import (
	"log"
	"net/http"

	"myproject/database"
	"myproject/models"
	"myproject/utils"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	news, err := database.GetAllNews()
	if err != nil {
		log.Printf("Ошибка при получении новостей: %v", err)
		http.Error(w, "Ошибка при получении новостей", http.StatusInternalServerError)
		return
	}

	type NewsWithComments struct {
		models.News
		CommentCount int
	}

	var newsWithComments []NewsWithComments
	for _, n := range news {
		comments, err := database.GetCommentsByNewsID(n.ID)
		commentCount := 0
		if err == nil {
			commentCount = len(comments)
		}

		newsWithComments = append(newsWithComments, NewsWithComments{
			News:         n,
			CommentCount: commentCount,
		})
	}

	data := map[string]interface{}{
		"News": newsWithComments,
	}

	utils.RenderTemplate(w, r, "index.html", data)
} 