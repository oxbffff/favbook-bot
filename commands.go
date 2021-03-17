package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
)

func processAllCommand(msg *tgbotapi.MessageConfig, update *tgbotapi.Update) error {
	allBooks, err := getAllBooks(update.Message.Chat.ID)
	if err != nil {
		msg.Text = fmt.Sprintf(errorMsg, err.Error())
		return err
	}

	var booksList string
	for i, bookInfo := range allBooks {
		booksList += fmt.Sprintf("%d) %s - %d\n", i+1, bookInfo.Title, bookInfo.Score)
	}
	msg.Text = booksList

	return nil
}

func processAddCommand(msg *tgbotapi.MessageConfig, update *tgbotapi.Update) error {
	bookTitle, err := getTitle(update.Message.Text)
	if err != nil {
		msg.Text = fmt.Sprintf(errorMsg, err.Error())
		return err
	}
	err = addNewBook(update.Message.Chat.ID, bookTitle)
	if err != nil {
		msg.Text = fmt.Sprintf(errorMsg, "unexpected error while adding the book")
		return err
	}

	msg.Text = "Book was added"

	return nil
}

func processDeleteCommand(msg *tgbotapi.MessageConfig, update *tgbotapi.Update) error {
	bookTitle, err := getTitle(update.Message.Text)
	if err != nil {
		msg.Text = fmt.Sprintf(errorMsg, err.Error())
		return err
	}

	err = deleteBook(update.Message.Chat.ID, bookTitle)
	if err != nil {
		msg.Text = fmt.Sprintf(errorMsg, "unexpected error while deleting the book")
		return err
	}

	msg.Text = "Book was deleted"

	return nil
}

func processScoreCommand(msg *tgbotapi.MessageConfig, update *tgbotapi.Update) error {
	splittedText := strings.Split(update.Message.Text, " ")
	fmt.Println(splittedText)
	if len(splittedText) != 3 {
		err := errors.New("wrong command form")
		msg.Text = fmt.Sprintf(errorMsg, err.Error())
		return err
	}

	bookNumber, err := strconv.Atoi(splittedText[1])
	if err != nil || bookNumber <= 0 {
		err := errors.New("second argument is wrong")
		msg.Text = fmt.Sprintf(errorMsg, err.Error())
		return err
	}
	score, err := strconv.Atoi(splittedText[2])
	if err != nil || score < 0 {
		err := errors.New("third argument is wrong")
		msg.Text = fmt.Sprintf(errorMsg, err.Error())
		return err
	}

	err = updateBookScore(update.Message.Chat.ID, int64(bookNumber)-1, score)
	switch err {
	case nil:
		msg.Text = "Info was updated"
		return nil
	case mongo.ErrNoDocuments:
		msg.Text = fmt.Sprintf(errorMsg, "no matches found for this index")
		return err
	default:
		msg.Text = fmt.Sprintf(errorMsg, "unexpected error")
		return err
	}
}
