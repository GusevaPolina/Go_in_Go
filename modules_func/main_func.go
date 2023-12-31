package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(myTheme{theme.LightTheme()})
	myWindow := myApp.NewWindow("Go in Go: the coolest version")

	// Start playing background music
	playBackgroundMusic("background.mp3")

	// Grid size input field with default value
	gridSizeInput := widget.NewEntry()
	gridSizeInput.SetText(fmt.Sprintf("%d", gridSize)) // Set default grid size
	gridSizeInput.OnSubmitted = func(value string) {
		newGridSize, err := strconv.Atoi(value)
		if err != nil || newGridSize < 1 {
			dialog.ShowError(fmt.Errorf("please, enter a valid grid size"), myWindow)
			gridSizeInput.SetText(fmt.Sprintf("%d", gridSize)) // Reset to current grid size
			return
		}

		// Save the new grid size and update the UI
		gridSize = newGridSize
		regenerateGrid(gridSize, myWindow, gridSizeInput)
	}
	regenerateGrid(gridSize, myWindow, gridSizeInput)

	myWindow.SetOnClosed(func() {
		stopTimer() // Stop the timer when the window is closed
	})
	myWindow.Resize(fyne.NewSize(windowWidth, windowHeight*1.1))
	myWindow.ShowAndRun()
}
