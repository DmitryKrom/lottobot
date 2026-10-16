package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"lottobot/kino"
	"lottobot/lotto"
	"lottobot/witz"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	CINEMA            = "Kino"
	LOTTO             = "Lotto"
	LOTTO6AUS49       = "Lotto_6aus49"
	EUROJACK          = "Eurojackpot"
	MENUE             = "Menue"
	BOTSCHAFT         = "/botschaft"
	WITZ              = "Witz"
	TODO              = "Todo"
	EUROJACKPOTBEGINN = "23.03.2012"
	LOTTOBEGINN       = "23.03.2012"
	SCREENSHOTNAMEJP  = "jackpot.png"
	SCREENSHOTLOTTO   = "lotto.png"
	btnStyleRed       = "danger"
	btnStyleGreen     = "success"
	btnStyleBlue      = "primary"
)

type Scope string

var handlerScope int
var JackpotDatumScope Scope
var LottoDatumScope Scope
var Conf_tok map[string]string
var jpDates = make([]string, 0)
var lDates = make([]string, 0)
var kinos = make([]*kino.ShowFilm, 0)
var menuMarkUp = NewMenu(MENUE, btnStyleBlue, " - Wählen sie Ihre Botart - ")
var kinoMarkUp = NewMenu(CINEMA, btnStyleBlue, " - Bitte, Filter wählen - ")
var lottoMarkUp = NewMenu(LOTTO, btnStyleBlue, " - Bitte, Lottoart wählen - ")
var witzMarkUp = NewMenu(WITZ, btnStyleBlue, " - Ok, hier ist Noch einer - ")
var todo = NewMenu(TODO, btnStyleGreen, " TODO ")
var einkauf = NewMenu("Einkauf", btnStyleBlue, " - Wählen Sie eine Liste -")
var lbmittel = NewMenu("Lebensmittel", btnStyleBlue, " - Wählen Sie eine Kategorie! - ")

var obstgemüse = []string{"Auberginie", "Zucchini", "Paprika", "Petersilie", "Dill", "Tomaten", "Gurken", "Melone",
	"Äpfel", "Bananen", "Kiwi", "Birnen", "Weintrauben", "Wassermelone", "Zitronen", "Kopfsalat", "Möhre", "Kartoffel", "r. Zwiebel", "w. Zwiebel", "Schnittlauch", "Süßkartoffel", "Knoblauch", "Erdbeeren", "Avocado"}
var milch = []string{"Milch", "Schmand", "Schlagsahne", "Käse a.S.", "Käse gerieb.", "Hartkäse", "Joghurt", "gr. Joghurt",
	"Weißkäse", "Mozarella", "Butter", "Quark", "Schmelzkäse", "Feta"}
var backwaren = []string{"schw. Brot", "berg. Landbrot", "Brötchen", "Toastbrot", "Backbrötchen", "Croissant", "Mehl", "Hefe", "Backpulver", "Vanillinzucker", "Zuckerpuder", "Schmelzschokolade", "Vanilie"}
var fleischfisch = []string{"Brustfilet", "H.Schenkel", "H.Flügel", "Hackfl. Rind", "Hackfl. Schwein", "Fisch", "Fischstäbchen", "fr.Lachs", "fr.Forelle", "Lachs geräuchert", "Garnellen", "Hering", "Wurst", "Schinken"}
var getreide = []string{"Buchweizen", "Grieß", "Haferflocken", "Müsli", "Kartoffelstärke"}
var sueß = []string{"Schokolade", "Lolli", "Chips"}
var tiefkühl = []string{"Fischstäbchen", "Pelmeni", "Himbeeren", "Blaubeeren", "Lachsfilet", "Eis", "Plombir"}
var konserv = []string{"Thunfisch", "Hering", "Frühstücksfleisch", "", "Lachsfilet", "Eis", "Plombir"}

var lebensMittel = map[string][]string{
	"Obst":     obstgemüse,
	"Milch":    milch,
	"Back":     backwaren,
	"Fleisch":  fleischfisch,
	"Getreide": getreide,
	"Süß":      sueß,
	"Tiefkühl": tiefkühl,
	"Konserve": konserv,
}

//var lebensmittel = [][]string{obstgemüse, milch, backwaren, fleischfisch, getreide, sueß, tiefkühl, konserv}

var lmZettel map[string]int

var inlineKeyboardMarkup = &models.InlineKeyboardMarkup{}
var messageIds []int
var monthOnly bool
var monthList bool

func main() {
	fmt.Println("Hallo From Lottobot ", time.Now().Format("02.01.2006"))
	JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
	LottoDatumScope = Scope(time.Now().Format("02.01.2006"))
	jpDates = getDatesJP(string(JackpotDatumScope))
	lDates = getDatesJP(string(LottoDatumScope))
	fmt.Println(lDates)
	// os.Exit(1)
	bot_token := getToken()
	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
		bot.WithCallbackQueryDataHandler("Filme Heute", bot.MatchTypeContains, handlerFilm),
		bot.WithCallbackQueryDataHandler("Filme Woche", bot.MatchTypeContains, handlerFilm),
		bot.WithWorkers(10),
		// bot.WithCallbackQueryDataHandler("Einkauf", bot.MatchTypeContains, handler_einkauf),
	}
	b, err := bot.New(bot_token, opts...)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	ctx := context.TODO()
	b.RegisterHandler(bot.HandlerTypeMessageText, CINEMA, bot.MatchTypeExact, handler_cinema)
	b.RegisterHandler(bot.HandlerTypeMessageText, LOTTO, bot.MatchTypeExact, handler_lotto_default)
	b.RegisterHandler(bot.HandlerTypeMessageText, LOTTO6AUS49, bot.MatchTypeExact, handler_lotto)
	b.RegisterHandler(bot.HandlerTypeMessageText, EUROJACK, bot.MatchTypeExact, handler_eurojackpot)
	b.RegisterHandler(bot.HandlerTypeMessageText, BOTSCHAFT, bot.MatchTypeExact, handler_botschaft)
	b.RegisterHandler(bot.HandlerTypeMessageText, MENUE, bot.MatchTypeExact, handler_menue)
	b.RegisterHandler(bot.HandlerTypeMessageText, WITZ, bot.MatchTypeExact, handler_witz)
	b.RegisterHandler(bot.HandlerTypeMessageText, TODO, bot.MatchTypeExact, handler_todo)
	b.RegisterHandler(bot.HandlerTypeMessageText, "Einkauf", bot.MatchTypeExact, handler_einkauf)
	b.RegisterHandler(bot.HandlerTypeMessageText, "Lebensmittel", bot.MatchTypeExact, handler_lebensmittel)
	b.RegisterHandler(bot.HandlerTypeMessageText, "Obst/Gemüse", bot.MatchTypeExact, handler_obst)

	botCom_kino := models.BotCommand{
		Command:     CINEMA,
		Description: "kino",
	}
	botCom_lotto := models.BotCommand{
		Command:     LOTTO,
		Description: "lotto",
	}
	botCom_lotto6aus49 := models.BotCommand{
		Command:     LOTTO6AUS49,
		Description: "lotto6aus49",
	}
	botCom_eurojackpot := models.BotCommand{
		Command:     EUROJACK,
		Description: "eurojackpot",
	}
	botCom_menue := models.BotCommand{
		Command:     MENUE,
		Description: "Menü",
	}
	botCommands := []models.BotCommand{botCom_kino, botCom_lotto6aus49, botCom_eurojackpot, botCom_lotto, botCom_menue}

	b.SetMyCommands(ctx, &bot.SetMyCommandsParams{
		Commands: botCommands,
	})

	b.Start(ctx)
}

