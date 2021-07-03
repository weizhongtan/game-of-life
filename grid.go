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
