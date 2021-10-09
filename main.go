package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Name  string `json:"name" bson:"name"`
	ID    string `json:"id" bson:"id"`
	Email string `json:"email" bson:"email"`
	Posts []Post `json:"posts" bson:"posts"`
}
type Post struct {
	userID   string              `jsong:"userID" bson:"userID"`
	ID       string              `json:"id" bson:"id"`
	Caption  string              `json:"caption" bson:"caption"`
	PostTime primitive.Timestamp `json:"post_time" bson:"post_time"`
}

var client *mongo.Client

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server Running!")
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	collection := client.Database("api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(w).Encode(res)
}

func ListUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var users []User
	collection := client.Database("api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(`message: ` + err.Error()))
		return
	}
	defer res.Close(ctx)
	for res.Next(ctx) {
		var user User
		res.Decode(&user)
		users = append(users, user)
	}
	if err := res.Err(); err != nil {
		w.WriteHeader(404)
		return
	}
	json.NewEncoder(w).Encode(users)
}

func NoAction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	fmt.Fprintf(w, "No action found.")
}

func Serve(w http.ResponseWriter, r *http.Request) {
	var listUser = regexp.MustCompile(`^\/users\/(\d+)$`)
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost:
		fmt.Fprintf(w, "POST method! User created!")
		CreateUser(w, r)
		return
	case r.Method == http.MethodGet && listUser.MatchString(r.URL.Path):
		fmt.Fprintf(w, "Get method! User listed!")
		ListUser(w, r)
		return
	default:
		NoAction(w, r)
		return
	}
}
func CreatePost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var post Post
	json.NewDecoder(r.Body).Decode(&post)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := client.Database("api").Collection("users")
	res, err := collection.UpdateOne(ctx, bson.M{"id": post.userID},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "posts", Value: post}}},
		})
	json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Updated %v Documents!\n", res.ModifiedCount)
}
func ListPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	var post Post
	collection := client.Database("api").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	res, err := collection.Find(ctx, bson.M{"posts.id": r.URL.Path.rsplit('/', 1)})
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte(`message: ` + err.Error()))
		return
	}
	defer res.Close(ctx)
	if err := res.Err(); err != nil {
		w.WriteHeader(404)
		return
	}
	json.NewEncoder(w).Encode(post)
}
func Posts(w http.ResponseWriter, r *http.Request) {
	var listUser = regexp.MustCompile(`^\/posts\/(\d+)$`)
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodPost:
		fmt.Fprintf(w, "POST method! Post created!")
		CreatePost(w, r)
		return
	case r.Method == http.MethodGet && listUser.MatchString(r.URL.Path):
		fmt.Fprintf(w, "Get method! Post listed!")
		ListPost(w, r)
		return
	default:
		NoAction(w, r)
		return
	}
}

func Handle() {
	router := http.ServeMux{}
	router.HandleFunc("/", homePage)
	router.HandleFunc("/users", Serve)
	router.HandleFunc("/users/", Serve)
	router.HandleFunc("/posts", Posts)
	router.HandleFunc("/posts/", Posts)
	log.Fatal(http.ListenAndServe(":8081", &router))
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	Handle()
}
