package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type Context struct {
	Token string
}

//tg means telegram, all stuff related to telegram like bot api, http server for bot etc.

func (ctxRef *Context) GetUpdates() {
	ctx := *ctxRef

	bot, err := tgbotapi.NewBotAPI(ctx.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			if update.Message.Text == "rooms" {

			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}

func (ctxRef *Context) setWebhook() {
	ctx := *ctxRef
	bot, err := tgbotapi.NewBotAPI(ctx.Token)
	if err != nil {
		log.Panic(err)
	}

	params := make(tgbotapi.Params)
	params["url"] = "https://34e0-85-174-195-48.ngrok.io"

	request, err := bot.MakeRequest("setWebhook", params)
	if err != nil {
		return
	}
	_ = len(request.Description)
}

func (ctxRef *Context) DeleteWebhook() {
	ctx := *ctxRef
	bot, err := tgbotapi.NewBotAPI(ctx.Token)
	if err != nil {
		log.Panic(err)
	}

	request, err := bot.MakeRequest("deleteWebhook", make(tgbotapi.Params))
	if err != nil {
		return
	}
	_ = len(request.Description)
}
