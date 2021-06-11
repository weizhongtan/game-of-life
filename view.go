package main

import (
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

type View struct {
	screen tcell.Screen
}

func NewView() *View {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	v := View{s}
	return &v
}

func (v *View) quit() {
	v.screen.Fini()
	os.Exit(0)
}

func (v *View) sync() {
	v.screen.Sync()
}

func (v *View) drawText(x1, y1, x2, y2 int, style tcell.Style, text string) {
	row, col := y1, x1
	for _, r := range text {
		v.screen.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func (v *View) drawBox(x1, y1, x2, y2 int, text string) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorPurple)
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			v.screen.SetContent(col, row, ' ', nil, style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		v.screen.SetContent(col, y1, tcell.RuneHLine, nil, style)
		v.screen.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		v.screen.SetContent(x1, row, tcell.RuneVLine, nil, style)
		v.screen.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		v.screen.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		v.screen.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		v.screen.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		v.screen.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}

	v.drawText(x1+1, y1+1, x2-1, y2-1, style, text)
}

func (v *View) drawCell(x, y int) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset)
	v.screen.SetContent(x, y, tcell.RuneBlock, nil, style)
	v.screen.SetContent(x+1, y, tcell.RuneBlock, nil, style)
}
