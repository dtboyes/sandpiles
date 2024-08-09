package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Board [][]int

// for some reason, my laptop runs parallel really slowly
// I've tried running this using other computers and I did notice a speedup
// when running on my laptop, I got about 32 seconds for parallel and 29 seconds for serial on ./sandpile 1000 50000 central
// the same command on another computer resulted in 20 seconds for parallel and 27 seconds for serial, so I'm not sure what might be happening
func main() {
	var recursive bool = true

	size, _ := strconv.Atoi(os.Args[1])
	pile, _ := strconv.Atoi(os.Args[2])
	mode := os.Args[3]

	start := time.Now()

	if recursive { // just a fun recursive example that I left in (not necessary for assignment so recursion is off by default)
		var b Board = InitializeEmptyBoard(size)
		b.RecursiveSandpiles(size, pile, mode)
		c := b.DrawToCanvas()
		c.SaveToPNG("sandpiles.png")
		elapsed := time.Since(start)
		fmt.Println("Time taken for recursive: ", elapsed)
	} else {
		var b1 Board = InitializeEmptyBoard(size)
		var b2 Board = InitializeEmptyBoard(size)
		b1.InitializeBoard(size, pile, mode)
		b2 = CopyBoard(b1)

		start1 := time.Now()
		b1.ComputeSteadyState(true)
		elapsed1 := time.Since(start1)
		c1 := b1.DrawToCanvas()
		c1.SaveToPNG("parallel.png")
		fmt.Println("Time taken for parallel: ", elapsed1)

		fmt.Println("Now computing serially...")
		start2 := time.Now()
		b2.ComputeSteadyState(false)
		elapsed2 := time.Since(start2)
		c2 := b2.DrawToCanvas()
		c2.SaveToPNG("serial.png")
		fmt.Println("Time taken for serial: ", elapsed2)
	}

}
