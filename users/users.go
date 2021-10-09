package users

import post "appointy/post"

type User struct {
	Name  string      `json:"name" bson:"name"`
	ID    string      `json:"id" bson:"id"`
	Email string      `json:"email" bson:"email"`
	Posts []post.Post `json:"posts" bson:"posts"`
}
