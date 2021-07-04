package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

func getUI(isRunning bool) []string {
	var toggleMsg string
	if isRunning {
		toggleMsg = "pause"
	} else {
		toggleMsg = "run"
	}
	lines := []string{
		"[left click] toggle cell",
		fmt.Sprintf("[space]      %s simulation", toggleMsg),
		"[r]          reset",
		"[esc/q]      exit",
	}
	return lines
}

type View struct {
	screen  tcell.Screen
	currGen Grid
	running bool
}

func NewView() *View {
	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	// adjust for screen size
	// width - box left + right borders
	// height - box top + bottom borders + controls UI
	sw, sh := s.Size()
	w, h := (sw/2)-1, sh-(2+len(getUI(false)))

	g := NewGrid(w, h)
	v := View{s, g, false}
	return &v
}

func (v *View) quit() {
	v.screen.Fini()
	os.Exit(0)
}

func (v *View) sync() {
	v.screen.Sync()
}

func (v *View) drawText(x1, y1, x2, y2 int, text string) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset)
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

// drawTextLine draws a line of text at a given line number
func (v *View) drawTextLine(line int, text string) {
	w := (v.currGen.width() * 2) + 2
	// %-*s explained:
	// "-" right justify
	// "*" pass width in
	v.drawText(0, line, w, line+1, fmt.Sprintf("%-*s", w, text))
}

func (v *View) drawBox(x1, y1, x2, y2 int, color tcell.Color) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	var style tcell.Style
	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			// checkerboard pattern
			var color tcell.Color
			if row%2^((col+1)/2)%2 == 0 {
				color = tcell.ColorSlateGray
			} else {
				color = tcell.ColorDarkSlateGray
			}
			style = tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(color)
			v.screen.SetContent(col, row, ' ', nil, style)
		}
	}

	style = tcell.StyleDefault.Foreground(color).Background(tcell.ColorReset)

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
}

func (v *View) toggleCell(x, y int) {
	v.currGen.toggleCell(x, y)
}

func wrap(val, min, max int) int {
	if val < min {
		return max
	}
	if val > max {
		return min
	}
	return val
}

func (v *View) update() {
	if v.running {
		currGen := v.currGen
		nextGen := NewGrid(currGen.width(), currGen.height())

		// count the neighbours for each position in the grid
		for i := 0; i < len(currGen); i++ {
			for j, col := 0, currGen[i]; j < len(col); j++ {
				count := 0
				// check all 9 neighbours including self
				for a := i - 1; a <= i+1; a++ {
					for b := j - 1; b <= j+1; b++ {
						// wrap around the matrix
						aWrap, bWrap := wrap(a, 0, currGen.width()-1), wrap(b, 0, currGen.height()-1)
						if currGen[aWrap][bWrap] == GridCellAlive {
							count++
						}
					}
				}
				// if the count of neighbours including self is 3
				//   or
				// the count of neighbours including self is 4 and self is alive
				//   -> next generation cell position is alive
				if count == 3 || (count == 4 && currGen[i][j] == GridCellAlive) {
					nextGen[i][j] = GridCellAlive
				} else {
					nextGen[i][j] = GridCellDead
				}
			}
		}

		v.currGen = nextGen
	}
}

func (v *View) reset() {
	v.currGen = NewGrid(v.currGen.width(), v.currGen.height())
}

func (v *View) render() {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset)

	// draw outer box
	var color tcell.Color
	if v.running {
		color = tcell.ColorGreen
	} else {
		color = tcell.ColorWhite
	}
	v.drawBox(0, 0, v.currGen.width()*2+1, v.currGen.height()+1, color)

	// render grid within outer box
	g := v.currGen
	for i := 0; i < len(g); i++ {
		col := g[i]
		for j := 0; j < len(col); j++ {
			if g[i][j] == GridCellAlive {
				// screen needs 2 cells per column
				// offset of one column and one row into the screen
				x := 1 + i*2
				y := 1 + j
				v.screen.SetContent(x, y, tcell.RuneBlock, nil, style)
				v.screen.SetContent(x+1, y, tcell.RuneBlock, nil, style)
			}
		}
	}
	// draw game status
	lines := getUI(v.running)
	for i, msg := range lines {
		v.drawTextLine(g.height()+2+i, msg)
	}
	// Update screen
	v.screen.Show()
}
