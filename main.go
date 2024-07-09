package main

import (
	"net/url"
	"os"

	log "github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/lostdusty/gobalt"
)

func main() {
	godotenv.Load()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	tgbotapi.SetLogger(NewLogger())
	bot.Debug = true // Enable debug mode
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Take the Chat ID and Text from the incoming message
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.ReplyToMessageID = update.Message.MessageID

		// Check if the message *is* a valid URL and grab it
		// Example of valid messages:
		// `https://www.youtube.com/watch?v=dQw4w9WgXcQ`
		u, err := url.ParseRequestURI(update.Message.Text)
		if err != nil {
			msg.Text += "Send a valid URL pls"
		} else {
			downloadMedia := gobalt.CreateDefaultSettings()
			downloadMedia.Url = u.String()
			result, err := gobalt.Run(downloadMedia)
			if err != nil {
				log.Error(err)
				msg.Text += "Oops, something went wrong"
				msg.Text += "\n" + err.Error()
			} else {
				// Return the url from cobalt to download the requested media.
				// result.URL -> https://us4-co.wuk.sh/api/stream?t=wTn-71aaWAcV2RBejNFN.....
				msg.Text += "status: " + result.Text
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Download here", result.URL)))
			}
		}

		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}
	}
}
