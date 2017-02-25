package vrmp

import (
	"time"

	"log"

	"appengine"
	"appengine/datastore"
)

// User Structure
type User struct {
	Name        string
	Bio         string
	Email       string
	LastUpdated time.Time
}

// Updates User Record
func updateUser(ctx appengine.Context, name, bio, email string) {
	user := &User{
		Name:        name,
		Bio:         bio,
		Email:       email,
		LastUpdated: time.Now(),
	}
	// use email as StringKey
	key := datastore.NewKey(ctx, "User", email, 0, nil)
	_, err := datastore.Put(ctx, key, user)
	if err != nil {
		// handle err
		log.Println(err)
		return
	}
}

// Gets User
func getUser(ctx appengine.Context, email string) (User, error) {
	key := datastore.NewKey(ctx, "User", email, 0, nil)
	var user User
	err := datastore.Get(ctx, key, &user)
	if err != nil {
		log.Println(err)
	}
	return user, err
}

// Gets all Users
func getAllUsers(ctx appengine.Context) []User {
	users := []User{}
	q := datastore.NewQuery("User")
	t := q.Run(ctx)
	for {
		var e User
		_, err := t.Next(&e)
		if err == datastore.Done {
			break // done
		}
		if err != nil {
			log.Printf("error fetching next person %v", err)
			break
		}
		// push to arr
		users = append(users, e)
	}
	return users
}
