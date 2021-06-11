package main

type Column []int

type Grid struct {
	cols Column
}

func NewGrid() *Grid {
	g := &Grid{}
	return g
}
