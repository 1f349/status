package notifier

import (
	tgBotApi5 "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
)

type Telegram struct {
	bot    *tgBotApi5.BotAPI
	chatId int64
}

func newTelegram() *Telegram {
	bot, err := tgBotApi5.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		panic(err)
	}
	chatId, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		panic(err)
	}
	return &Telegram{bot: bot, chatId: chatId}
}

func (t *Telegram) SendMessage(a string) {
	_, err := t.bot.Send(tgBotApi5.NewMessage(t.chatId, a))
	if err != nil {
		log.Printf("[Telegram] %s\n", err)
		return
	}
}
