package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

func updateKeyboardPage(chatID int64, offset int64) (string, *tgbotapi.InlineKeyboardMarkup, error) {
	client, err := getMongoClient()
	if err != nil {
		return "", nil, err
	}
	defer client.Disconnect(context.TODO())
	collection := client.Database(os.Getenv("DB")).Collection(os.Getenv("BOOKS_COLLECTION"))

	filter := bson.D{{"by_chat_id", chatID}}
	opts := options.FindOneOptions{
		Skip: &offset,
	}
	result := collection.FindOne(context.TODO(), filter, &opts)

	bookInfo := BookInfo{}
	result.Decode(&bookInfo)

	msgText := fmt.Sprintf(bookDescriptionMsg, bookInfo.Title, bookInfo.Score)
	keyboardPage := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("+", "inc_"+bookInfo.ID.String()),
			tgbotapi.NewInlineKeyboardButtonData("-", "dec_"+bookInfo.ID.String()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<", "prev_"+strconv.FormatInt(offset-1, 10)),
			tgbotapi.NewInlineKeyboardButtonData(">", "nxt_"+strconv.FormatInt(offset+1, 10)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Delete", "del_"+bookInfo.ID.String()),
		),
	)

	return msgText, &keyboardPage, nil
}
