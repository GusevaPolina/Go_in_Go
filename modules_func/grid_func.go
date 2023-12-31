package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

var gridSize int = 9

// Dot counters
var dotCount, blueDotCount, redDotCount int

type cellState int

const (
	empty cellState = iota
	red
	blue
)

type tappableArea struct {
	widget.BaseWidget
	onTap func()
}

func newTappableArea(x, y int, onTap func(x, y int)) *tappableArea {
	t := &tappableArea{}
	t.ExtendBaseWidget(t)
	t.onTap = func() { onTap(x, y) }
	return t
}

func (t *tappableArea) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}

func (t *tappableArea) Tapped(_ *fyne.PointEvent) {
	if t.onTap != nil {
		t.onTap()
	}
}

func drawGrid(container *fyne.Container, offsetX, offsetY float32, cellSize float32, gridSize int) {
	// Adjust the starting point based on the radius of the dots
	adjustedOffsetX := offsetX + cellSize/2
	adjustedOffsetY := offsetY + cellSize/2

	for i := 0; i < gridSize; i++ {
		// Draw vertical lines
		vLine := canvas.NewLine(color.Black)
		// vLine.StrokeWidth = 2
		vLine.Move(fyne.NewPos(adjustedOffsetX+float32(i)*cellSize, adjustedOffsetY))
		vLine.Resize(fyne.NewSize(2, cellSize*(float32(gridSize)-1)))
		container.Add(vLine)

		// Draw horizontal lines
		hLine := canvas.NewLine(color.Black)
		// hLine.StrokeWidth = 2
		hLine.Move(fyne.NewPos(adjustedOffsetX, adjustedOffsetY+float32(i)*cellSize))
		hLine.Resize(fyne.NewSize(cellSize*(float32(gridSize)-1), 2))
		container.Add(hLine)
	}
}

func clearGrid() {
	dotCount, blueDotCount, redDotCount = 0, 0, 0
	resetTimer()
}

// Additional methods for handling grid functionalities
// ...
