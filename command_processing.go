package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func processingAllCommmand(msg *tgbotapi.MessageConfig, update *tgbotapi.Update) error {
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

func processingAddCommmand(msg *tgbotapi.MessageConfig, update *tgbotapi.Update) error {
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

func processingDeleteCommmand(msg *tgbotapi.MessageConfig, update *tgbotapi.Update) error {
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

func processingMyCommmand(msg *tgbotapi.MessageConfig, update *tgbotapi.Update) error {
	allBooks, err := getAllBooks(update.Message.Chat.ID)
	if err != nil {
		msg.Text = fmt.Sprintf(errorMsg, err.Error())
		return err
	}
	fmt.Println(allBooks)
	msgText, keyboardPage, err := updateKeyboardPage(update.Message.Chat.ID, 0)
	if err != nil {
		msg.Text = fmt.Sprintf(errorMsg, err.Error())
		return err
	}
	msg.Text = msgText
	msg.ReplyMarkup = keyboardPage

	return nil
}
