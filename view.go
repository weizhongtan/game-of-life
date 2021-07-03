package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

type View struct {
	screen  tcell.Screen
	grid    *Grid
	running bool
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

	g := NewGrid()
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

func (v *View) drawBox(x1, y1, x2, y2 int) {
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
}

func (v *View) drawCell(x, y int) {
	v.grid.drawCell(x, y)
}

func (v *View) update() {
	if v.running {
		g := *v.grid

		// nb represents the number of neighbors for each position in the grid
		nb := *NewGrid()

		// count the neighbors for each position in the grid
		for i := 0; i < len(g); i++ {
			for j, col := 0, g[i]; j < len(col); j++ {
				count := 0
				// check neighbors
				for a := i - 1; a <= i+1; a++ {
					for b := j - 1; b <= j+1; b++ {
						if a >= 0 && a < GridMaxCols && b >= 0 && b < GridMaxRows {
							if g[a][b] == GridCellAlive {
								count++
							}
						}
					}
				}
				nb[i][j] = count
			}
		}

		// calculate next generation and store in nb
		for i := 0; i < len(g); i++ {
			for j, col := 0, g[i]; j < len(col); j++ {
				count := nb[i][j]

				// if the count of neighbors including self is 3
				//   or
				// the count of neighbors including self is 4 and self is alive
				//   -> next generation cell position is alive
				if count == 3 || (count == 4 && g[i][j] == GridCellAlive) {
					g[i][j] = GridCellAlive
				} else {
					g[i][j] = GridCellDead
				}
			}
		}
	}
}

func (v *View) render() {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorReset)
	// draw outer box
	v.drawBox(0, 0, GridMaxCols*2+1, GridMaxRows+1)
	// render grid within outer box
	g := *v.grid
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
	var runningMsg string
	if v.running {
		runningMsg = "running"
	} else {
		runningMsg = "paused"
	}
	lines := []string{
		"Controls:",
		"[left click] add cell",
		"[r] toggle simulation",
		"[esc] exit",
		fmt.Sprintf("simulation: %s", runningMsg),
	}
	for i, msg := range lines {
		v.drawText(0, i+GridMaxRows+2, (GridMaxCols*2 + 2), i+GridMaxRows+2, msg)
	}
	// Update screen
	v.screen.Show()
}