func getToken() string {

	f, err := os.Open("config/config.json")
	if err != nil {
		fmt.Printf("[ERROR] [OPEN config.json FAILED DUE TO: %v\n]", err)
	}
	//fmt.Println("File opened !!!")
	defer f.Close()
	reader := bufio.NewReader(f)
	if err := json.NewDecoder(reader).Decode(&Conf_tok); err != nil {
		fmt.Printf("[ERROR] [json DECODE FAILED DUE TO: %v\n]", err)
	}
	return Conf_tok["LOTTO_BOT_TOKEN"]
}

func handler_menue(ctx context.Context, b *bot.Bot, update *models.Update) {
	handlerScope = LevelMenu
	fmt.Println(update.Message.From.LastName)
	fmt.Println(update.Message.From.FirstName)
	message_Params := bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        menuMarkUp.Text,
		ReplyMarkup: menuMarkUp.KeyBoard,
	}
	_, err := b.SendMessage(ctx, &message_Params)
	if err != nil {
		fmt.Println("[ERROR] [SENDING MESSAGE FAILED DUE TO ]", err)
	}
}

// func handler2(ctx context.Context, b *bot.Bot, update *models.Update) {

// 	var gewinn_Text string
// 	var message_Params bot.SendMessageParams
// 	// fmt.Println("handler String: ", update.Message.Text)
// 	fmt.Println("in handler ", handlerScope)
// 	switch handlerScope {
// 	case "lebensmittel":
// 		handler_lebensmittel(ctx, b, update)
// 		return
// 	case "obst":
// 		handler_update_inline(ctx, b, update)
// 		return
// 	}

// 	if handlerScope == "obst" {
// 		handler_update_inline(ctx, b, update)
// 		return
// 	}

// 	msg_text := update.Message.Text
// 	fmt.Println("--------------- handler ----------------> ", msg_text)
// 	if msg_text == "NOCH" {
// 		if handlerScope == "witz" {
// 			handlerScope = "witz"
// 			handler_witz(ctx, b, update)
// 		}
// 		return
// 	}
// 	if msg_text == "Obst/Gemüse" {
// 		// handlerScope = "todo Einkauf"
// 		handler_obst(ctx, b, update)
// 		return
// 	}
// 	if msg_text == "Lebensmittel" {
// 		// handlerScope = "todo Einkauf"
// 		handler_lebensmittel(ctx, b, update)
// 		return
// 	}
// 	if msg_text == "Einkauf" {
// 		// handlerScope = "todo Einkauf"
// 		handler_einkauf(ctx, b, update)
// 		return
// 	}
// 	// if msg_text == "Einkaufszettel" {
// 	// 	handlerScope = "todo Einkauf"
// 	// 	handler_zettel(ctx, b, update)
// 	// 	return
// 	// }
// 	if msg_text == "Todo" {
// 		handlerScope = "todo"
// 		handler_todo(ctx, b, update)
// 		return
// 	}
// 	if msg_text == "Menu" {
// 		handlerScope = "menu"
// 		handler_menue(ctx, b, update)
// 		return
// 	}
// 	if msg_text == "Witz" {
// 		handlerScope = "witz"
// 		handler_witz(ctx, b, update)
// 		return
// 	}
// 	if msg_text == "Kino" {
// 		handlerScope = "cinema"
// 		handler_cinema(ctx, b, update)
// 		return
// 	}
// 	if strings.Contains(string(handlerScope), "cinema") {
// 		switch msg_text {
// 		case "HEUTE":
// 			handlerScope = "cinema heute"
// 		case "WOCHE":
// 			handlerScope = "cinema woche"
// 		case "VORVERKAUF":
// 			handlerScope = "cinema vorverkauf"
// 		case "EVENTS":
// 			handlerScope = "cinema events"
// 		case "VORSCHAU":
// 			handlerScope = "cinema vorschau"
// 		}
// 		handler_cinema(ctx, b, update)
// 		return
// 	}

// 	if msg_text == "Menue" {
// 		handler_menue(ctx, b, update)
// 		return
// 	}
// 	if msg_text == "Lotto" {
// 		handler_lotto_default(ctx, b, update)
// 		return
// 	}
// 	if msg_text == "eurojackpot" {
// 		handler_eurojackpot(ctx, b, update)
// 		return
// 	}
// 	if msg_text == "lotto-6aus49" {
// 		handler_lotto(ctx, b, update)
// 		return
// 	}
// 	if msg_text == "ZURÜCK" {
// 		handler_lotto_default(ctx, b, update)
// 		return
// 	}
// 	for x := 1; x <= 12; x++ {
// 		var m time.Month
// 		m = time.Month(x)
// 		if msg_text == m.String() {
// 			if handlerScope == "eurojackpot" {
// 				handler_eurojackpot(ctx, b, update)
// 				return
// 			}
// 			if handlerScope == "lotto" {
// 				handler_lotto(ctx, b, update)
// 				return
// 			}
// 		}
// 	}
// 	if msg_text == "vorher" || msg_text == "nachher" {
// 		if handlerScope == "eurojackpot" {
// 			handler_eurojackpot(ctx, b, update)
// 			return
// 		}
// 		if handlerScope == "lotto" {
// 			handler_lotto(ctx, b, update)
// 			return
// 		}
// 	}

// 	i, err := strconv.ParseInt(msg_text, 10, 64)
// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println(i)
// 	}

// 	if msg_text != "" {
// 		if err := check_datum(msg_text); err != nil {
// 			message_Params = bot.SendMessageParams{
// 				ChatID:          update.Message.Chat.ID,
// 				Text:            err.Error() + "\n",
// 				MessageEffectID: Conf_tok["FAILED_EFFECT"],
// 			}
// 		} else {
// 			artDatum := handlerScope + Scope(msg_text)
// 			gw, err := lotto.Get_lotto(string(artDatum))
// 			// gewinn_Text = build_message(gw, string(lotto_dat))
// 			if err != nil {
// 				fmt.Println("ERROR: ", err.Error())
// 			}
// 			if err != nil {
// 				if err.Error() == "screenshot" {
// 					fmt.Println("screenshot")
// 					switch handlerScope {
// 					case "eurojackpot":
// 						showJackpotScreenshot(update)
// 						return
// 					case "lotto-6aus49":
// 						showLottoScreenshot(update)
// 						return
// 					}

