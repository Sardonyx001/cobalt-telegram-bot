package main

import (
	"net/url"
	"os"
	"os/signal"
	"syscall"

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
		log.Fatalf("Couldn't start Telegram bot: %v", err)
	}

	// Handle graceful shutdown
	c := make(chan os.Signal, 3)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		bot.StopReceivingUpdates()
		os.Exit(0)
	}()

	bot.Debug = true
	log.Infof("Authorized on account %s", bot.Self.UserName)
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.ReplyToMessageID = update.Message.MessageID

			// Check if the message *is* a valid URL and grab it
			u, err := url.ParseRequestURI(update.Message.Text)
			log.Info("url.ParseRequestURI", "url", u)
			if err != nil {
				msg.Text += "Send a valid URL pls"
			} else {
				msg.Text += "Select a download option"
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

			msg := tgbotapi.NewEditMessageText(
				update.CallbackQuery.Message.Chat.ID,
				update.CallbackQuery.Message.MessageID,
				"",
			)

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
				var file tgbotapi.Chattable
				switch callback.Text {
				case "audio":
					file = tgbotapi.NewAudio(update.CallbackQuery.Message.Chat.ID, tgbotapi.FileURL(result.URL))
				default:
					file = tgbotapi.NewVideo(update.CallbackQuery.Message.Chat.ID, tgbotapi.FileURL(result.URL))
				}

				if _, err := bot.Send(file); err != nil {
					log.Fatal(err)
				}
			}

		}
	}
}
