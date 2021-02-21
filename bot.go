package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(".env not found")
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.OpenFile(os.Getenv("LOG_FILE"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(f)
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TOKEN_FAVBOOK_BOT"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = os.Getenv("BOT_DEBUG") == "true"

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ParseMode = "Markdown"

			switch update.Message.Command() {
			case "start", "help":
				msg.Text = startAndHelpMsg
			case "add":
				splittedText := strings.SplitN(update.Message.Text, " ", 2)
				if len(splittedText) == 1 {
					msg.Text = fmt.Sprintf(errorMsg, "empty title")
				} else {
					bookTitle := splittedText[1]
					err = addNewBook(update.Message.Chat.ID, bookTitle)
					if err != nil {
						msg.Text = fmt.Sprintf(errorMsg, "unexpected error while adding the book")
						log.Println(err)
					} else {
						msg.Text = "Book added"
					}
				}
			case "delete":
				fmt.Println(update.Message.Chat.ID, update.Message.Text)
			case "all":
				allBooks, err := getAllBooks(update.Message.Chat.ID)
				if err != nil {
					msg.Text = fmt.Sprintf(errorMsg, err.Error())
					log.Println(err)
				} else {
					var booksList string
					for i, bookInfo := range allBooks {
						booksList += fmt.Sprintf("%d) %s - %d\n", i+1, bookInfo.Title, bookInfo.Score)
					}
					msg.Text = booksList
				}
			case "my":
				fmt.Println(update.Message.Chat.ID, update.Message.Text)
			default:
				msg.Text = unknownCommandMsg
			}

			bot.Send(msg)
		}
	}
}