// 				}
// 				message_Params = bot.SendMessageParams{
// 					ChatID:          update.Message.Chat.ID,
// 					Text:            "Ich konnte keine Lottozahlen finden.\nBitte, Datum prüfen !!!\n",
// 					MessageEffectID: Conf_tok["SUCCESS_EFFECT"],
// 				}
// 				_, err := b.SendMessage(ctx, &message_Params)
// 				if err != nil {
// 					fmt.Println("[ERROR] [SENDING MESSAGE FAILED DUE TO ]", err)
// 				}
// 				return
// 			}
// 			if len(gw.Lotto_zahlen) == 0 || len(gw.Zusatz_zahlen) == 0 {
// 				message_Params = bot.SendMessageParams{
// 					ChatID:          update.Message.Chat.ID,
// 					Text:            "Ich konnte keine Lottozahlen finden.\nBitte, Datum prüfen !!!\n",
// 					MessageEffectID: Conf_tok["SUCCESS_EFFECT"],
// 				}
// 				_, err := b.SendMessage(ctx, &message_Params)
// 				if err != nil {
// 					fmt.Println("[ERROR] [SENDING MESSAGE FAILED DUE TO ]", err)
// 				}
// 				return
// 			} else {
// 				for i := 0; i < len(gw.Lotto_zahlen); i++ {
// 					if gw.Lotto_zahlen[i] == 0 {
// 						message_Params = bot.SendMessageParams{
// 							ChatID:          update.Message.Chat.ID,
// 							Text:            "Eine oder mehrere Zahlen sind \"0\"\n",
// 							MessageEffectID: Conf_tok["SUCCESS_EFFECT"],
// 						}
// 						_, err := b.SendMessage(ctx, &message_Params)
// 						if err != nil {
// 							fmt.Println("[ERROR] [SENDING MESSAGE FAILED DUE TO ]", err)
// 						}
// 						return
// 					}
// 				}

// 			}
// 			gewinn_Text = build_message(gw, msg_text)

// 			message_Params = bot.SendMessageParams{
// 				ChatID:          update.Message.Chat.ID,
// 				Text:            gewinn_Text + "\n",
// 				MessageEffectID: Conf_tok["SUCCESS_EFFECT"],
// 			}
// 		}
// 	}

// 	_, err = b.SendMessage(ctx, &message_Params)
// 	if err != nil {
// 		fmt.Println("[ERROR] [SENDING MESSAGE FAILED DUE TO ]", err)
// 	}
// }

func build_message(gw lotto.Play, lotto_dat string) string {
	var text = "Gewinnzahlen vom " + lotto_dat + " \n"
	for i := 0; i < len(gw.Lotto_zahlen); i++ {
		text = text + " " + strconv.Itoa(gw.Lotto_zahlen[i])
	}
	text = text + "\n" + "Zusatzzahlen: \n"
	for i := 0; i < len(gw.Zusatz_zahlen); i++ {
		text = text + " " + strconv.Itoa(gw.Zusatz_zahlen[i])
	}
	return text
}

func check_datum(dat string) error {
	d, _ := time.Parse("02.01.2006", dat)
	if d.After(time.Now()) {
		return errors.New("DAS DATUM LIEGT IN DER ZUKUNFT")
	}
	if strings.Contains(dat, " ") {
		return errors.New("DAS DATUM DARF KEINE LEERSTELLEN BEINHALTEN")
	}
	if strings.Count(dat, ".") != 2 {
		return errors.New("TAG UND MONAT TRENNUNG MIT \".\" FEHLT")
	}
	d_m_y := strings.Split(dat, ".")

	fmt.Println(d_m_y)
	if len(d_m_y[0]) < 2 || len(d_m_y[1]) < 2 {
		return errors.New("TAG- UND MONATANGABEN MÜSSEN 2-STELLIG SEIN")
	}

	return nil
}
func handler_zettel(ctx context.Context, b *bot.Bot, update *models.Update) {
	del_Params := bot.DeleteMessagesParams{
		ChatID:     update.Message.Chat.ID,
		MessageIDs: messageIds,
	}
	_, err := b.DeleteMessages(ctx, &del_Params)
	if err != nil {
		fmt.Println(" ----------------- Delete Message Err ------------------- ", err)
	}
	photo_Params := bot.DeleteChatPhotoParams{
		ChatID: update.Message.Chat.ID,
	}
	_, err = b.DeleteChatPhoto(ctx, &photo_Params)
	if err != nil {
		fmt.Println(" ----------------- Delete Photos Err ------------------- ", err)
	}

	einkaufKb := NewMenu("Lebensmittel", btnStyleGreen, " ************ EINKAUFSLISTE ************** ")
	mes_Params := bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        einkaufKb.Text,
		ReplyMarkup: einkaufKb.KeyBoard,
	}
	_, err = b.SendMessage(ctx, &mes_Params)
	if err != nil {
		fmt.Println("[ERROR Einkauf] [SEND Message failed DUE TO ]", err)
	}
	lmZettel = make(map[string]int)
	line := 2
	num := 0

outerloop:
	for i := 1; i <= len(lebensMittel["Obst"]); i++ {
		if i == (len(lebensMittel["Obst"])/line)+1 {
			break outerloop
		}
		inlineKeyboard := make([]models.InlineKeyboardButton, line)
		inlineKeyboardMarkup.InlineKeyboard = append(inlineKeyboardMarkup.InlineKeyboard, inlineKeyboard)
		for y := 1; y <= line; y++ {
			actI := (((i * line) - line) + y) - 1
			if actI >= len(lebensMittel["Obst"]) {
				break outerloop
			}
			// fmt.Println(y, " ", actI, " ", num, "   ", i)
			inlineKeyboardMarkup.InlineKeyboard[num][y-1].Text = lebensMittel["Obst"][actI]
			inlineKeyboardMarkup.InlineKeyboard[num][y-1].CallbackData = lebensMittel["Obst"][actI]
		}
		num++
	}

	mes_Params = bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        " *************** bitte tippen **************** ",
		ReplyMarkup: inlineKeyboardMarkup,
	}
	_, err = b.SendMessage(ctx, &mes_Params)
	if err != nil {
		fmt.Println("[ERROR Einkauf] [SEND Message failed DUE TO ]", err)
		mes_Params = bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		}
		_, err = b.SendMessage(ctx, &mes_Params)
	}
}
func handler_einkauf(ctx context.Context, b *bot.Bot, update *models.Update) {
	handlerScope = LevelEinkauf
	mes_Params := bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        einkauf.Text,
		ReplyMarkup: einkauf.KeyBoard,
	}
	_, err := b.SendMessage(ctx, &mes_Params)
	if err != nil {
		fmt.Println("[ERROR handler_einkauf] [SEND Message failed DUE TO ]", err)
	}
}

func handler_lebensmittel(ctx context.Context, b *bot.Bot, update *models.Update) {
	handlerScope = LevelLebensmittel
	mes_Params := bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        lbmittel.Text,
		ReplyMarkup: lbmittel.KeyBoard,
	}
	_, err := b.SendMessage(ctx, &mes_Params)
	if err != nil {
		fmt.Println("[ERROR LEBENSMITTEL] [SEND Message failed DUE TO ]", err)
	}
}

func handler_update_inline(ctx context.Context, b *bot.Bot, update *models.Update) {
	callBackData := update.CallbackQuery.Data
	fmt.Println("---------> update inline : ", callBackData)
	for _, l := range lebensMittel["Obst"] {
		if strings.Contains(callBackData, l) {
			lmZettel[l]++
			for x := range inlineKeyboardMarkup.InlineKeyboard {
				for y := range inlineKeyboardMarkup.InlineKeyboard[x] {
					if strings.Contains(inlineKeyboardMarkup.InlineKeyboard[x][y].Text, l) {
						inlineKeyboardMarkup.InlineKeyboard[x][y].Text = l + " (" + strconv.Itoa(lmZettel[l]) + ")"
					}
				}
			}
		}
	}

	edit_Params := bot.EditMessageReplyMarkupParams{
		ChatID:          update.CallbackQuery.Message.Message.Chat.ID,
		MessageID:       update.CallbackQuery.Message.Message.ID,
		InlineMessageID: "hallo",
		ReplyMarkup:     inlineKeyboardMarkup,
	}

	_, err := b.EditMessageReplyMarkup(ctx, &edit_Params)
	if err != nil {
		fmt.Println("Edit ", err)
	}
}

