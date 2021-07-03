package main

import "github.com/gdamore/tcell/v2"

type Controller struct {
	view       *View
	buttonDown bool
}

func (c *Controller) processInput() {
	view := c.view

	// don't block if there are no events to handle
	if !view.screen.HasPendingEvent() {
		return
	}

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
		} else if ev.Rune() == ' ' {
			c.view.running = !c.view.running
		}
	case *tcell.EventMouse:
		x, y := ev.Position()
		button := ev.Buttons()
		// Only process button events, not wheel events
		button &= tcell.ButtonMask(0xff)

		switch button {
		case tcell.ButtonNone:
			if c.buttonDown {
				// offset to account for borders
				view.toggleCell(x-1, y-1)
				c.buttonDown = false
			}
		default:
			if !c.buttonDown {
				c.buttonDown = true
			}
		}
	}
}
