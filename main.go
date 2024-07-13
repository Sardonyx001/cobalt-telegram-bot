package main

import (
	"net/url"
	"os"

	log "github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
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
				// downloadMedia := gobalt.CreateDefaultSettings()
				// downloadMedia.Url = u.String()
				// result, err := gobalt.Run(downloadMedia)
				// if err != nil {
				// 	log.Error(err)
				// 	msg.Text += "Oops, something went wrong"
				// 	msg.Text += "\n" + err.Error()
				// } else {
				// 	// Return the url from cobalt to download the requested media.
				// 	// result.URL -> https://us4-co.wuk.sh/api/stream?t=wTn-71aaWAcV2RBejNFN.....
				// 	msg.Text += "status: " + result.Text
				// 	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				// 		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Download here", result.URL)))
				// }
			}
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				log.Fatal(err)
			}

			msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, callback.Text)
			if _, err := bot.Send(msg); err != nil {
				log.Fatal(err)
			}
		}
	}
}