func handler_obst(ctx context.Context, b *bot.Bot, update *models.Update) {
	messageIds = append(messageIds, update.Message.ID)

	//handler_lebensmittel(ctx, b, update)
	handlerScope = LevelObst
	data := update.Message.Text
	fmt.Println("hS: ", handlerScope)
	fmt.Println("msgT", data)
	lmZettel = make(map[string]int)
	line := 2
	num := 0

outerloop:
	for i := 1; i <= len(lebensMittel["Obst"]); i++ {
		if i == (len(lebensMittel["Obst"])/line)+1 {
			break outerloop
		}
		inlineKeyboard := make([]models.InlineKeyboardButton, line)
		inlineKeyboardMarkup.InlineKeyboard = append(inlineKeyboardMarkup.InlineKeyboard, inlineKeyboard)
		for y := 1; y <= line; y++ {
			actI := (((i * line) - line) + y) - 1
			if actI >= len(lebensMittel["Obst"]) {
				break outerloop
			}
			// fmt.Println(y, " ", actI, " ", num, "   ", i)
			inlineKeyboardMarkup.InlineKeyboard[num][y-1].Text = lebensMittel["Obst"][actI]
			inlineKeyboardMarkup.InlineKeyboard[num][y-1].CallbackData = lebensMittel["Obst"][actI]
		}
		num++
	}

	mes_Params := bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        " *************** bitte tippen **************** ",
		ReplyMarkup: inlineKeyboardMarkup,
	}
	_, err := b.SendMessage(ctx, &mes_Params)
	if err != nil {
		fmt.Println("[ERROR Einkauf] [SEND Message failed DUE TO ]", err)
		mes_Params = bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		}
		_, err = b.SendMessage(ctx, &mes_Params)
	}
	fmt.Println(lmZettel)
}

func handler_todo(ctx context.Context, b *bot.Bot, update *models.Update) {
	messageIds = append(messageIds, update.Message.ID)
	msgText := update.Message.Text
	fmt.Println("hS: ", handlerScope)
	fmt.Println("handler_todo msgT", msgText)
	// if handlerScope == "todo Einkauf" {

	// }
	handlerScope = LevelTodo

	mes_Params := bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        todo.Text,
		ReplyMarkup: todo.KeyBoard,
	}
	_, err := b.SendMessage(ctx, &mes_Params)
	if err != nil {
		fmt.Println("[ERROR TODO] [SEND Message failed DUE TO ]", err)
	}
}

func handler_witz(ctx context.Context, b *bot.Bot, update *models.Update) {
	messageIds = append(messageIds, update.Message.ID)
	fmt.Println("handler_cinema msg_text: ", update.Message.Text)
	//if update.Message.Text == "/witz" {
	handlerScope = LevelWitz
	//}

	witz := witz.GetWitz()

	mes_Params := bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        witz.Text,
		ReplyMarkup: witzMarkUp.KeyBoard,
	}
	_, err := b.SendMessage(ctx, &mes_Params)
	if err != nil {
		fmt.Println("[ERROR Witz] [SEND Message failed DUE TO ]", err)
	}
}

func handlerFilm(ctx context.Context, b *bot.Bot, update *models.Update) {

	fmt.Println("handlerFilm msg_text: ", update.CallbackQuery.Data)
	var msg string
	var titel string
	var poster string
	var description string
	for i := range kinos {
		if strings.Contains(update.CallbackQuery.Data, kinos[i].Title) {
			titel = kinos[i].Title
			poster = kinos[i].Poster
			description = kinos[i].Description
			msg = "Zeiten: \n"
			for x := range kinos[i].Showtimes {
				msg = msg + "\t" + kinos[i].Showtimes[x].Datetime + "\n"
			}
			break
		}
	}
	b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
		ShowAlert:       false,
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   "Sie haben gewählt: " + titel + "\n" + msg,
	})
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Text:   description,
	})

	foto_Params := bot.SendPhotoParams{
		ChatID: update.CallbackQuery.Message.Message.Chat.ID,
		Photo:  &models.InputFileString{Data: poster},
		//MessageEffectID: Conf_tok["SUCCESS_EFFECT"],
	}
	_, err := b.SendPhoto(ctx, &foto_Params)
	if err != nil {
		fmt.Println("[ERROR CINESTAR] [SENDING MESSAGE FAILED DUE TO ]", err)
	}
}

