package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookInfo struct {
	ID       primitive.ObjectID `bson:"_id"`
	ByChatID int64              `bson:"by_chat_id"`
	Title    string             `bson:"title"`
	Score    int                `bson:"score"`
}
