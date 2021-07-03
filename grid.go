package main

// Grid represents alive cells as 1 and dead cells as 0
type Grid [][]int

const (
	GridMaxCols   = 50
	GridMaxRows   = 25
	GridCellAlive = 1
	GridCellDead  = 0
)

func NewGrid(width, height int) *Grid {
	grid := Grid{}
	for i := 0; i < width; i++ {
		row := []int{}
		for j := 0; j < height; j++ {
			row = append(row, 0)
		}
		grid = append(grid, row)
	}
	return &grid
}

func (g Grid) toggleCell(x, y int) {
	if x >= 0 && x < g.width()*2 && y >= 0 && y < g.height() {
		// round to nearest even value, then scale down to grid size
		x2 := (x - (x % 2)) / 2
		val := g[x2][y]
		if val == GridCellDead {
			g[x2][y] = GridCellAlive
		} else {
			g[x2][y] = GridCellDead
		}
	}
}

func (g Grid) width() int {
	return len(g)
}

func (g Grid) height() int {
	return len(g[0])
}
