package main

const (
	startAndHelpMsg = "This bot can help you manage your book library.\n" +
		"First of all, you can save the book you have read, rate it and compare it with other books.\n" +
		"Commands:\n" +
		"/add `title` where title is the title of the book \n" +
		"/delete `title` where title is the title of the book \n" +
		"/all - bot sends you a text message with all your books and their score\n" +
		"/score `book_number_in_list` `score` - book rating update\n" +
		"Additional features in development"
	unknownCommandMsg = "Unknow command. Use `/start` or `/help` for more info"
	errorMsg          = "An error occurred: %s"

	logMsg = "Name: %s\nChatID: %d\nCommand: %s\nCallback: %s\nError description: %s\n"
)
