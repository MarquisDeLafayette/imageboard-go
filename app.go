package vrmp

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Initializes App
func init() {
	http.Handle("/", getHandlers())
}

// Handles Page/Function Requests
func getHandlers() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", handleIndex).Methods("GET")

	router.HandleFunc("/userview", handleUserView).Methods("GET")
	router.HandleFunc("/useredit", handleUserEdit).Methods("GET")
	router.HandleFunc("/userupdate", handleUserUpdate).Methods("POST")
	
	router.HandleFunc("/imageupload", handleImageUpload).Methods("GET")
	router.HandleFunc("/guestimageupload", handleGuestImageUpload).Methods("GET")
	router.HandleFunc("/imageupload_complete", handleImageUploadComplete).Methods("POST")
	router.HandleFunc("/imageview", handleImageView).Methods("GET")
	router.HandleFunc("/serveimage", handleImageServe).Methods("GET")
	router.HandleFunc("/imagedelete", handleImageDelete).Methods("GET")
	router.HandleFunc("/imageedit", handleCaptionEdit).Methods("POST")

	return router
}