package middleware

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

var (
	Store = sessions.NewCookieStore([]byte("secret-key-123"))
)

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session-name")
		auth, ok := session.Values["authenticated"].(bool)

		if !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func AdminRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := Store.Get(r, "session-name")
		auth, ok := session.Values["authenticated"].(bool)
		isAdmin, adminOk := session.Values["is_admin"].(bool)

		if !ok || !auth || !adminOk || !isAdmin {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func IsAuthenticated(r *http.Request) bool {
	session, err := Store.Get(r, "session-name")
	if err != nil {
		log.Println("Ошибка получения сессии:", err)
		return false
	}
	
	auth, ok := session.Values["authenticated"].(bool)
	return ok && auth
}

func GetCurrentUsername(r *http.Request) string {
	session, err := Store.Get(r, "session-name")
	if err != nil {
		log.Println("Ошибка получения сессии:", err)
		return ""
	}
	
	username, ok := session.Values["username"].(string)
	if !ok {
		return ""
	}
	return username
}

func GetCurrentUserID(r *http.Request) int {
	session, err := Store.Get(r, "session-name")
	if err != nil {
		log.Println("Ошибка получения сессии:", err)
		return 0
	}
	
	userID, ok := session.Values["user_id"].(int)
	if !ok {
		return 0
	}
	return userID
} 