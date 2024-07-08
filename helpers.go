package main

import (
	"os"

	log "github.com/charmbracelet/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
