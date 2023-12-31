package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

type cellState int

const (
	empty cellState = iota
	red
	blue
)

type tappableArea struct {
	widget.BaseWidget
	onTap func(x, y int)
	x, y  int
}

func newTappableArea(x, y int, onTap func(x, y int), cellStates [][]cellState) *tappableArea {
	t := &tappableArea{
		onTap: func(x, y int) {
			if cellStates[x][y] == empty {
				onTap(x, y)
			}
		},
		x: x,
		y: y,
	}
	t.ExtendBaseWidget(t)
	return t
}

func (t *tappableArea) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.Transparent)
	return widget.NewSimpleRenderer(rect)
}

func (t *tappableArea) Tapped(_ *fyne.PointEvent) {
	if t.onTap != nil {
		t.onTap(t.x, t.y)
	}
}

// Grid manages the game grid.
type Grid struct {
	container     *fyne.Container
	cellStates    [][]cellState
	dotsContainer *fyne.Container
	cellSize      int
	gridOffsetX   float32
	gridOffsetY   float32
	gridSize      int
	blueDotCount  *int
	redDotCount   *int
	dotCount      *int
	onDotPlaced   func() // Callback function
	timer         *Timer
	gameWindow    *GameWindow
}

func NewGrid(gridSize, cellSize int, gridOffsetX, gridOffsetY float32, onDotPlaced func(), timer *Timer, gameWindow *GameWindow) *Grid {
	cellStates := make([][]cellState, gridSize)
	for i := range cellStates {
		cellStates[i] = make([]cellState, gridSize)
	}

	return &Grid{
		container:     container.NewWithoutLayout(),
		cellStates:    cellStates,
		dotsContainer: container.NewWithoutLayout(),
		cellSize:      cellSize,
		gridOffsetX:   gridOffsetX,
		gridOffsetY:   gridOffsetY,
		gridSize:      gridSize,
		blueDotCount:  new(int),
		redDotCount:   new(int),
		dotCount:      new(int),
		onDotPlaced:   onDotPlaced,
		timer:         timer,
		gameWindow:    gameWindow,
	}
}

func (g *Grid) DrawGrid() {
	// Adjust the starting point based on the radius of the dots
	adjustedOffsetX := g.gridOffsetX + float32(g.cellSize)/2
	adjustedOffsetY := g.gridOffsetY + float32(g.cellSize)/2

	for i := 0; i < g.gridSize; i++ {
		// Draw vertical lines
		vLine := canvas.NewLine(color.Black)
		vLine.Move(fyne.NewPos(adjustedOffsetX+float32(i)*float32(g.cellSize), adjustedOffsetY))
		vLine.Resize(fyne.NewSize(2, float32(g.cellSize)*(float32(g.gridSize)-1)))
		g.container.Add(vLine)

		// Draw horizontal lines
		hLine := canvas.NewLine(color.Black)
		hLine.Move(fyne.NewPos(adjustedOffsetX, adjustedOffsetY+float32(i)*float32(g.cellSize)))
		hLine.Resize(fyne.NewSize(float32(g.cellSize)*(float32(g.gridSize)-1), 2))
		g.container.Add(hLine)
	}

	// Create and add tappable areas for each grid cell
	for y := 0; y < g.gridSize; y++ {
		for x := 0; x < g.gridSize; x++ {
			area := newTappableArea(x, y, func(x, y int) {
				if g.cellStates[x][y] != empty {
					// Skip if the cell is already filled
					return
				}

				var currentColor cellState
				if *g.dotCount%2 == 0 {
					currentColor = blue
				} else {
					currentColor = red
				}

				// Place a dot of the determined color
				g.PlaceDot(currentColor, x, y)
				// Refresh the grid container to show the new dot
				g.container.Refresh()
			}, g.cellStates) // Pass g.cellStates as the fourth argument

			area.Resize(fyne.NewSize(float32(g.cellSize), float32(g.cellSize)))
			area.Move(fyne.NewPos(float32(x)*float32(g.cellSize)+g.gridOffsetX, float32(y)*float32(g.cellSize)+g.gridOffsetY))
			g.container.Add(area)
		}
	}
}

// Additional methods for handling grid functionalities
// ...
