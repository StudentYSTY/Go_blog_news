package utils

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"myproject/middleware"
)

var (
	Templates *template.Template
)

func LoadTemplates() {
	var err error
	
	funcMap := template.FuncMap{
		"formatDate": FormatDate,
		"isAuthor": IsAuthor,
	}
	
	Templates = template.New("").Funcs(funcMap)
	Templates, err = Templates.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Ошибка загрузки шаблонов: %v", err)
	}
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	
	data["IsAuthenticated"] = middleware.IsAuthenticated(r)
	data["CurrentUsername"] = middleware.GetCurrentUsername(r)
	data["CurrentUserID"] = middleware.GetCurrentUserID(r)
	
	err := Templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Printf("Ошибка при отображении шаблона %s: %v", tmpl, err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}

func FormatDate(t time.Time) string {
	return t.Format("02.01.2006 15:04")
}

func IsAuthor(currentUserID, authorID int) bool {
	return currentUserID == authorID
} 