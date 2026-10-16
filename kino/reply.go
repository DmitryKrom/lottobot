package kino

import (
	"fmt"

	"github.com/go-telegram/bot/models"
)

type MyKeyBoardButton struct {
	KeyBoardButton models.KeyboardButton
}

func NewKeyBoardButton() *MyKeyBoardButton {
	return &MyKeyBoardButton{
		KeyBoardButton: models.KeyboardButton{},
	}
}

func (kb *MyKeyBoardButton) Text(text string) {
	kb.KeyBoardButton.Text = text
}

func (kb *MyKeyBoardButton) Style(style string) {
	kb.KeyBoardButton.Style = style
}

type MyKeyBoardLine struct {
	KeyBoardLine []models.KeyboardButton
}

func NewKeyBoardLine(count int, lineNum int) MyKeyBoardLine {
	return MyKeyBoardLine{
		KeyBoardLine: make([]models.KeyboardButton, count),
	}
}
func (kbl *MyKeyBoardLine) AppendButton(kb *MyKeyBoardButton) {
	kbl.KeyBoardLine = append(kbl.KeyBoardLine, kb.KeyBoardButton)
}

type MyKeyBoard struct {
	KeyBoard [][]models.KeyboardButton
}

func NewKeyBoard(lines int) *MyKeyBoard {
	return &MyKeyBoard{
		KeyBoard: make([][]models.KeyboardButton, lines),
	}
}

func (kb *MyKeyBoard) Build(btns []MyKeyBoardButton, lines int, itemsProLine int) *models.ReplyKeyboardMarkup {
	var rkb = &models.ReplyKeyboardMarkup{}
	anfang := 0
	for i := range lines {
		line := NewKeyBoardLine(itemsProLine, i)
		kb.KeyBoard[i] = make([]models.KeyboardButton, 0)
		for x, btn := range btns[anfang : anfang+itemsProLine] {
			fmt.Println("			x: ", x, "  anfang : ", anfang)
			line.AppendButton(&btn)
		}
		kb.KeyBoard[i] = line.KeyBoardLine
		anfang = anfang + itemsProLine
	}
	rkb.Keyboard = kb.KeyBoard
	return rkb
}
