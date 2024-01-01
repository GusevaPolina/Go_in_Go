package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strconv"
)

// GameWindow represents the main game window.
type GameWindow struct {
	window       fyne.Window
	grid         *Grid
	timer        *Timer
	musicPlayer  *MusicPlayer
	windowWidth  int
	windowHeight int
	gridSize     int
	// UI components
	timeElapsedLabel  *widget.Label
	blueDotCountLabel *widget.Label
	redDotCountLabel  *widget.Label
	gridSizeInput     *widget.Entry
	backgroundImage   *canvas.Image
	gameEndBanner     *fyne.Container
}

func NewGameWindow(app fyne.App) *GameWindow {
	mainWindow := app.NewWindow("Go in Go: the coolest version")

	gw := &GameWindow{
		window:       mainWindow,
		windowWidth:  500,
		windowHeight: 500,
		gridSize:     9,
		// Initialize labels
		blueDotCountLabel: widget.NewLabel("Blue Dots: 0"),
		redDotCountLabel:  widget.NewLabel("Red Dots: 0"),
		timeElapsedLabel:  widget.NewLabel("Time: 0s"),
		// Initialize the gameEndBanner
		gameEndBanner: createGameEndBanner(),
		// Initialize gridSizeInput
		gridSizeInput: widget.NewEntry(),
	}

	// Initialize Timer
	gw.timer = NewTimer(gw.timeElapsedLabel)

	// Initialize MusicPlayer
	gw.musicPlayer = NewMusicPlayer("../background.mp3")

	// Configure gridSizeInput
	gw.gridSizeInput.SetText(fmt.Sprintf("%d", gw.gridSize))
	gw.gridSizeInput.OnSubmitted = func(value string) {
		newGridSize, err := strconv.Atoi(value)
		if err != nil || newGridSize < 1 {
			dialog.ShowError(fmt.Errorf("Invalid grid size"), gw.window)
			gw.gridSizeInput.SetText(fmt.Sprintf("%d", gw.gridSize)) // Reset to current grid size
			return
		}
		gw.RegenerateGrid(newGridSize)
	}

	// Call RegenerateGrid here after initializing all components
	gw.RegenerateGrid(gw.gridSize)

	return gw
}

// Show makes the GameWindow visible.
func (gw *GameWindow) Show() {
	gw.window.Show()
}

func (gw *GameWindow) Configure() {
	gw.window.SetTitle("Go in Go: the coolest version")
	gw.window.Resize(fyne.NewSize(float32(gw.windowWidth), float32(gw.windowHeight)*1.1))

	// Initialize Music Player
	gw.musicPlayer = NewMusicPlayer("../background.mp3") // Assuming you have a NewMusicPlayer function
	gw.musicPlayer.Play()                             // Start music (if needed)

	// Initialize other UI components
	gw.blueDotCountLabel = widget.NewLabel("Blue Dots: 0")
	gw.redDotCountLabel = widget.NewLabel("Red Dots: 0")

	// Initialize Timer and its label
	gw.timeElapsedLabel = widget.NewLabel("Time: 0s")
	gw.timer = NewTimer(gw.timeElapsedLabel)
	gw.timer.Start() // Start the timer

	gw.gridSizeInput = widget.NewEntry()
	gw.gridSizeInput.SetText(fmt.Sprintf("%d", gw.gridSize))
	gw.gridSizeInput.OnSubmitted = func(value string) {
		newGridSize, err := strconv.Atoi(value)
		if err != nil || newGridSize < 1 {
			dialog.ShowError(fmt.Errorf("Invalid grid size"), gw.window)
			gw.gridSizeInput.SetText(fmt.Sprintf("%d", gw.gridSize)) // Reset to current grid size
			return
		}
		gw.RegenerateGrid(newGridSize)
	}

	gw.RegenerateGrid(gw.gridSize)
}

func (gw *GameWindow) UpdateDotCounters() {
	gw.blueDotCountLabel.SetText(fmt.Sprintf("Blue Dots: %d", *gw.grid.blueDotCount))
	gw.redDotCountLabel.SetText(fmt.Sprintf("Red Dots: %d", *gw.grid.redDotCount))
}

