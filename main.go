package main

import (
	"net/url"
	"os"

	log "github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/lostdusty/gobalt"
)

var keyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✨ auto", "auto"),
		tgbotapi.NewInlineKeyboardButtonData("🎶 audio", "audio"),
	),
)

func main() {
	godotenv.Load()
	tgbotapi.SetLogger(NewLogger())
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true
	log.Infof("Authorized on account %s", bot.Self.UserName)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Select a download option")
			msg.ReplyToMessageID = update.Message.MessageID

			// Check if the message *is* a valid URL and grab it
			u, err := url.ParseRequestURI(update.Message.Text)
			log.Info("url.ParseRequestURI", "url", u)
			if err != nil {
				msg.Text += "Send a valid URL pls"
			} else {
				msg.ReplyMarkup = keyboard
			}
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				log.Fatal(err)
			}

			msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "")

			downloadMedia := gobalt.CreateDefaultSettings()
			downloadMedia.Url = update.CallbackQuery.Message.ReplyToMessage.Text

			if callback.Text == "audio" {
				downloadMedia.AudioOnly = true
			}
			result, err := gobalt.Run(downloadMedia)
			if err != nil {
				log.Error(err)
				msg.Text += "Oops, something went wrong"
				msg.Text += "\n" + err.Error()
				if _, err := bot.Send(msg); err != nil {
					log.Fatal(err)
				}
			} else {
				document := tgbotapi.NewDocument(update.CallbackQuery.Message.Chat.ID, tgbotapi.FileURL(result.URL))
				_, err := bot.Send(document)
				if err != nil {
					log.Fatal(err)
				}
			}

		}
	}
}
