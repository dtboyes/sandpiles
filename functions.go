package main

import (
	"math/rand"
	"runtime"
)

// ComputeSteadyState computes the steady state of the board either in parallel or serially
func (b Board) ComputeSteadyState(parallel bool) Board {
	if parallel {
		numProcs := runtime.NumCPU()
		runtime.GOMAXPROCS(numProcs)
		b.SandpileMultiProc(numProcs)
	} else {
		b.SandpileSerial()
	}
	return b
}

// SandpileMultiProc continues to check if the board is topple-able, then topples it accordingly
// it creates a map, which maps each processor in numProcs to a board containing the top and bottom slices of each sub_board
// it then creates a channel to wait for all of the goroutines to finish, then adds the top and bottom slices to the board
func (b Board) SandpileMultiProc(numProcs int) Board {
	rows := nRows(b)
	cols := nCols(b)
	for b.CanToppleBoard() {
		ch := make(chan bool, numProcs)

		// creation and intiialization of map to hold top and bottom slices of each sub_board
		top_and_bottom := make(map[int]Board)
		for i := 0; i < numProcs; i++ {
			top_and_bottom[i] = make(Board, 2)
			for j := range top_and_bottom[i] {
				top_and_bottom[i][j] = make([]int, cols)
			}
		}

		// numProcs goroutines to topple each sub_board and update the top and bottom slices
		for i := 0; i < numProcs; i++ {
			start := i * rows / numProcs
			end := (i + 1) * rows / numProcs
			sub_board := b[start:end]

			go sub_board.SandpileSingleProc(top_and_bottom, i, ch)
		}
		// waiting for all goroutines to finish
		for i := 0; i < numProcs; i++ {
			<-ch
		}

		// adding the top and bottom slices to the board using AddTwoArrays function
		for i := 0; i < numProcs; i++ {
			start := i * rows / numProcs
			end := (i + 1) * rows / numProcs
			if start == 0 {
				b[end] = AddTwoArrays(b[end], top_and_bottom[i][1])
			} else if i == numProcs-1 {
				b[start-1] = AddTwoArrays(b[start-1], top_and_bottom[i][0])
			} else {
				b[start-1] = AddTwoArrays(b[start-1], top_and_bottom[i][0])
				b[end] = AddTwoArrays(b[end], top_and_bottom[i][1])
			}
		}
	}

	return b
}

// SandpileSingleProc topples the board once for a single sub_board
// the call to ToppleParallel handles the updating of the rows contained above and below the sub_board
func (b Board) SandpileSingleProc(top_and_bottom map[int]Board, index int, ch chan bool) {
	for r := range b {
		for c := range b[r] {
			b.ToppleParallel(r, c, top_and_bottom, index)
		}
	}

	ch <- true
}

// SandpileSerial continues to check if the board is topple-able, then topples it accordingly cell by cell
func (b Board) SandpileSerial() Board {
	for b.CanToppleBoard() {
		for r := range b {
			for c := range b[r] {
				b.ToppleSerial(r, c)
			}
		}
	}
	return b
}

// CanToppleBoard checks if the board can be toppled by checking if any cell has 4 or more grains of sand
func (b Board) CanToppleBoard() bool {
	for r := range b {
		for c := range b[r] {
			if b[r][c] >= 4 {
				return true
			}
		}
	}
	return false
}

// ToppleParallel is nearly identical to ToppleSerial, except that it takes in a map for top and bottom slices, and the index of the processor associated with those slices
// it updates top and bottom slices, and the board, accordingly
func (b Board) ToppleParallel(r, c int, top_and_bottom map[int]Board, index int) {
	if b[r][c] >= 4 {
		b[r][c] -= 4
		if r == 0 {
			top_and_bottom[index][0][c]++
		}
		if r == len(b)-1 {
			top_and_bottom[index][1][c]++
		}
		if r > 0 {
			b[r-1][c]++
		}
		if r < len(b)-1 {
			b[r+1][c]++
		}
		if c > 0 {
			b[r][c-1]++
		}
		if c < len(b[r])-1 {
			b[r][c+1]++
		}
	}
}

// ToppleSerial just topples a single cell in the board
func (b Board) ToppleSerial(r, c int) {
	if b[r][c] >= 4 {
		b[r][c] -= 4
		if r > 0 {
			b[r-1][c]++
		}
		if r < len(b)-1 {
			b[r+1][c]++
		}
		if c > 0 {
			b[r][c-1]++
		}
		if c < len(b[r])-1 {
			b[r][c+1]++
		}
	}
}

// AddTwoArrays adds two arrays together and returns the result
func AddTwoArrays(a, b []int) []int {
	for i := range a {
		a[i] += b[i]
	}
	return a
}

// nRows and nCols are helper functions to get the number of rows and columns in a board
func nRows(b Board) int {
	return len(b)
}

func nCols(b Board) int {
	return len(b[0])
}

// InitializeBoard initializes the board with a pile of sand in the center or randomly
func (b Board) InitializeBoard(size, pile int, mode string) {
	if mode == "central" {
		b[size/2][size/2] = pile
	} else if mode == "random" {
		x := make([]int, 100)
		y := make([]int, 100)
		for i := 0; i < 100; i++ {
			x[i] = rand.Intn(size)
			y[i] = rand.Intn(size)
		}
		for index := 0; index < pile; index++ {
			i := rand.Intn(100)
			b[x[i]][y[i]]++
		}
	}
}

// InitializeEmptyBoard initializes an empty board of size size
func InitializeEmptyBoard(size int) Board {
	b := make(Board, size)
	for i := range b {
		b[i] = make([]int, size)
	}
	return b
}

// CopyBoard copies all of the elements of a board into a new board
func CopyBoard(b Board) Board {
	c := make(Board, len(b))
	for i := range b {
		c[i] = make([]int, len(b[i]))
		copy(c[i], b[i])
	}
	return c
}

// NOT NECESSARY FOR ASSIGNMENT ------------------------------------------------
// RecursiveSandpiles does the same as SandpileSerial but recursively
// it handles the initialization of the board and the recursive calls on each pile
func (b Board) RecursiveSandpiles(size, pile int, mode string) {
	if mode == "central" {
		b[size/2][size/2] = pile
		ToppleRecursive(b, size/2, size/2)
	} else if mode == "random" {
		x := make([]int, 5)
		y := make([]int, 5)
		for i := 0; i < 5; i++ {
			x[i] = rand.Intn(size)
			y[i] = rand.Intn(size)
		}
		for index := 0; index < pile; index++ {
			i := rand.Intn(5)
			b[x[i]][y[i]]++
		}
		for i := 0; i < 5; i++ {
			ToppleRecursive(b, x[i], y[i])
		}
	}
}

// ToppleRecursive is a recursive function that topples a cell and calls itself on the cells above, below, to the left, and to the right
func ToppleRecursive(b Board, r, c int) {
	for b[r][c] >= 4 {
		b[r][c] -= 4
		if r > 0 {
			b[r-1][c]++
			ToppleRecursive(b, r-1, c)
		}
		if r < len(b)-1 {
			b[r+1][c]++
			ToppleRecursive(b, r+1, c)
		}
		if c > 0 {
			b[r][c-1]++
			ToppleRecursive(b, r, c-1)
		}
		if c < len(b[r])-1 {
			b[r][c+1]++
			ToppleRecursive(b, r, c+1)
		}
	}
}

// END OF BLOCK ---------------------------------------------------------------