func handler_cinema(ctx context.Context, b *bot.Bot, update *models.Update) {
	//var foto_Params bot.SendPhotoParams
	fmt.Println("handler_cinema msg_text: ", update.Message.Text)
	if update.Message.Text == "/kino" || update.Message.Text == "Kino" {
		handlerScope = LevelKino
	}
	if update.Message.Text == "heute" || update.Message.Text == "HEUTE" {
		handlerScope = LevelKinoHeute
	}
	if update.Message.Text == "woche" || update.Message.Text == "WOCHE" {
		handlerScope = LevelKinoWoche
	}

	switch handlerScope {
	case LevelKino:
		filterContainer(ctx, b, update)
		return
	case LevelKinoHeute:
		filterHeute(ctx, b, update)
		return
	case LevelKinoWoche:
		filterWoche(ctx, b, update)
		return
	// case "cinema vorschau":
	// 	filterVorschau(ctx, b, update)
	// 	return
	// case "cinema events":
	// 	filterEvents(ctx, b, update)
	// 	return
	// case "cinema vorverkauf":
	// 	filterVorverkauf(ctx, b, update)
	// 	return
	default:
		kinos := kino.GetKino(handlerScope)
		kinoGesamt := 0
		for i := 0; i < len(kinos); i++ {
			kinoGesamt = i + 1
		}
		mes_Params := bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf(" -- gesamt: %d Filme --", kinoGesamt),
		}
		_, err := b.SendMessage(ctx, &mes_Params)
		if err != nil {
			fmt.Println("[ERROR CINESTAR] [SEND Message failed DUE TO ]", err)
		}
		// fmt.Printf("--> Cinema Text %d gesendet!\n", kinoGesamt)

		fmt.Println()
		return
	}
	// handlerScope = "cinema"

}
func filterContainer(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("filterContainer msg_text: ", update.Message.Text)
	mes_Params := bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        kinoMarkUp.Text,
		ReplyMarkup: kinoMarkUp.KeyBoard,
	}
	_, err := b.SendMessage(ctx, &mes_Params)
	if err != nil {
		fmt.Println("[ERROR CINESTAR] [SEND Message failed DUE TO ]", err)
	}
	fmt.Println("Scope kino : ", handlerScope)
}
func filterHeute(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("filterHeute ---------- > handlerScope: ", handlerScope)
	handlerScope = LevelKinoHeute
	kinos = kino.GetKino(LevelKinoHeute)

	if len(kinos) == 0 {
		mes_Params := bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Heute keine Fime mehr",
		}
		_, err := b.SendMessage(ctx, &mes_Params)
		if err != nil {
			fmt.Println("keine Filme mehr ", err)
		}
		return
	}

	inlineKeyboardMarkup := &models.InlineKeyboardMarkup{}

	for i := 0; i < len(kinos); i++ {
		inlineKeyboard := make([]models.InlineKeyboardButton, 1)
		inlineKeyboardMarkup.InlineKeyboard = append(inlineKeyboardMarkup.InlineKeyboard, inlineKeyboard)
		inlineKeyboardMarkup.InlineKeyboard[i][0].Text = kinos[i].Title
		inlineKeyboardMarkup.InlineKeyboard[i][0].CallbackData = "Filme Heute " + kinos[i].Title
	}
	mes_Params := bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Heute im Kino",
		ReplyMarkup: inlineKeyboardMarkup,
	}
	_, err := b.SendMessage(ctx, &mes_Params)
	if err != nil {
		fmt.Println("[ERROR CINESTAR] [SEND Message failed DUE TO ]", err)
		mes_Params = bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		}
		_, err = b.SendMessage(ctx, &mes_Params)
		return
	}
}
func filterWoche(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("Scope kino : ", handlerScope)
	handlerScope = LevelKinoWoche
	kinos = kino.GetKino(LevelKinoWoche)
	if len(kinos) == 0 {
		mes_Params := bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Diese Woche keine Fime mehr",
		}
		_, err := b.SendMessage(ctx, &mes_Params)
		if err != nil {
			fmt.Println("keine Filme mehr ", err)
		}
		return
	}

	kinoGesamt := 0
	inlineKeyboardMarkup := &models.InlineKeyboardMarkup{}

	for i := 0; i < len(kinos); i++ {
		inlineKeyboard := make([]models.InlineKeyboardButton, 1)
		inlineKeyboardMarkup.InlineKeyboard = append(inlineKeyboardMarkup.InlineKeyboard, inlineKeyboard)
		// fmt.Println(kinos[i].Title, " - ", len(kinos[i].Title))
		inlBtnText := kinos[i].Title
		if len(kinos[i].Title) >= 64 {
			inlBtnText = "Komischer Name"
			mes_Params := bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   kinos[i].Title,
			}
			_, err := b.SendMessage(ctx, &mes_Params)
			if err != nil {
				fmt.Println(err)
			}
		}
		inlineKeyboardMarkup.InlineKeyboard[i][0].Text = inlBtnText
		inlineKeyboardMarkup.InlineKeyboard[i][0].CallbackData = "Filme Woche " + inlBtnText
		// fmt.Println("-->  gesendet!   ", inlineKeyboardMarkup.InlineKeyboard[i][0].Text, "  ", i)
		kinoGesamt = i + 1
	}
	fmt.Println(kinoGesamt)
	mes_Params := bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Diese Woche im Kino",
		ReplyMarkup: inlineKeyboardMarkup,
	}
	_, err := b.SendMessage(ctx, &mes_Params)
	if err != nil {
		fmt.Println("[ERROR CINESTAR] [SEND Message failed DUE TO ]", err)
		mes_Params = bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   err.Error(),
		}
		_, err = b.SendMessage(ctx, &mes_Params)
		return
	}
}
func filterVorschau(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("Scope kino : ", handlerScope)
}
func filterVorverkauf(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("Scope kino : ", handlerScope)
}
func filterEvents(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("Scope kino : ", handlerScope)
}

func handler_eurojackpot(ctx context.Context, b *bot.Bot, update *models.Update) {
	text := update.Message.Text
	handlerScope = LevelJP
	initButton := "init"
	for x := 1; x <= 12; x++ {
		var m time.Month
		m = time.Month(x)
		if text == m.String() {
			initButton = "month"
			scopeDat := "01." + checkDayMonthLen(strconv.Itoa(x)) + ".2026"
			JackpotDatumScope = Scope(scopeDat)
			fmt.Println(monthOnly, " - ", len(jpDates))
			if (monthOnly || (!monthOnly && len(jpDates) > 0)) && monthList {
				monthOnly = false
				monthList = false
				initButton = "init"
			} else if (!monthOnly && len(jpDates) >= 0) && !monthList {
				jpDates = nil
				monthOnly = true
				monthList = true
				fmt.Println("+++ nur Monate +++")
			}
			fmt.Println(" nach if ----> ", monthOnly, " - ", len(jpDates))
		}
	}

	if text == "<< <" {
		JackpotDatumScope = newMonth("-", string(JackpotDatumScope))
	}
	if text == "> >>" {
		JackpotDatumScope = newMonth("+", string(JackpotDatumScope))
		actScope := time.Now()
		t, err := time.Parse("02.01.2006", string(JackpotDatumScope))
		if err != nil {
			fmt.Println(err)
		}
		if t.After(actScope) {
			fmt.Println(JackpotDatumScope, " > ", actScope)
			dateButtons := getDateButtons2(handlerScope, "01.07.2026", initButton)
			message_Params := bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        " Monat liegt in der Zukunft.\nWählen Sie das Datum ! ",
				ReplyMarkup: dateButtons,
			}
			_, err := b.SendMessage(ctx, &message_Params)
			if err != nil {
				fmt.Println("[ERROR] [handler_lotto] : ", err)
			}
			JackpotDatumScope = Scope(actScope.Format("02.01.2006"))
			return
		}
	}
	fmt.Println(" update.Message: ", text)
	if text == "<< <<" {
		JackpotDatumScope = newYear("-", string(JackpotDatumScope))
		fmt.Println(" ----- --------- ------------- ", JackpotDatumScope)
		// JackpotDatumScope = "01.06.2025"
	}
	if text == ">> >>" {
		JackpotDatumScope = newYear("+", string(JackpotDatumScope))
		// JackpotDatumScope = "01.06.2025"
	}

	jpDates = getDatesJP(string(JackpotDatumScope))
	fmt.Println(" mO: ", monthOnly, " ------ >>>> initButton : ", initButton, "  len jpDates: ", len(jpDates), "  ", string(JackpotDatumScope))
	dateButtons := getDateButtons2(handlerScope, string(JackpotDatumScope), initButton)
	message_Params := bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        " Wählen Sie das Datum ! ",
		ReplyMarkup: dateButtons,
	}
	_, err := b.SendMessage(ctx, &message_Params)
	if err != nil {
		fmt.Println("[ERROR] [handler_lotto] : ", err)
	}
}

