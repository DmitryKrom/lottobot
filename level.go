package main

import (
	"context"
	"fmt"
	"lottobot/lotto"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	Januar    = "January"
	February  = "February"
	March     = "March"
	April     = "April"
	May       = "May"
	June      = "June"
	July      = "July"
	August    = "August"
	September = "September"
	Oktober   = "Oktober"
	November  = "November"
	December  = "December"
)

const (
	LevelStart = iota
	LevelMenu
	LevelLotto
	LevelJP
	Level6aus49
	LevelWitz
	LevelKino
	LevelKinoHeute
	LevelKinoWoche
	LevelTodo
	LevelEinkauf
	LevelLebensmittel
	LevelObst
	LevelFleisch
	LevelWäsche
	LevelTermin
)

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {

	switch handlerScope {
	case LevelMenu:
		eingabe := update.Message.Text
		fmt.Println("Level: ", handlerScope)
		l, swiched := checkLevelSwitch(eingabe, ctx, b, update)
		if swiched {
			sendFaultMessage(LevelMenu, fmt.Sprintf("Level %d sitched to Level %d ", LevelMenu, l), ctx, b, update)
			return
		}

	case LevelLebensmittel:
		eingabe := update.CallbackQuery.Data
		_, ok := lebensMittel[eingabe]
		if !ok {
			sendFaultMessage(LevelLebensmittel, "Ich kenne diese Lebensmittelkategorie nicht!!!", ctx, b, update)
		}
		handler_update_inline(ctx, b, update)
		return
	case LevelJP:
		current := time.Now().Format("02.01.2006")
		LottoDatumScope = Scope(current)
		eingabe := update.Message.Text
		fmt.Println("Level: ", handlerScope)
		l, swiched := checkLevelSwitch(eingabe, ctx, b, update)
		if swiched {
			fmt.Println(LevelJP, " - switched ")
			sendFaultMessage(LevelJP, fmt.Sprintf("Level %d sitched to Level %d ", LevelJP, l), ctx, b, update)
			return
		}
		switch eingabe {
		case "<< <":
			handler_eurojackpot(ctx, b, update)
			return
		case "> >>":
			handler_eurojackpot(ctx, b, update)
			return
		case "<< <<":
			handler_eurojackpot(ctx, b, update)
			return
		case ">> >>":
			handler_eurojackpot(ctx, b, update)
			return
		case Januar:
			handler_eurojackpot(ctx, b, update)
			return
		case February:
			handler_eurojackpot(ctx, b, update)
			return
		case March:
			handler_eurojackpot(ctx, b, update)
			return
		case April:
			handler_eurojackpot(ctx, b, update)
			return
		case May:
			fmt.Println("case May !!!")
			handler_eurojackpot(ctx, b, update)
			return
		case June:
			handler_eurojackpot(ctx, b, update)
			return
		case July:
			handler_eurojackpot(ctx, b, update)
			return
		case August:
			handler_eurojackpot(ctx, b, update)
			return
		case September:
			handler_eurojackpot(ctx, b, update)
			return
		case Oktober:
			handler_eurojackpot(ctx, b, update)
			return
		case November:
			handler_eurojackpot(ctx, b, update)
			return
		case December:
			handler_eurojackpot(ctx, b, update)
			return
		}
		if err := check_datum(eingabe); err != nil {
			sendFaultMessage(LevelJP, err.Error(), ctx, b, update)
			return
		}
		fmt.Println(" -------- ----- ", swiched)
		gw, err := lotto.Get_lotto(LevelJP, eingabe)

		gewinn_Text := build_message(gw, eingabe)

		message_Params := bot.SendMessageParams{
			ChatID:          update.Message.Chat.ID,
			Text:            gewinn_Text + "\n",
			MessageEffectID: Conf_tok["SUCCESS_EFFECT"],
		}
		_, err = b.SendMessage(ctx, &message_Params)
		if err != nil {
			fmt.Println("[ERROR] [SENDING MESSAGE FAILED DUE TO ]", err)
		}

	case Level6aus49:
		current := time.Now().Format("02.01.2006")
		JackpotDatumScope = Scope(current)
		eingabe := update.Message.Text
		fmt.Println("Level: ", handlerScope)
		l, swiched := checkLevelSwitch(eingabe, ctx, b, update)
		if swiched {
			fmt.Println(Level6aus49, " - switched ")
			sendFaultMessage(Level6aus49, fmt.Sprintf("Level %d sitched to Level %d ", LevelJP, l), ctx, b, update)
			return
		}
		switch eingabe {
		case "<< <":
			handler_lotto(ctx, b, update)
			return
		case "> >>":
			handler_lotto(ctx, b, update)
			return
		case "<< <<":
			handler_lotto(ctx, b, update)
			return
		case ">> >>":
			handler_lotto(ctx, b, update)
			return
		case Januar:
			handler_lotto(ctx, b, update)
			return
		case February:
			handler_lotto(ctx, b, update)
			return
		case March:
			handler_lotto(ctx, b, update)
			return
		case April:
			handler_lotto(ctx, b, update)
			return
		case May:
			handler_lotto(ctx, b, update)
			return
		case June:
			handler_lotto(ctx, b, update)
			return
		case July:
			handler_lotto(ctx, b, update)
			return
		case August:
			handler_lotto(ctx, b, update)
			return
		case September:
			handler_lotto(ctx, b, update)
			return
		case Oktober:
			handler_lotto(ctx, b, update)
			return
		case November:
			handler_lotto(ctx, b, update)
			return
		case December:
			handler_lotto(ctx, b, update)
			return
		}
		if err := check_datum(eingabe); err != nil {
			sendFaultMessage(Level6aus49, err.Error(), ctx, b, update)
			return
		}
		fmt.Println(" -------- ----- ", swiched)
		gw, err := lotto.Get_lotto(Level6aus49, eingabe)

		gewinn_Text := build_message(gw, eingabe)

		message_Params := bot.SendMessageParams{
			ChatID:          update.Message.Chat.ID,
			Text:            gewinn_Text + "\n",
			MessageEffectID: Conf_tok["SUCCESS_EFFECT"],
		}
		_, err = b.SendMessage(ctx, &message_Params)
		if err != nil {
			fmt.Println("[ERROR] [SENDING MESSAGE FAILED DUE TO ]", err)
		}

	case LevelWitz:
		eingabe := update.Message.Text
		fmt.Println("Level: ", handlerScope)
		l, swiched := checkLevelSwitch(eingabe, ctx, b, update)
		if swiched {
			sendFaultMessage(handlerScope, fmt.Sprintf("Level %d sitched to Level %d ", LevelJP, l), ctx, b, update)
			return
		}
		fmt.Println(" ----------- ----- ", swiched)
		handler_witz(ctx, b, update)
		return

	case LevelLotto:
		eingabe := update.Message.Text
		fmt.Println("Level: ", handlerScope)
		l, swiched := checkLevelSwitch(eingabe, ctx, b, update)
		if swiched {
			sendFaultMessage(handlerScope, fmt.Sprintf("Level %d sitched to Level %d ", LevelLotto, l), ctx, b, update)
			return
		}
		fmt.Println(" ----------- ----- ", swiched)
		switch eingabe {
		case "NOCH":
			handler_witz(ctx, b, update)
			return
		case "Menue":
			handler_menue(ctx, b, update)
			return
		case "/menu":
			handler_menue(ctx, b, update)
			return
		case "Lotto-6aus49":
			handler_lotto(ctx, b, update)
			return
		case "Eurojackpot":
			handler_eurojackpot(ctx, b, update)
			return
		default:
			sendFaultMessage(LevelLotto, "Unbekannt", ctx, b, update)
		}

	case LevelTodo:
	case LevelEinkauf:
	case LevelObst:
		eingabe := update.CallbackQuery.Data
		fmt.Println("EINGABE: ", eingabe)
		if !checkInlineReplyBtnText(eingabe, LevelObst) {
			sendFaultMessage(LevelObst, "Ich hab das nicht in der Liste.\nBitte, Buttons benutzen!!!", ctx, b, update)
			return
		}
		go handler_update_inline(ctx, b, update)
		return
	case LevelKino:
		eingabe := update.Message.Text
		switch eingabe {
		case "HEUTE":
			handler_cinema(ctx, b, update)
			return
		case "WOCHE":
			handler_cinema(ctx, b, update)
			return
		}
	case LevelKinoHeute:
		eingabe := update.Message.Text
		switch eingabe {
		case "HEUTE":
			handler_cinema(ctx, b, update)
			return
		case "WOCHE":
			handler_cinema(ctx, b, update)
			return
		}
	case LevelKinoWoche:
		eingabe := update.Message.Text
		switch eingabe {
		case "HEUTE":
			handler_cinema(ctx, b, update)
			return
		case "WOCHE":
			handler_cinema(ctx, b, update)
			return
		}

	}
	// eingabe := update.Message.Text
	// if eingabe == "Witz" {
	// 	handlerScope = LevelWitz
	// 	handler_witz(ctx, b, update)
	// 	return
	// }
	// if eingabe == "/menu" {
	// 	handlerScope = LevelMenu
	// 	fmt.Println("menue2")
	// 	handler_menue(ctx, b, update)
	// 	return
	// }

}

