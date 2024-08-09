package main

import (
	"canvas"
)

func (b Board) DrawToCanvas() canvas.Canvas {
	if b == nil {
		panic("Can't Draw a nil Board.")
	}

	// set a new square canvas
	c := canvas.CreateNewCanvas(len(b), len(b[0]))

	// create a black background
	c.SetFillColor(canvas.MakeColor(0, 0, 0))
	c.ClearRect(0, 0, len(b), len(b[0]))
	c.Fill()

	for i := range b {
		for j := range b[i] {
			if b[i][j] == 0 {
				c.SetFillColor(canvas.MakeColor(0, 0, 0))
			}
			if b[i][j] == 1 {
				c.SetFillColor(canvas.MakeColor(85, 85, 85))
			}
			if b[i][j] == 2 {
				c.SetFillColor(canvas.MakeColor(170, 170, 170))
			}
			if b[i][j] == 3 {
				c.SetFillColor(canvas.MakeColor(255, 255, 255))
			}
			c.ClearRect(i, j, i+1, j+1)
			c.Fill()
		}
	}
	// we want to return an image!
	return c
}
