package main

import (
	"log"
	"os"
	"strconv"
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
		if update.CallbackQuery != nil {
			splittedText := strings.Split(update.CallbackQuery.Data, "_")
			callback := splittedText[0]

			switch callback {
			case "nxt", "prev":
				offset, err := strconv.Atoi(splittedText[1])
				if err != nil {
					log.Printf(
						logMsg,
						update.Message.Chat.FirstName,
						update.Message.Chat.ID,
						"",
						"next|prev",
						err.Error(),
					)
				} else {
					msgText, keyboardPage, err := updateKeyboardPage(update.CallbackQuery.Message.Chat.ID, int64(offset))
					if err != nil {
						log.Printf(
							logMsg,
							update.Message.Chat.FirstName,
							update.Message.Chat.ID,
							"",
							"next|prev",
							err.Error(),
						)
					} else {
						msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "")
						msg.Text = msgText
						msg.ReplyMarkup = keyboardPage
						msg.ParseMode = "Markdown"
						bot.Send(msg)
					}
				}
			}

			bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))
		}
		if update.Message != nil {
			if update.Message.IsCommand() {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				msg.ParseMode = "Markdown"

				switch update.Message.Command() {
				case "start", "help":
					msg.Text = startAndHelpMsg
				case "add":
					err = processingAddCommmand(&msg, &update)
					if err != nil {
						log.Printf(
							logMsg,
							update.Message.Chat.FirstName,
							update.Message.Chat.ID,
							"/add",
							"",
							err.Error(),
						)
					}
				case "delete":
					err = processingDeleteCommmand(&msg, &update)
					if err != nil {
						log.Printf(
							logMsg,
							update.Message.Chat.FirstName,
							update.Message.Chat.ID,
							"/delete",
							"",
							err.Error(),
						)
					}
				case "all":
					err = processingAllCommmand(&msg, &update)
					if err != nil {
						log.Printf(
							logMsg,
							update.Message.Chat.FirstName,
							update.Message.Chat.ID,
							"/all",
							"",
							err.Error(),
						)
					}
				case "my":
					err = processingMyCommmand(&msg, &update)
					if err != nil {
						log.Printf(
							logMsg,
							update.Message.Chat.FirstName,
							update.Message.Chat.ID,
							"/my",
							"",
							err.Error(),
						)
					}
				default:
					msg.Text = unknownCommandMsg
				}

				bot.Send(msg)
			}
		}
	}
}
