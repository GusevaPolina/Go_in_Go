package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

const (
	windowWidth  = 500
	windowHeight = 500
)

// Custom theme
type myTheme struct {
	fyne.Theme
}

func (m myTheme) ButtonColor() color.Color {
	return color.NRGBA{R: 0, G: 0, B: 0, A: 255} // Black
}

func (m myTheme) TextColor() color.Color {
	return color.White
}

func regenerateGrid(gridSize int, myWindow fyne.Window, gridSizeInput *widget.Entry) {
	cellSize := 400 / gridSize
	clearGrid()

	gridContainer := container.NewWithoutLayout()
	dotsContainer := container.NewWithoutLayout()

	backgroundImage := canvas.NewImageFromFile("image.png")
	backgroundImage.FillMode = canvas.ImageFillContain
	backgroundImage.Resize(fyne.NewSize(windowWidth, windowHeight))
	backgroundImage.Move(fyne.NewPos(0, 0))
	gridContainer.Add(backgroundImage)

	gridOffsetX := float64(windowWidth) * 0.1
	gridOffsetY := float64(windowHeight) * 0.1

	// Initialize cell states
	cellStates := make([][]cellState, gridSize)
	for i := range cellStates {
		cellStates[i] = make([]cellState, gridSize)
	}

	blueDotCountLabel := widget.NewLabel("Blue Dots: 0")
	redDotCountLabel := widget.NewLabel("Red Dots: 0")

	timeElapsedLabel = widget.NewLabel("Time: 0s")

	// Create a CanvasText for the banner
	gameEndBanner := canvas.NewText("The game is over", color.NRGBA{R: 0, G: 0, B: 0, A: 255}) // Black color
	gameEndBanner.TextStyle = fyne.TextStyle{Bold: true}
	gameEndBanner.Alignment = fyne.TextAlignCenter
	gameEndBanner.TextSize = 24 // Adjust the text size as needed
	gameEndBanner.Hide()        // Initially hidden

	// Function to check if the game is over and update banner visibility
	checkAndUpdateGameEnd := func() {
		if blueDotCount+redDotCount == gridSize*gridSize {
			gameEndBanner.Show()
			stopTimer()
		} else {
			gameEndBanner.Hide()
		}
	}

	// Function to update dot counters
	updateDotCountLabels := func() {
		blueDotCountLabel.SetText(fmt.Sprintf("Blue Dots: %d", blueDotCount))
		redDotCountLabel.SetText(fmt.Sprintf("Red Dots: %d", redDotCount))
	}

	// Create tappable areas
	for y := 0; y < gridSize; y++ {
		for x := 0; x < gridSize; x++ {
			area := newTappableArea(x, y, func(x, y int) {
				if cellStates[x][y] == empty {
					currentColor := red
					if dotCount%2 == 0 {
						currentColor = blue
					}
					placeDot(currentColor, x, y, cellStates, dotsContainer, float32(cellSize), float32(gridOffsetX), float32(gridOffsetY), &dotCount)
					if dotCount > 1 { // Skip the first step
						checkAndFillClusters(cellStates, dotsContainer, cellSize, float32(gridOffsetX), float32(gridOffsetY))
					}
					updateDotCountLabels()
				}
				checkAndUpdateGameEnd() // Update banner visibility
			})
			area.Resize(fyne.NewSize(float32(cellSize), float32(cellSize)))
			area.Move(fyne.NewPos(float32(x*cellSize)+float32(gridOffsetX), float32(y*cellSize)+float32(gridOffsetY)))
			gridContainer.Add(area)
		}
	}

	// Modify deleteButton's click handler to also hide the banner
	deleteButton := widget.NewButton("Delete all dots", func() {
		for x := range cellStates {
			for y := range cellStates[x] {
				cellStates[x][y] = empty
			}
		}
		dotsContainer.RemoveAll()
		dotsContainer.Refresh()
		clearGrid()
		updateDotCountLabels()
		gameEndBanner.Hide()                 // Hide the banner when all dots are deleted
		timeElapsedLabel.SetText("Time: 0s") // Reset the timer label text
	})

	// Top Layout
	topBar := container.NewHBox(deleteButton, blueDotCountLabel, redDotCountLabel, gridSizeInput, layout.NewSpacer(), timeElapsedLabel)

	// Draw the grid
	drawGrid(gridContainer, float32(gridOffsetX), float32(gridOffsetY), float32(cellSize), gridSize)

	// Layout for the banner to span the entire width
	bannerContainer := container.NewHBox(layout.NewSpacer(), gameEndBanner, layout.NewSpacer())

	// Adjust the content layout to include the banner
	// Use a VBox layout to position the banner in the middle vertically
	content := container.NewVBox(
		container.NewBorder(topBar, nil, nil, nil, container.NewMax(gridContainer, dotsContainer)),
		bannerContainer,
	)

	myWindow.SetContent(content)
}

// Additional methods for handling GameWindow functionalities
// ...
