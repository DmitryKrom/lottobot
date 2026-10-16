package main

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func sendFaultMessage(level int, msg string, ctx context.Context, b *bot.Bot, update *models.Update) {
	switch level {
	case LevelLebensmittel:
		sendTextMessage(msg, ctx, b, update)
	case LevelJP:
		sendTextMessage(msg, ctx, b, update)
	case Level6aus49:
		sendTextMessage(msg, ctx, b, update)
	case LevelObst:
		sendTextCallBack(msg, ctx, b, update)
		// sendTextMessage(msg, ctx, b, update)
	case LevelWitz:
		sendTextMessage(msg, ctx, b, update)
	case LevelMenu:
		sendTextMessage(msg, ctx, b, update)
	}
}

func sendTextCallBack(txt string, ctx context.Context, b *bot.Bot, update *models.Update) {
	mes_Params := bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   txt,
	}
	_, err := b.SendMessage(ctx, &mes_Params)
	if err != nil {
		fmt.Println("[ERROR CINESTAR] [SEND Message failed DUE TO ]", err)
	}
}

func sendTextMessage(txt string, ctx context.Context, b *bot.Bot, update *models.Update) {
	mes_Params := bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   txt,
	}
	_, err := b.SendMessage(ctx, &mes_Params)
	if err != nil {
		fmt.Println("[ERROR CINESTAR] [SEND Message failed DUE TO ]", err)
	}
}
