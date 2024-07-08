package main

import (
	"net/url"
	"os"

	log "github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/lostdusty/gobalt"
)

type Logger struct {
	logger *log.Logger
}

func (l *Logger) Println(v ...interface{}) {
	l.logger.Error(v)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

func NewLogger() tgbotapi.BotLogger {
	logger := log.New(os.Stdout)
	logger.SetLevel(log.DebugLevel)
	return &Logger{logger: logger}
}

func main() {
	godotenv.Load()
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		panic(err)
	}

	tgbotapi.SetLogger(NewLogger())
	bot.Debug = true

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	// Start polling Telegram for updates.
	updates := bot.GetUpdatesChan(updateConfig)

	// Let's go through each update that we're getting from Telegram.
	for update := range updates {
		// Telegram can send many types of updates depending on what your Bot
		// is up to. We only want to look at messages for now, so we can
		// discard any other updates.
		if update.Message == nil {
			continue
		}

		// Take the Chat ID and Text from the incoming message
		// and use it to create a new message.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		// We'll also say that this message is a reply to the previous message.
		msg.ReplyToMessageID = update.Message.MessageID

		// Check if the message *is* a valid URL and grab it
		// Example of valid messages:
		// `https://www.youtube.com/watch?v=dQw4w9WgXcQ`
		u, err := url.ParseRequestURI(update.Message.Text)
		if err != nil {
			msg.Text += "Send a valid URL pls"
		} else {

			//Creates a Settings struct with default values, and save it to downloadMedia variable.
			downloadMedia := gobalt.CreateDefaultSettings()

			// Sets the URL, you MUST set one before downloading the media.
			downloadMedia.Url = u.String()

			// After changing the url, Run() will make the necessary requests to cobalt to download your media
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

		// Okay, we're sending our message off! We don't care about the message
		// we just sent, so we'll discard it.
		if _, err := bot.Send(msg); err != nil {
			// Note that panics are a bad way to handle errors. Telegram can
			// have service outages or network errors, you should retry sending
			// messages or more gracefully handle failures.
			panic(err)
		}
	}
}