func checkInlineReplyBtnText(t string, level int) bool {
	switch level {
	case LevelObst:
		for x := range inlineKeyboardMarkup.InlineKeyboard {
			for y := range inlineKeyboardMarkup.InlineKeyboard[x] {
				btnText := inlineKeyboardMarkup.InlineKeyboard[x][y]
				fmt.Println("\tbtnText: ", btnText)
				btnText.Text = strings.Fields(btnText.Text)[0]
				if strings.Contains(t, btnText.Text) {
					fmt.Println("\t    btnText: ", btnText)
					return true
				}
			}
		}
		return false
	}
	return true
}

func checkLevelSwitch(eingabe string, ctx context.Context, b *bot.Bot, update *models.Update) (int, bool) {

	switch handlerScope {
	case LevelMenu:
		switch eingabe {
		case "/eurojackpot":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_eurojackpot(ctx, b, update)
			return LevelJP, true
		case "Eurojackpot":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_eurojackpot(ctx, b, update)
			return LevelJP, true
		case "/lotto":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_lotto_default(ctx, b, update)
			return LevelLotto, true
		case "Lotto":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_lotto_default(ctx, b, update)
			return LevelLotto, true
		case "Witz":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_witz(ctx, b, update)
			return LevelWitz, true
		case "Todo":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_todo(ctx, b, update)
			return LevelTodo, true
		case "Kino":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_cinema(ctx, b, update)
			return LevelKino, true
		default:
			return handlerScope, false
		}
	case LevelStart:
	case LevelLotto:
		switch eingabe {
		case "/lotto_6aus49":
			handler_lotto(ctx, b, update)
			return Level6aus49, true
		case "/eurojackpot":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_eurojackpot(ctx, b, update)
			return LevelJP, true
		case "Eurojackpot":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_eurojackpot(ctx, b, update)
			return LevelJP, true
		case "Lotto-6aus49":
			handler_lotto(ctx, b, update)
			return Level6aus49, true
		case "Menue":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_menue(ctx, b, update)
			return LevelMenu, true
		case "/menu":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_menue(ctx, b, update)
			return LevelMenu, true
		case "Witz":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_witz(ctx, b, update)
			return LevelWitz, true
		case "Todo":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_todo(ctx, b, update)
			return LevelTodo, true
		case "Kino":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_cinema(ctx, b, update)
			return LevelKino, true
		default:
			return handlerScope, false
		}
	case LevelJP:
		switch eingabe {
		case "/lotto":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_lotto_default(ctx, b, update)
			return LevelLotto, true
		case "Lotto":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_lotto_default(ctx, b, update)
			return LevelLotto, true
		case "/lotto_6aus49":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_lotto(ctx, b, update)
			return Level6aus49, true
		case "/eurojackpot":
			handler_eurojackpot(ctx, b, update)
			return LevelJP, true
		case "Lotto-6aus49":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_lotto(ctx, b, update)
			return Level6aus49, true
		case "/menu":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_menue(ctx, b, update)
			return LevelMenu, true
		case "Menue":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_menue(ctx, b, update)
			return LevelMenu, true
		case "Witz":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_witz(ctx, b, update)
			return LevelWitz, true
		case "Todo":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_todo(ctx, b, update)
			return LevelTodo, true
		case "Kino":
			JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
			handler_cinema(ctx, b, update)
			return LevelKino, true
		default:
			return handlerScope, false
		}
	case Level6aus49:
		JackpotDatumScope = Scope(time.Now().Format("02.01.2006"))
		switch eingabe {
		case "/lotto":
			handler_lotto_default(ctx, b, update)
			return LevelLotto, true
		case "/lotto_6aus49":
			handler_lotto(ctx, b, update)
			return Level6aus49, true
		case "/eurojackpot":
			handler_eurojackpot(ctx, b, update)
			return LevelJP, true
		case "Eurojackpot":
			handler_eurojackpot(ctx, b, update)
			return LevelJP, true
		case "/menu":
			handler_menue(ctx, b, update)
			return LevelMenu, true
		case "Menue":
			handler_menue(ctx, b, update)
			return LevelMenu, true
		case "Witz":
			handler_witz(ctx, b, update)
			return LevelWitz, true
		case "Todo":
			handler_todo(ctx, b, update)
			return LevelTodo, true
		case "Kino":
			handler_cinema(ctx, b, update)
			return LevelKino, true
		default:
			return handlerScope, false
		}
	case LevelLebensmittel:
	case LevelWitz:
	case LevelEinkauf:
	case LevelKino:
	case LevelKinoHeute:
	case LevelKinoWoche:
	}

	return handlerScope, false
}
