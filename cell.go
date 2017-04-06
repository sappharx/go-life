package main

import (
	"math/rand"
	"time"
	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	squarePoints = []float32{
		// bottom left right-triangle
		-0.5,  0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,

		// top right right-triangle
		-0.5,  0.5, 0,
		0.5,  0.5, 0,
		0.5, -0.5, 0,
	}
)

type cell struct {
	drawable	uint32

	alive     bool
	nextState bool

	x int
	y int
}

func (c *cell) draw() {
	if !c.alive {
		return
	}

	gl.BindVertexArray(c.drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(squarePoints) / 3))
}

/**
 * checkState determines the state of the cell for the next tick of the game
 */
func (c *cell) checkState(cells [][]*cell) {
	c.alive = c.nextState
	c.nextState = c.alive

	liveCount := c.liveNeighbors(cells)

	if c.alive {
		// any live cell with fewer than two live neighbors dies, as if caused by underpopulation
		if liveCount < 2 {
			c.nextState = false
		}

		// any live cell with two or three live neighbors lives on to the next generation
		if liveCount == 2 || liveCount == 3	{
			c.nextState = true
		}

		// any live cell with more than three live neighbors dies, as if by overpopulation
		if liveCount > 3 {
			c.nextState = false
		}
	} else {
		// any dead cell with exactly three live neighbors becomes a live cell, as if by reproduction
		if liveCount == 3 {
			c.nextState = true
		}
	}
}

/**
 * liveNeighbors returns the number of the live neighbors for a cell
 */
func (c *cell) liveNeighbors(cells [][]*cell) int {
	var liveCount int

	add := func(x, y int) {
		// if we're at an edge, check the other side of the board
		if x == len(cells) {
			x = 0
		} else if x == -1 {
			x = len(cells) - 1
		}

		if y == len(cells[x]) {
			y = 0
		} else if y == -1 {
			y = len(cells[x]) - 1
		}

		if cells[x][y].alive {
			liveCount++
		}
	}

	add(c.x - 1, c.y)		// to the left
	add(c.x + 1, c.y)		// to the right
	add(c.x, c.y + 1)		// up
	add(c.x, c.y - 1)		// down
	add(c.x - 1, c.y + 1)	// top-left
	add(c.x + 1, c.y + 1)	// top-right
	add(c.x - 1, c.y - 1)	// bottom-left
	add(c.x + 1, c.y - 1)	// bottom-right

	return liveCount
}


func makeCells() [][]*cell {
	rand.Seed(time.Now().UnixNano())

	cells := make([][]*cell, rows, rows)

	for x := 0; x < rows; x++ {
		for y := 0; y < columns; y++ {
			c := newCell(x, y)

			c.alive = rand.Float64() < threshold
			c.nextState = c.alive

			cells[x] = append(cells[x], c)
		}
	}

	return cells
}

func newCell(x, y int) *cell {
	points := make([]float32, len(squarePoints), len(squarePoints))
	copy(points, squarePoints)

	for i := 0; i < len(points); i++ {
		var position float32
		var size float32

		switch i % 3 {
		case 0:	// x-index of a point
			size = 1.0 / float32(columns)
			position = float32(x) * size
		case 1:	// y-index of a point
			size = 1.0/ float32(rows)
			position = float32(y) * size
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (position * 2) - 1
		} else {
			points[i] = ((position + size) * 2) - 1
		}
	}

	return &cell{
		drawable: makeVao(points),

		x: x,
		y: y,
	}
}
