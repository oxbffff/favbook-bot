package main

const (
	startAndHelpMsg = "This bot can help you manage your book library.\n" +
		"First of all, you can save the book you have read, rate it and compare it with other books.\n" +
		"Commands:\n" +
		"/add `title` where title is the title of the book \n" +
		"/delete `title` where title is the title of the book \n" +
		"/all - bot sends you a text message with all your books and their score\n" +
		"/my - bot sends you a inline keyboard to manage your saved books\n" +
		"Additional features in development"
	unknownCommandMsg = "Unknow command. Use `/start` or `/help` for more info"
	errorMsg   = "An error occurred: %s"
)
