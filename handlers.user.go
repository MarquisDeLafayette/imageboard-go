package vrmp

import (
	"fmt"
	"net/http"
	"time"
)

// Handle User Update
func handleUserUpdate(w http.ResponseWriter, r *http.Request) {
	// POST
	email := getUserEmail(r)
	if email == "" {
		redirToLogin(w, r)
		return
	}
	name := r.FormValue("name")
	bio := r.FormValue("bio")
	updateUser(getContext(r), name, bio, email)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Handle Profile View
func handleUserView(w http.ResponseWriter, r *http.Request) {
	userEmail := r.FormValue("email")
	if userEmail == "" {
		userEmail = getUserEmail(r)
	}
	user, err := getUser(getContext(r), userEmail)
	if err != nil {
		//http.Redirect(w,r, "/useredit", http.StatusFound)
	}
	images, err := getImageRecordsByEmail(getContext(r), userEmail)
	if err != nil {
		fmt.Fprintf(w, "error getting images. %v", err)
		return
	}
	tmplData := map[string]interface{}{
		"username": getUserEmail(r),
		"name":        user.Name,
		"bio":         user.Bio,
		"email":       user.Email,
		"lastUpdated": user.LastUpdated.Format(time.UnixDate),
		"loginURL":  getLoginURL(r, "/"),
		"logoutURL": getLogoutURL(r, ""),
		"images":      images,

	}
	renderTemplate(w, "user_view.html", tmplData)
}

// Handle Profile Eidt
func handleUserEdit(w http.ResponseWriter, r *http.Request) {
	email := getUserEmail(r)
	if email == "" {
		redirToLogin(w, r)
		return
	}
	// pre-populate input fields if record exists
	user, err := getUser(getContext(r), email)
	var tmplData map[string]string

	if err != nil {
		tmplData = map[string]string{
			"name":        "<enter your name>",
			"bio":         "<enter your bio>",
			"lastUpdated": "not found",
		}
	} else {
		tmplData = map[string]string{
			"name":        user.Name,
			"bio":         user.Bio,
			"lastUpdated": user.LastUpdated.Format(time.UnixDate),
			"username": getUserEmail(r),
			"loginURL":  getLoginURL(r, "/"),
			"logoutURL": getLogoutURL(r, ""),
		}
	}
	tmplData["email"] = email
	renderTemplate(w, "user_edit.html", tmplData)
}
