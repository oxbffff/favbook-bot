package main

import (
	"context"
	"errors"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	opts := options.Find()
	opts.SetSort(bson.D{{"score", -1}})
	cur, err := collection.Find(context.TODO(), filter, opts)
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

func deleteBook(chatID int64, bookTitle string) error {
	filter := bson.D{{"by_chat_id", chatID}, {"title", bookTitle}}

	client, err := getMongoClient()
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())
	collection := client.Database(os.Getenv("DB")).Collection(os.Getenv("BOOKS_COLLECTION"))

	result, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("cannot delete this book or book is not exist")
	}

	return nil
}

func getTitle(text string) (string, error) {
	splittedText := strings.SplitN(text, " ", 2)
	if len(splittedText) == 1 {
		return "", errors.New("empty title")
	}

	bookTitle := splittedText[1]

	return bookTitle, nil
}

func updateBookScore(chatID int64, offset int64, score int) error {
	client, err := getMongoClient()
	if err != nil {
		return err
	}
	defer client.Disconnect(context.TODO())
	collection := client.Database(os.Getenv("DB")).Collection(os.Getenv("BOOKS_COLLECTION"))

	findFilter := bson.D{{"by_chat_id", chatID}}
	opts := options.FindOneOptions{
		Skip: &offset,
	}
	findResult := collection.FindOne(context.TODO(), findFilter, &opts)
	if err = findResult.Err(); err != nil {
		return err
	}
	bookInfo := BookInfo{}
	findResult.Decode(&bookInfo)

	updateFilter := bson.D{{"_id", bookInfo.ID}}
	update := bson.D{{"$set", bson.D{{"score", score}}}}
	_, err = collection.UpdateOne(context.TODO(), updateFilter, update)
	if err != nil {
		return errors.New("cannot update score")
	}

	return nil
}