func handler_lotto(ctx context.Context, b *bot.Bot, update *models.Update) {
	fmt.Println("handler_lotto String: ", update.Message.Text)
	text := update.Message.Text
	handlerScope = Level6aus49
	initButton := "init"
	for x := 1; x <= 12; x++ {
		var m time.Month
		m = time.Month(x)
		if text == m.String() {
			initButton = "month"
			scopeDat := "01." + checkDayMonthLen(strconv.Itoa(x)) + ".2026"
			LottoDatumScope = Scope(scopeDat)
			fmt.Println(monthOnly, " - ", len(jpDates))
			if (monthOnly || (!monthOnly && len(lDates) > 0)) && monthList {
				monthOnly = false
				monthList = false
				initButton = "init"
			} else if (!monthOnly && len(lDates) >= 0) && !monthList {
				lDates = nil
				monthOnly = true
				monthList = true
				fmt.Println("+++ nur Monate +++")
			}
			fmt.Println(" nach if ----> ", monthOnly, " - ", len(lDates))
		}
	}
	if text == "<< <" {
		LottoDatumScope = newMonth("-", string(LottoDatumScope))
	}
	if text == "> >>" {
		LottoDatumScope = newMonth("+", string(LottoDatumScope))
		actScope := time.Now()
		t, err := time.Parse("02.01.2006", string(LottoDatumScope))
		if err != nil {
			fmt.Println(err)
		}
		if t.After(actScope) {
			fmt.Println(LottoDatumScope, " > ", actScope)
			dateButtons := getDateButtons2(handlerScope, "01.07.2026", initButton)
			message_Params := bot.SendMessageParams{
				ChatID:      update.Message.Chat.ID,
				Text:        " Monat liegt in der Zukunft.\nWählen Sie das Datum ! ",
				ReplyMarkup: dateButtons,
			}
			_, err := b.SendMessage(ctx, &message_Params)
			if err != nil {
				fmt.Println("[ERROR] [handler_lotto] : ", err)
			}
			LottoDatumScope = Scope(actScope.Format("02.01.2006"))
			return
		}
	}
	fmt.Println(" update.Message: ", text)
	if text == "<< <<" {
		LottoDatumScope = newYear("-", string(LottoDatumScope))
		fmt.Println(" ----- --------- ------------- ", LottoDatumScope)
		// JackpotDatumScope = "01.06.2025"
	}
	if text == ">> >>" {
		LottoDatumScope = newYear("+", string(LottoDatumScope))
		// JackpotDatumScope = "01.06.2025"
	}

	lDates = getDatesJP(string(LottoDatumScope))
	fmt.Println(" mO: ", monthOnly, " ------ >>>> initButton : ", initButton, "  len jpDates: ", len(lDates), "  ", string(LottoDatumScope))
	dateButtons := getDateButtons2(handlerScope, string(LottoDatumScope), initButton)
	//dateButtons := getDateButtons(lDates)
	message_Params := bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        " Wählen Sie das Datum ! ",
		ReplyMarkup: dateButtons,
	}
	_, err := b.SendMessage(ctx, &message_Params)
	if err != nil {
		fmt.Println("[ERROR] [handler_lotto] : ", err)
	}
}

func handler_lotto_default(ctx context.Context, b *bot.Bot, update *models.Update) {
	handlerScope = LevelLotto

	message_Params := bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        lottoMarkUp.Text,
		ReplyMarkup: lottoMarkUp.KeyBoard,
	}
	_, err := b.SendMessage(ctx, &message_Params)
	if err != nil {
		fmt.Println("[ERROR] [handler_lotto_defoault] : ", err)
	}
}

func handler_botschaft(ctx context.Context, b *bot.Bot, update *models.Update) {
	// msg := botschaft.Get_Botschaft(Conf_tok["BOTSCHAFT"])
	// handlerScope = "botschaft"
	// err := os.WriteFile("screenshot.jpeg", msg.Screenshot, 0644)
	// if err != nil {
	// 	fmt.Printf("Error writing screenshot: %v ", err)
	// }

	// file, _ := os.Open("screenshot1.jpeg")
	// defer file.Close()

	// message_Params := bot.SendPhotoParams{
	// 	ChatID: update.Message.Chat.ID,
	// 	Photo: &models.InputFileUpload{
	// 		// Filename: "screenshot.jpeg",
	// 		Data: file,
	// 	},
	// }
	// _, err = b.SendPhoto(ctx, &message_Params)
	// if err != nil {
	// 	fmt.Println("[ERROR] [SENDING MESSAGE FAILED DUE TO ]", err)
	// }
}

func newYear(op string, date string) Scope {
	current := time.Now().Format("02.01.2006")
	t, err := time.Parse("02.01.2006", date)
	if err != nil {
		fmt.Println("[new month parsing date fail:] ", err)
	}

	switch op {
	case "-":
		t = t.AddDate(-1, 0, 0)
		if current < t.Format("02.01.2006") {
			fmt.Println(" Neues Jahr ----------<> ", t.Format("02.01.2006"))
			return Scope(current)
		}
		newDate := t.Format("02.01.2006")
		return Scope(newDate)
	case "+":
		t = t.AddDate(1, 0, 0)
		if current < t.Format("02.01.2006") {
			return Scope(current)
		}
		newDate := t.Format("02.01.2006")
		return Scope(newDate)
	}

	return ""
}

func newMonth(op string, date string) Scope {
	fmt.Println(" ----- date: ", date, "  ----------")
	current := time.Now().Format("02.01.2006")
	t, err := time.Parse("02.01.2006", date)
	if err != nil {
		fmt.Println("[new month parsing date fail:] ", err)
	}
	actMonth := t.Month().String()
	actYear := t.Year()
	fmt.Println(actYear, " - ", actMonth)
	newMonth := 0
	switch op {
	case "-":
		for x := 1; x <= 12; x++ {
			var m time.Month
			m = time.Month(x)
			fmt.Println("\t", m.String())
			if actMonth == m.String() {
				newMonth = (x - 1)
				if newMonth == 0 {
					newMonth = 12
					actYear--
				}
				scopeDat := "01." + checkDayMonthLen(strconv.Itoa(newMonth)) + "." + strconv.Itoa(actYear)
				if scopeDat > current {
					return Scope(date)
				}
				JackpotDatumScope = Scope(scopeDat)
				fmt.Println("\t\n", JackpotDatumScope)
				return JackpotDatumScope
			}
		}
	case "+":

		for x := 1; x <= 12; x++ {
			var m time.Month
			m = time.Month(x)
			if actMonth == m.String() {
				newMonth = (x + 1)
				if newMonth == 13 {
					newMonth = 1
					actYear++
				}
				fmt.Println("\t", m.String(), "  x: ", x, " newMonth: ", newMonth, " actYear: ", actYear)
				scopeDat := "01." + checkDayMonthLen(strconv.Itoa(newMonth)) + "." + strconv.Itoa(actYear)
				JackpotDatumScope = Scope(scopeDat)
				fmt.Println("\tJackpotDatumScope: ", JackpotDatumScope)
				return JackpotDatumScope
			}
		}
	}

	return ""
}

