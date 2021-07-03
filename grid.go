package main

// Grid represents alive cells as 1 and dead cells as 0
type Grid [][]int

const (
	GridMaxCols   = 15
	GridMaxRows   = 15
	GridCellAlive = 1
	GridCellDead  = 0
)

func NewGrid() *Grid {
	grid := Grid{}
	for i := 0; i < GridMaxCols; i++ {
		row := []int{}
		for j := 0; j < GridMaxRows; j++ {
			row = append(row, 0)
		}
		grid = append(grid, row)
	}
	return &grid
}

func (g *Grid) drawCell(x, y int) {
	if x < GridMaxCols*2 && y < GridMaxRows {
		// round to nearest even value, then scale down to grid size
		x2 := (x - (x % 2)) / 2
		(*g)[x2][y] = GridCellAlive
	}
}