func createGameEndBanner() *fyne.Container {
	// Define the text style for the regular and highlighted parts
	regularTextStyle := fyne.TextStyle{Bold: true}
	highlightedTextStyle := fyne.TextStyle{Bold: true, Italic: true}
	// Define the colors
	regularColor := color.Black
	highlightedColor := color.NRGBA{R: 85, G: 208, B: 225, A: 255} // Blue color
	// Define the fonts
	regularFont := float32(24)
	highlightedFont := float32(32)

	// Create separate text objects for each part of the banner
	textG := canvas.NewText("G", highlightedColor)
	textG.TextStyle, textG.TextSize = highlightedTextStyle, highlightedFont
	textGame := canvas.NewText("ame ", regularColor)
	textGame.TextStyle, textGame.TextSize = regularTextStyle, regularFont

	textO := canvas.NewText("O", highlightedColor)
	textO.TextStyle, textO.TextSize = highlightedTextStyle, highlightedFont
	textOver := canvas.NewText("ver!", regularColor)
	textOver.TextStyle, textOver.TextSize = regularTextStyle, regularFont

	// Add the text objects to the container
	bannerContainer := container.NewHBox(layout.NewSpacer(), textG, textGame, textO, textOver, layout.NewSpacer())

	return bannerContainer
}

func (gw *GameWindow) RegenerateGrid(newGridSize int) {
	gw.gridSize = newGridSize

	// Reset the timer and dot counters
	gw.timer.Reset()
	gw.blueDotCountLabel.SetText("Blue Dots: 0")
	gw.redDotCountLabel.SetText("Red Dots: 0")
	gw.gameEndBanner.Hide() // Hide the game end banner

	// Initialize and draw the new grid
	gw.grid = NewGrid(gw.gridSize, 400/gw.gridSize, float32(gw.windowWidth)*0.1, float32(gw.windowHeight)*0.1, gw.UpdateDotCounters, gw.timer, gw)
	gw.grid.DrawGrid()
	gw.grid.onDotPlaced = gw.UpdateDotCounters

	resetButton := widget.NewButton("Go try again", func() {
		// Reset the grid
		for x := range gw.grid.cellStates {
			for y := range gw.grid.cellStates[x] {
				gw.grid.cellStates[x][y] = empty
			}
		}
		gw.grid.dotsContainer.RemoveAll()
		gw.grid.dotsContainer.Refresh()

		// Reset dot counters
		*gw.grid.blueDotCount = 0
		*gw.grid.redDotCount = 0
		*gw.grid.dotCount = 0
		gw.blueDotCountLabel.SetText("Blue Dots: 0")
		gw.redDotCountLabel.SetText("Red Dots: 0")

		// Reset the timer
		gw.timer.Reset()

		// Optionally, if you have game over or other status indicators, reset them
		// ...

		// Hide the game end banner when the reset button is clicked
		gw.gameEndBanner.Hide()

		// Optionally, redraw the grid if needed
		gw.grid.DrawGrid()
	})

	// Load and set the background image
	gw.backgroundImage = canvas.NewImageFromFile("../background.png")
	gw.backgroundImage.FillMode = canvas.ImageFillContain
	gw.backgroundImage.Resize(fyne.NewSize(float32(gw.windowWidth), float32(gw.windowHeight)))
	gw.backgroundImage.Move(fyne.NewPos(0, 0))

	// Hide the game end banner only if it's initialized
	if gw.gameEndBanner != nil {
		gw.gameEndBanner.Hide()
	}

	// Update the main container to include the new grid
	mainContainer := container.NewWithoutLayout()
	mainContainer.Add(gw.backgroundImage)
	mainContainer.Add(gw.grid.container)
	mainContainer.Add(gw.grid.dotsContainer) // Make sure dotsContainer is part of the main content

	// Layout for the top bar with the timer at the right of the dot counters
	topBar := container.NewHBox(
		resetButton,
		gw.blueDotCountLabel,
		gw.redDotCountLabel,
		gw.gridSizeInput,
		layout.NewSpacer(),  // Spacer pushes the timer to the right
		gw.timeElapsedLabel, // Timer label on the far right
	)

	// Combine the top bar with the main container
	// Use a VBox layout to position the banner in the middle vertically
	content := container.NewVBox(
		container.NewBorder(topBar, nil, nil, nil, mainContainer),
		gw.gameEndBanner,
	)

	// Set the content of the window
	gw.window.SetContent(content)
}

// Cleanup performs any necessary cleanup tasks for the GameWindow.
func (gw *GameWindow) Cleanup() {
	if gw.timer != nil {
		gw.timer.Stop()
	}
	if gw.musicPlayer != nil {
		gw.musicPlayer.Stop()
	}
	// Additional cleanup for other resources
	// ...
}

// Additional methods for handling GameWindow functionalities
// ...