// Die Funktion lädt KeyBoardButtons mit Jahr, Monat und Spieltagen - in LevelJP
func initialButton(d string) *models.ReplyKeyboardMarkup {
	kb := &models.ReplyKeyboardMarkup{}
	jahrBefore := "<< <<"
	jahrAfter := ">> >>"
	monthBefore := "<< <"
	monthAfter := "> >>"
	datum, _ := time.Parse("02.01.2006", d)
	year := datum.Year()
	month := datum.Month()
	btnsProLine := 3
	btnLine := 0
	kbLine := make([]models.KeyboardButton, btnsProLine)
	kb.Keyboard = append(kb.Keyboard, kbLine)
	lineContent := make([]models.KeyboardButton, btnsProLine)
	kb.Keyboard[btnLine] = append(kb.Keyboard[btnLine], lineContent...)
	kb.Keyboard[btnLine][0].Text = jahrBefore
	kb.Keyboard[btnLine][1].Text = strconv.Itoa(year)
	kb.Keyboard[btnLine][2].Text = jahrAfter
	btnLine++
	kbLine = make([]models.KeyboardButton, btnsProLine)
	kb.Keyboard = append(kb.Keyboard, kbLine)
	kb.Keyboard[btnLine][0].Text = monthBefore
	kb.Keyboard[btnLine][1].Text = month.String() //strconv.Itoa(int(month))
	kb.Keyboard[btnLine][2].Text = monthAfter
	fmt.Println("+++ len dates jp: ", len(jpDates))
	if len(jpDates) > 0 {
		if len(jpDates) >= 4 {
			for i := 0; i < len(jpDates); i++ {
				maxDateBtnProLine := 3
				if (len(jpDates) - i) < maxDateBtnProLine {
					maxDateBtnProLine = (len(jpDates) - i)
				}
				kbLine = make([]models.KeyboardButton, maxDateBtnProLine)
				kb.Keyboard = append(kb.Keyboard, kbLine)
				btnLine++
				for y := 0; y < maxDateBtnProLine; y++ {
					// kb.Keyboard[len(kb.Keyboard)-1][y].Text = jpDates[i]
					kb.Keyboard[btnLine][y].Text = jpDates[i]
					i++
				}
				fmt.Println(" +++ ", i)
				i--
			}
		}
		if len(jpDates) < 4 {
			for i := 0; i < len(jpDates); i++ {
				maxDateBtnProLine := 2
				if (len(jpDates) - i) < maxDateBtnProLine {
					maxDateBtnProLine = (len(jpDates) - i)
				}
				kbLine = make([]models.KeyboardButton, maxDateBtnProLine)
				kb.Keyboard = append(kb.Keyboard, kbLine)
				btnLine++
				for y := 0; y < maxDateBtnProLine; y++ {
					kb.Keyboard[btnLine][y].Text = jpDates[i]
					i++
				}
				i--
			}
		}
	}
	kbLine = make([]models.KeyboardButton, btnsProLine)
	kb.Keyboard = append(kb.Keyboard, kbLine)
	kb.Keyboard[len(kb.Keyboard)-1][0].Text = "Menue"
	kb.Keyboard[len(kb.Keyboard)-1][1].Text = "Lotto"
	// kb.InputFieldPlaceholder = checkDayMonthLen(kb.Keyboard[2][1].Text) + "." + checkDayMonthLen(kb.Keyboard[1][1].Text) + "." + kb.Keyboard[0][1].Text
	return kb
}

func initialButton6aus49(d string) *models.ReplyKeyboardMarkup {
	kb := &models.ReplyKeyboardMarkup{}
	jahrBefore := "<< <<"
	jahrAfter := ">> >>"
	monthBefore := "<< <"
	monthAfter := "> >>"
	datum, _ := time.Parse("02.01.2006", d)
	year := datum.Year()
	month := datum.Month()
	btnsProLine := 3
	btnLine := 0
	kbLine := make([]models.KeyboardButton, btnsProLine)
	kb.Keyboard = append(kb.Keyboard, kbLine)
	lineContent := make([]models.KeyboardButton, btnsProLine)
	kb.Keyboard[btnLine] = append(kb.Keyboard[btnLine], lineContent...)
	kb.Keyboard[btnLine][0].Text = jahrBefore
	kb.Keyboard[btnLine][1].Text = strconv.Itoa(year)
	kb.Keyboard[btnLine][2].Text = jahrAfter
	btnLine++
	kbLine = make([]models.KeyboardButton, btnsProLine)
	kb.Keyboard = append(kb.Keyboard, kbLine)
	kb.Keyboard[btnLine][0].Text = monthBefore
	kb.Keyboard[btnLine][1].Text = month.String() //strconv.Itoa(int(month))
	kb.Keyboard[btnLine][2].Text = monthAfter
	fmt.Println("+++ len dates jp: ", len(lDates))
	if len(lDates) > 0 {
		if len(lDates) >= 4 {
			for i := 0; i < len(lDates); i++ {
				maxDateBtnProLine := 3
				if (len(lDates) - i) < maxDateBtnProLine {
					maxDateBtnProLine = (len(lDates) - i)
				}
				kbLine = make([]models.KeyboardButton, maxDateBtnProLine)
				kb.Keyboard = append(kb.Keyboard, kbLine)
				btnLine++
				for y := 0; y < maxDateBtnProLine; y++ {
					// kb.Keyboard[len(kb.Keyboard)-1][y].Text = jpDates[i]
					kb.Keyboard[btnLine][y].Text = lDates[i]
					i++
				}
				fmt.Println(" +++ ", i)
				i--
			}
		}
		if len(lDates) < 4 {
			for i := 0; i < len(lDates); i++ {
				maxDateBtnProLine := 2
				if (len(lDates) - i) < maxDateBtnProLine {
					maxDateBtnProLine = (len(lDates) - i)
				}
				kbLine = make([]models.KeyboardButton, maxDateBtnProLine)
				kb.Keyboard = append(kb.Keyboard, kbLine)
				btnLine++
				for y := 0; y < maxDateBtnProLine; y++ {
					kb.Keyboard[btnLine][y].Text = lDates[i]
					i++
				}
				i--
			}
		}
	}
	kbLine = make([]models.KeyboardButton, btnsProLine)
	kb.Keyboard = append(kb.Keyboard, kbLine)
	kb.Keyboard[len(kb.Keyboard)-1][0].Text = "Menue"
	kb.Keyboard[len(kb.Keyboard)-1][1].Text = "Lotto"
	// kb.InputFieldPlaceholder = checkDayMonthLen(kb.Keyboard[2][1].Text) + "." + checkDayMonthLen(kb.Keyboard[1][1].Text) + "." + kb.Keyboard[0][1].Text
	return kb
}

