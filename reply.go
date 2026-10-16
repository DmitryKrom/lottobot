package main

import "github.com/go-telegram/bot/models"

type Menu struct {
	Text     string
	KeyBoard models.ReplyMarkup
}

func NewMenu(art string, color string, replyText string) Menu {
	var menu = Menu{}
	switch art {
	case MENUE:
		menu.Text = replyText
		kb := &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{{Text: "Kino", Style: color}},
				{{Text: "Lotto", Style: color}},
				// {{Text: "Strom", Style: color}},
				{{Text: "Todo", Style: color}},
				{{Text: "Witz", Style: color}},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
			Selective:       true,
		}
		menu.KeyBoard = kb
	case CINEMA:
		menu.Text = replyText
		kb := &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "HEUTE", Style: color}, {Text: "WOCHE", Style: color},
				},
				// {
				// 	{Text: "VORVERKAUF", Style: color}, {Text: "EVENTS", Style: color},
				// },
				{
					{Text: "Menue", Style: color},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
			Selective:       true,
		}
		menu.KeyBoard = kb
	case LOTTO:
		menu.Text = replyText
		kb := &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "Lotto-6aus49", Style: color}, {Text: "Eurojackpot", Style: color},
				},
				{
					{Text: "Menue", Style: color},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
			Selective:       true,
		}
		menu.KeyBoard = kb
	case WITZ:
		menu.Text = replyText
		kb := &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "NOCH", Style: color},
				},
				{
					{Text: "Menue", Style: color},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
			Selective:       true,
		}
		menu.KeyBoard = kb

	case TODO:
		menu.Text = replyText
		kb := &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "Einkauf", Style: color}, {Text: "Termine", Style: color},
				},
				{
					{Text: "Menue", Style: color},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
			Selective:       true,
		}
		menu.KeyBoard = kb
	case "Einkauf":
		menu.Text = replyText
		kb := &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "Lebensmittel", Style: color}, {Text: "Wäsche", Style: color},
				},
				{
					{Text: "Kleidung", Style: color}, {Text: "Küche", Style: color},
				},
				{
					{Text: "Zurück", Style: color}, {Text: "Menue", Style: color},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
			Selective:       true,
		}
		menu.KeyBoard = kb

	case "Lebensmittel":
		menu.Text = replyText
		kb := &models.ReplyKeyboardMarkup{
			Keyboard: [][]models.KeyboardButton{
				{
					{Text: "Obst/Gemüse", Style: color}, {Text: "Milchprodukte", Style: color}, {Text: "Backwaren", Style: color},
				},
				{
					{Text: "Fisch/Fleisch", Style: color}, {Text: "Getränke", Style: color}, {Text: "Getreideprodukte", Style: color},
				},
				{
					{Text: "Süßes", Style: color}, {Text: "Tiefkühlkost", Style: color}, {Text: "Konserven", Style: color},
				},
				{
					{Text: "Menue", Style: color}, {Text: "OK", Style: btnStyleGreen},
				},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: false,
			Selective:       true,
		}
		menu.KeyBoard = kb

	}

	return menu
}
