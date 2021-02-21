package main

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func addNewBook(chatID int64, bookTitle string) error {
	client, err := getMongoClient()
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())
	collection := client.Database(os.Getenv("DB")).Collection(os.Getenv("BOOKS_COLLECTION"))

	bookInfo := BookInfo{ID: primitive.NewObjectID(), ByChatID: chatID, Title: bookTitle, Score: 0}
	_, err = collection.InsertOne(context.TODO(), bookInfo)
	if err != nil {
		return err
	}

	return nil
}

func getAllBooks(chatID int64) ([]BookInfo, error) {
	filter := bson.D{{"by_chat_id", chatID}}
	allBooks := []BookInfo{}

	client, err := getMongoClient()
	if err != nil {
		return allBooks, err
	}
	defer client.Disconnect(context.TODO())
	collection := client.Database(os.Getenv("DB")).Collection(os.Getenv("BOOKS_COLLECTION"))

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return allBooks, err
	}

	for cur.Next(context.TODO()) {
		bookInfo := BookInfo{}
		err = cur.Decode(&bookInfo)
		if err != nil {
			return allBooks, err
		}
		allBooks = append(allBooks, bookInfo)
	}
	cur.Close(context.TODO())

	if len(allBooks) == 0 {
		return allBooks, errors.New("found 0 books")
	}

	return allBooks, nil
}
