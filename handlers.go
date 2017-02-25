package vrmp

import (
	"net/http"
	"time"
	"strconv"
	"log"
)

// Index Handler
func handleIndex(w http.ResponseWriter, r *http.Request) {
	// get page number
	pageString := r.URL.Query().Get("page")
	log.Println("User requested page", pageString)
	if pageString == "" {
		pageString = "1"
	}

	// convert str to int
	pageNumber, _ := strconv.Atoi(pageString)

	// get all user data
	users := getAllUsers(getContext(r))
	images := getAllImages(getContext(r), pageNumber)	

	tmplData := map[string]interface{}{
		"time":      time.Now().Format(time.UnixDate),
		"username":  getUserEmail(r),
		"loginURL":  getLoginURL(r, "/"),
		"logoutURL": getLogoutURL(r, ""),
		"users": 	 users,
		"images":    images,
	}
	renderTemplate(w, "index.html", tmplData)
}
