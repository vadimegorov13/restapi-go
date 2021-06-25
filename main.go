package main

import (
	"encoding/json" // encoding for json file
	"log"           // can log errors
	"net/http"      // can work woth http stuff to create api and stuff

	"math/rand" // generate a rand number
	"strconv"   //string converter

	"github.com/gorilla/mux"
)

// Post struct (Model)
type Post struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Text   string  `json:"text"`
	Author *Author `json:"author"`
}

// Author struct (Model)
type Author struct {
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
}

// Init posts var as a slice Post struct
// Slice is a variable of length array
var posts []Post

// Get all posts (getPosts)
func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Need this or else json will return just a text
	json.NewEncoder(w).Encode(posts)
}

// Get one post by id (getPost)
func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Need this or else json will return just a text
	params := mux.Vars(r)                              // Get params
	// Loop through the posts and find the wirght one by id
	for _, post := range posts {
		if post.ID == params["id"] {
			json.NewEncoder(w).Encode(post)
			return
		}
	}
	json.NewEncoder(w).Encode(&Post{}) // Give response of one post
}

// Create post (createPost)
func createPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Need this or else json will return just a text
	var post Post
	_ = json.NewDecoder(r.Body).Decode(&post)
	post.ID = strconv.Itoa(rand.Intn(10000000)) // Mock ID - not safe for production because it can create the same id
	posts = append(posts, post)
	json.NewEncoder(w).Encode(post) // Give the response of one post
}

// Update post (updatePost)
func updatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Need this or else json will return just a text
	params := mux.Vars(r)                              // Get params
	for index, post := range posts {
		if post.ID == params["id"] {
			posts = append(posts[:index], posts[index+1:]...) // Its kinda like a slice in js
			var post Post
			_ = json.NewDecoder(r.Body).Decode(&post)
			post.ID = params["id"]
			posts = append(posts, post)
			json.NewEncoder(w).Encode(post) // Give the response of one post
			return
		}
	}
}

// Delete Post (deletePost)
func deletePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // Need this or else json will return just a text
	params := mux.Vars(r)                              // Get params
	for index, post := range posts {
		if post.ID == params["id"] {
			posts = append(posts[:index], posts[index+1:]...) // Its kinda like a slice in js
			break
		}
	}
	json.NewEncoder(w).Encode(posts) // Give the response of all the posts
}

func main() {
	// Init Router
	router := mux.NewRouter()

	// Mock data - @todo - implement data base
	posts = append(posts, Post{ID: "1", Title: "Post 1", Text: "This is my first post", Author: &Author{
		First_name: "Vadim", Last_name: "Egorov",
	}})

	posts = append(posts, Post{ID: "2", Title: "Post 2", Text: "This is a post by bob", Author: &Author{
		First_name: "Bob", Last_name: "Bob",
	}})

	// Router Handlers / Endpoints
	router.HandleFunc("/api/posts", getPosts).Methods("GET")
	router.HandleFunc("/api/posts/{id}", getPost).Methods("GET")
	router.HandleFunc("/api/posts", createPost).Methods("POST")
	router.HandleFunc("/api/posts/{id}", updatePost).Methods("PUT")
	router.HandleFunc("/api/posts/{id}", deletePost).Methods("DELETE")

	// log.Fatal will throw an error if something goes wrong
	log.Fatal(http.ListenAndServe(":8000", router)) // Same as app.listen in node js
}
