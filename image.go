package vrmp

import (
	"log"
	"time"

	"net/http"
	"net/url"

	"appengine"
	"appengine/blobstore"
	"appengine/datastore"

)

// Image Structure
type ImageRecord struct {
	Name         string
	BlobKey      string
	Email        string
	TimeUploaded time.Time
}

// Gets Image Upload URL
func getImageUploadURL(ctx appengine.Context) string {
	uploadURL, err := blobstore.UploadURL(ctx, "/imageupload_complete", nil)
	if err != nil {
		log.Println("Error generating upload URL", err)
		return ""
	}
	return uploadURL.String()
}

// Gets Image by Blobkey
func serveImageByKey(w http.ResponseWriter, key string) {
	blobstore.Send(w, appengine.BlobKey(key))
}

// Gets Image Parse Upload Key
func getImageParseUploadKey(ctx appengine.Context, r *http.Request) (string, url.Values) {
	blobs, otherVals, err := blobstore.ParseUpload(r)
	if err != nil {
		log.Println("Error getting blob key!", err)
		return "", otherVals
	}
	file := blobs["file"]
	if len(file) == 0 {
		log.Println("Error! len file is 0", err)
		return "", otherVals
	}
	return string(file[0].BlobKey), otherVals
}

// Delete Image Blob
func deleteImageBlob(ctx appengine.Context, blobkey string) {
	blobstore.Delete(ctx, appengine.BlobKey(blobkey))
}

// Deletes Record of Image
func deleteImageRecord(ctx appengine.Context, blobkey string) {
	key := datastore.NewKey(ctx, "ImageRecord", blobkey, 0, nil)
	datastore.Delete(ctx, key)
}

// Updates Database of Images
func updateImageRecord(ctx appengine.Context, name, blobkey, email string) {
	img := &ImageRecord{
		Name:         name,
		BlobKey:      blobkey,
		Email:        email,
		TimeUploaded: time.Now(),
	}
	// use blobkey as key
	key := datastore.NewKey(ctx, "ImageRecord", blobkey, 0, nil)
	_, err := datastore.Put(ctx, key, img)
	if err != nil {
		// handle err
		log.Println(err)
		return
	}
}

// Gets Record of Images
func getImageRecord(ctx appengine.Context, blobKey string) (ImageRecord, error) {
	key := datastore.NewKey(ctx, "ImageRecord", blobKey, 0, nil)
	var img ImageRecord
	err := datastore.Get(ctx, key, &img)
	if err != nil {
		log.Println(err)
	}
	return img, err
}

// Gets Images Based off Email
func getImageRecordsByEmail(ctx appengine.Context, email string) ([]ImageRecord, error) {
	q := datastore.NewQuery("ImageRecord").
		Filter("Email = ", email).
		Order("-TimeUploaded")
	var results []ImageRecord
	_, err := q.GetAll(ctx, &results)
	if err != nil {
		log.Println(err)
	}
	return results, err
}

// Gets all Images Ordered by Date Uploaded
func getAllImages(ctx appengine.Context, pageNumber int) []ImageRecord {
	// Create a query for all Person entities.
	q := datastore.NewQuery("ImageRecord").
		Order("-TimeUploaded")

	images := []ImageRecord{}

	// Iterate over the results.
	t := q.Run(ctx) 
	for {
		var i ImageRecord
		_, err := t.Next(&i)
		if err == datastore.Done {
			break
		}
		if err != nil {
			break
		}
		// push to arr
		images = append(images, i)
	}

	// make a subarray
	// XXX: hard code to display 10 images at a time
	fromIndex := (pageNumber - 1) * 10
	toIndex := fromIndex + 10

	if (toIndex > len(images)){
		if(fromIndex > len(images)){
			return images[0:0]
		}else{
			return images[fromIndex:len(images)]
		}
	} else {
		return images[fromIndex:toIndex]
	}
}