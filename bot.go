package main

import (
	"log"
	"os"

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
			switch update.Message.Command() {
			case "help", "start":
				msg.Text = startAndHelpMsg
			default:
				msg.Text = "Unknow command"
			}
			bot.Send(msg)
		}
	}
}
