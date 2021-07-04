package main

// Grid represents alive cells as 1 and dead cells as 0
type Grid [][]int

const (
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

func (g *Grid) toggleCell(x, y int) {
	grid := *g
	if x >= 0 && x < grid.width()*2 && y >= 0 && y < grid.height() {
		// round to nearest even value, then scale down to grid size
		x2 := (x - (x % 2)) / 2
		val := grid[x2][y]
		if val == GridCellDead {
			grid[x2][y] = GridCellAlive
		} else {
			grid[x2][y] = GridCellDead
		}
	}
}

func (g *Grid) width() int {
	return len(*g)
}

func (g *Grid) height() int {
	return len((*g)[0])
}
