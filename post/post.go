package post

import "go.mongodb.org/mongo-driver/bson/primitive"

type Post struct {
	UserID   string              `jsong:"userID" bson:"userID"`
	ID       string              `json:"id" bson:"id"`
	Caption  string              `json:"caption" bson:"caption"`
	PostTime primitive.Timestamp `json:"post_time" bson:"post_time"`
}
