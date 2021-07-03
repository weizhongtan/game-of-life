package main

import (
	"github.com/gdamore/tcell/v2"
)

func main() {
	view := NewView()

	// Event loop
	buttonDown := false
	for {
		// convert grid into setContent calls
		view.render()

		// Update screen
		view.screen.Show()

		// Poll event
		ev := view.screen.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			view.sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				view.quit()
			} else if ev.Key() == tcell.KeyCtrlL {
				view.screen.Sync()
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				view.screen.Clear()
			}
		case *tcell.EventMouse:
			x, y := ev.Position()
			button := ev.Buttons()
			// Only process button events, not wheel events
			button &= tcell.ButtonMask(0xff)

			switch button {
			case tcell.ButtonNone:
				if buttonDown {
					view.drawCell(x, y)
					buttonDown = false
				}
			default:
				if !buttonDown {
					buttonDown = true
				}
			}
		}
	}
}
