package vrmp

import (
	"io"
	"log"
	"net/http"
	"time"
)

// Handles Image Upload Page
func handleImageUpload(w http.ResponseWriter, r *http.Request) {
	postURL := getImageUploadURL(getContext(r))
	tmplData := map[string]string{
		"username": getUserEmail(r),
		"uploadURL": postURL,
		"loginURL":  getLoginURL(r, "/"),
		"logoutURL": getLogoutURL(r, ""),
	}
	renderTemplate(w, "image_upload.html", tmplData)
}

// Handles Image Upload Page for Guests
func handleGuestImageUpload(w http.ResponseWriter, r *http.Request) {
	postURL := getImageUploadURL(getContext(r))
	tmplData := map[string]string{
		"uploadURL": postURL,
		"loginURL":  getLoginURL(r, "/"),
		"logoutURL": getLogoutURL(r, ""),
	}
	renderTemplate(w, "image_upload.html", tmplData)
}

// Handles Image Upload Page After It's Done
func handleImageUploadComplete(w http.ResponseWriter, r *http.Request) {
	ctx := getContext(r)
	// XXX: blobstore parse must be called before r.FormValue
	blobKey, otherValues := getImageParseUploadKey(ctx, r)
	name := otherValues.Get("name")
	log.Println("Storing", name, " blob key:", blobKey)
	if blobKey != "" {
		updateImageRecord(ctx, name, blobKey, getUserEmail(r))
		http.Redirect(w, r, "/imageview?blobkey="+blobKey, http.StatusFound)
	}	
	http.Redirect(w, r, "/", http.StatusFound)
}

// Handle Image View Page
func handleImageView(w http.ResponseWriter, r *http.Request) {
	userEmail := getUserEmail(r)
	blobKey := r.FormValue("blobkey")
	if blobKey == "" {
		return
	}
	img, err := getImageRecord(getContext(r), blobKey)
	if err != nil {
		io.WriteString(w, "image record not found")
		return
	}
	tmplData := map[string]string{
		"name":         img.Name,
		"blobKey":      img.BlobKey,
		"email":        img.Email,
		"timeUploaded": img.TimeUploaded.Format(time.UnixDate),
		"userEmail":    userEmail,
		"loginURL": 	getLoginURL(r, "/"),
		"logoutURL": 	getLogoutURL(r, ""),
		"username": 	getUserEmail(r),
	}
	renderTemplate(w, "image_view.html", tmplData)
}

// Handles Serving Images
func handleImageServe(w http.ResponseWriter, r *http.Request) {
	blobKey := r.FormValue("blobkey")
	if blobKey != "" {
		serveImageByKey(w, blobKey)
	}
}

// Handles Image Deletion
func handleImageDelete(w http.ResponseWriter, r *http.Request) {
	email := getUserEmail(r)
	blobkey := r.FormValue("blobkey")
	ctx := getContext(r)

	record, err := getImageRecord(ctx, blobkey)
	if err != nil {
		io.WriteString(w, "Error getting blob")
		log.Println(err)
		return
	}

	if record.Email != email {
		io.WriteString(w, "You can only delete images uploaded by you")
		return
	}

	// delete the blob
	deleteImageBlob(ctx, blobkey)
	// delete the image record
	deleteImageRecord(ctx, blobkey)

	http.Redirect(w, r, "/userview", http.StatusFound)
}

// Handles the Editing of Captions
func handleCaptionEdit(w http.ResponseWriter, r *http.Request) {
	email := getUserEmail(r)
	blobkey := r.FormValue("blobkey")
	ctx := getContext(r)
	caption := r.FormValue("captionEdit")

	record, err := getImageRecord(ctx, blobkey)

	if err != nil {
		io.WriteString(w, "Error getting blob")
		log.Println(err)
		return
	}

	if record.Email != email {
		io.WriteString(w, "You can only edit images uploaded by you")
		return
	}

	updateImageRecord(ctx, caption, blobkey, email)

	http.Redirect(w, r, "/userview", http.StatusFound)
}