// Die Funktion lädt KeyBoardButtons mit Jahr und Monatsnamen - in LevelJP
func monthButton(d string) *models.ReplyKeyboardMarkup {
	kb := &models.ReplyKeyboardMarkup{}
	jahrBefore := "<< <<"
	jahrAfter := ">> >>"
	datum, _ := time.Parse("02.01.2006", d)
	year := datum.Year()
	btnsProLine := 3
	btnLine := 0
	kbLine := make([]models.KeyboardButton, btnsProLine)
	kb.Keyboard = append(kb.Keyboard, kbLine)
	lineContent := make([]models.KeyboardButton, btnsProLine)
	kb.Keyboard[btnLine] = append(kb.Keyboard[btnLine], lineContent...)
	kb.Keyboard[0][0].Text = jahrBefore
	kb.Keyboard[0][1].Text = strconv.Itoa(year)
	kb.Keyboard[0][2].Text = jahrAfter
	for i := 1; i <= 12; i++ {
		kbLine = make([]models.KeyboardButton, btnsProLine)
		kb.Keyboard = append(kb.Keyboard, kbLine)
		for x := range 3 {
			var m time.Month
			m = time.Month(i)
			kb.Keyboard[len(kb.Keyboard)-1][x].Text = m.String()
			i++
			if i == 12 {
				m = time.Month(i)
				kb.Keyboard[len(kb.Keyboard)-1][x].Text = m.String()
			}
		}
		i--
	}
	kbLine = make([]models.KeyboardButton, btnsProLine)
	kb.Keyboard = append(kb.Keyboard, kbLine)
	kb.Keyboard[len(kb.Keyboard)-1][0].Text = "Menue"
	kb.Keyboard[len(kb.Keyboard)-1][1].Text = "Lotto"
	// kb.InputFieldPlaceholder = checkDayMonthLen(kb.Keyboard[2][1].Text) + "." + checkDayMonthLen(kb.Keyboard[1][1].Text) + "." + kb.Keyboard[0][1].Text
	return kb

}
func getButtonJackpot(d string, initButton string) *models.ReplyKeyboardMarkup {
	fmt.Println("btn: ", initButton)
	// kb := &models.ReplyKeyboardMarkup{}
	switch initButton {
	case "init":
		return initialButton(d)
	case "month":
		return monthButton(d)
	}
	return nil
}
func getButtonLotto(d string, initButton string) *models.ReplyKeyboardMarkup {
	fmt.Println("btn: ", initButton)
	// kb := &models.ReplyKeyboardMarkup{}
	switch initButton {
	case "init":
		return initialButton6aus49(d)
	case "month":
		return monthButton(d)
	}
	return nil
}
func getDateButtons2(s int, date string, initButton string) *models.ReplyKeyboardMarkup {
	// kb := &models.ReplyKeyboardMarkup{}
	switch s {
	case Level6aus49:
		return getButtonLotto(date, initButton)
	case LevelJP:
		return getButtonJackpot(date, initButton)
	}
	return nil
}
func getDateButtons(dates []string) *models.ReplyKeyboardMarkup {
	fmt.Println("getting Buttons ... ", dates)
	btnCount := 5
	btnLine := 0
	kb := &models.ReplyKeyboardMarkup{}
	kbLine := make([]models.KeyboardButton, btnCount)
	kb.Keyboard = append(kb.Keyboard, kbLine)
	lineContent := make([]models.KeyboardButton, btnCount)
	kb.Keyboard[btnLine] = append(kb.Keyboard[btnLine], lineContent...)
	// fmt.Println("Länge: ", len(kb.Keyboard[0]), " ", len(euroJackpotDates[jpotDat]))
	idx := 0
	for x := range dates {
		fmt.Println(dates[x])
		if idx >= btnCount {
			fmt.Println(dates[x])
			idx = 0
			btnLine++
			kbLine := make([]models.KeyboardButton, btnCount)
			kb.Keyboard = append(kb.Keyboard, kbLine)
			lineContent := make([]models.KeyboardButton, btnCount)
			kb.Keyboard[btnLine] = append(kb.Keyboard[btnLine], lineContent...)
		}

		kb.Keyboard[btnLine][idx].Text = dates[x]
		idx++
	}
	kbLine = make([]models.KeyboardButton, btnCount)
	kb.Keyboard = append(kb.Keyboard, kbLine)
	lineContent = make([]models.KeyboardButton, btnCount)
	lineContent[0].Text = "vorher"
	lineContent[1].Text = "ZURÜCK"
	lineContent[2].Text = "nachher"
	btnLine++
	kb.Keyboard[btnLine] = append(kb.Keyboard[btnLine], lineContent...)
	// for i := range kb.Keyboard {
	// 	for y := range kb.Keyboard[i] {
	// 		fmt.Print(kb.Keyboard[i][y].Text, "  ")
	// 	}
	// 	fmt.Println()
	// }

	return kb
}

// Die Funktion gibt eine Liste von Spieltgen im aktuell gewählten Monat, im jeweiligen Spiel zurück.
func getDatesJP(today string) []string {
	var dates = make([]string, 0)
	var monate []time.Month
	now, _ := time.Parse("02.01.2006", time.Now().Format("02.01.2006"))
	t, _ := time.Parse("02.01.2006", today)
	fmt.Println(t.Month().String(), "   ", handlerScope)
	year := t.Year()    // Beispieljahr
	month1 := t.Month() // Beispielmonat
	// month2 := t.Month() - 1
	monate = append(monate, month1)
	// Erster Tag des Monats
	var spielTage []time.Time
	for i := range monate {
		firstDay := time.Date(year, monate[i], 1, 0, 0, 0, 0, time.UTC)
		// Letzter Tag des Monats
		lastDay := firstDay.AddDate(0, 1, -1)
		// Iteriere durch alle Tage des Monats
		for d := firstDay; d.Before(lastDay.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
			if d.After(now) {
				// fmt.Printf("%s after %s ", d.Format("02.01.2006"), now.Format("02.01.2006"))
				break
			}
			switch handlerScope {
			case LevelJP:
				if d.Weekday() == time.Tuesday || d.Weekday() == time.Friday {
					onlyFr, _ := time.Parse("02.01.2006", "29.03.2022")
					if d.Before(onlyFr) && d.Weekday() == time.Tuesday {
						fmt.Println("Dienstag !!!!!!!!!!!!!!!!!")
						break
					}
					spielTage = append(spielTage, d)
				}
			case Level6aus49:
				if d.Weekday() == time.Wednesday || d.Weekday() == time.Saturday {
					spielTage = append(spielTage, d)
				}
			}
		}
	}

	// Ausgabe der Dienstage
	for _, spielTag := range spielTage {
		fmt.Println(spielTag.Format("02.01.2006"))
		dates = append(dates, spielTag.Format("02.01.2006"))
	}
	return dates
}

func showJackpotScreenshot(update *models.Update) {
	file, err := os.Open(SCREENSHOTNAMEJP)
	if err != nil {
		fmt.Println("[ERROR] Datei jackpot.png konnte nicht geöffnet werden -> ", err.Error())
	}
	defer file.Close()
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	// Formularfeld "photo"
	part, err := writer.CreateFormFile("photo", SCREENSHOTNAMEJP)
	if err != nil {
		fmt.Println("[ERROR] Kein FormFile konnte erstellt werden -> ", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("[ERROR] jackpot.png konnte nicht kopiert werden -> ", err)
	}
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto?chat_id=%d", getToken(), update.Message.Chat.ID), &requestBody)
	if err != nil {
		fmt.Println("[ERROR] NewRequest Fehler : ", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERROR] Request fehlgeschlagen : ", err)
	}
	defer resp.Body.Close()

}

func showLottoScreenshot(update *models.Update) {
	file, err := os.Open(SCREENSHOTLOTTO)
	if err != nil {
		fmt.Println("[ERROR] Datei jackpot.png konnte nicht geöffnet werden -> ", err.Error())
	}
	defer file.Close()
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	// Formularfeld "photo"
	part, err := writer.CreateFormFile("photo", SCREENSHOTLOTTO)
	if err != nil {
		fmt.Println("[ERROR] Kein FormFile konnte erstellt werden -> ", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("[ERROR] jackpot.png konnte nicht kopiert werden -> ", err)
	}
	writer.Close()
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto?chat_id=%d", getToken(), update.Message.Chat.ID), &requestBody)
	if err != nil {
		fmt.Println("[ERROR] NewRequest Fehler : ", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERROR] Request fehlgeschlagen : ", err)
	}
	defer resp.Body.Close()

}

func checkDayMonthLen(s string) string {
	if len(s) == 1 {
		return "0" + s
	}
	return s
}
