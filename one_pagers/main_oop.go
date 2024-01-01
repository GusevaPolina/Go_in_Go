package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	beepmp3 "github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"image/color"
	"os"
	"strconv"
	"time"
)

type cellState int

const (
	empty cellState = iota
	red
	blue
)

func (gw *GameWindow) UpdateDotCounters() {
	gw.blueDotCountLabel.SetText(fmt.Sprintf("Blue Dots: %d", *gw.grid.blueDotCount))
	gw.redDotCountLabel.SetText(fmt.Sprintf("Red Dots: %d", *gw.grid.redDotCount))
}

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
	timer         *Timer // Add this field
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

func (g *Grid) PlaceDot(cellColor cellState, x, y int) {
	g.cellStates[x][y] = cellColor
	var dotColor color.Color

	if cellColor == blue {
		dotColor = color.NRGBA{B: 255, A: 255}
		*g.blueDotCount++
	} else {
		dotColor = color.NRGBA{R: 255, A: 255}
		*g.redDotCount++
	}

	dotRadius := float32(g.cellSize) / 5
	dot := canvas.NewCircle(dotColor)
	dot.Resize(fyne.NewSize(dotRadius*2, dotRadius*2))
	dot.Move(fyne.NewPos((float32(x)+0.314)*float32(g.cellSize)+g.gridOffsetX, (float32(y)+0.314)*float32(g.cellSize)+g.gridOffsetY))
	g.dotsContainer.Add(dot)
	g.dotsContainer.Refresh()

	*g.dotCount++

	if g.onDotPlaced != nil {
		g.onDotPlaced()
	}

	// Check and fill clusters only after the second dot is placed
	if *g.dotCount > 1 {
		g.CheckAndFillClusters()
	}
	if *g.blueDotCount+*g.redDotCount == g.gridSize*g.gridSize {
		g.timer.Stop()                    // Stop the timer when all dots are placed
		g.gameWindow.gameEndBanner.Show() // Show the game end banner
	}
}

func (g *Grid) CheckAndFillClusters() {
	visited := make([][]bool, g.gridSize)
	for i := range visited {
		visited[i] = make([]bool, g.gridSize)
	}

	for x := 0; x < g.gridSize; x++ {
		for y := 0; y < g.gridSize; y++ {
			if g.cellStates[x][y] == empty && !visited[x][y] {
				cluster := g.findCluster(x, y, visited)
				borders := g.determineClusterBorders(cluster)
				g.fillClusterIfEnclosed(cluster, borders)
			}
		}
	}
}

func (g *Grid) findCluster(x, y int, visited [][]bool) []struct{ x, y int } {
	if x < 0 || x >= g.gridSize || y < 0 || y >= g.gridSize || visited[x][y] || g.cellStates[x][y] != empty {
		return nil
	}

	visited[x][y] = true
	cluster := []struct{ x, y int }{{x, y}}

	// Recursively search adjacent cells
	directions := []struct{ dx, dy int }{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	for _, dir := range directions {
		cluster = append(cluster, g.findCluster(x+dir.dx, y+dir.dy, visited)...)
	}

	return cluster
}

func (g *Grid) determineClusterBorders(cluster []struct{ x, y int }) map[cellState]bool {
	borders := make(map[cellState]bool)

	for _, cell := range cluster {
		for _, dir := range []struct{ dx, dy int }{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
			nx, ny := cell.x+dir.dx, cell.y+dir.dy

			if nx < 0 || ny < 0 || nx >= g.gridSize || ny >= g.gridSize {
				borders[empty] = true // Mark grid edge as a border
				continue
			}

			if g.cellStates[nx][ny] != empty {
				borders[g.cellStates[nx][ny]] = true
			}
		}
	}

	return borders
}

func (g *Grid) fillClusterIfEnclosed(cluster []struct{ x, y int }, borders map[cellState]bool) {
	// Check if cluster is enclosed by either one color or a combination of one color and grid edges
	if len(borders) == 1 || (len(borders) == 2 && borders[empty]) {
		var fillWith cellState
		colorFound := false

		for colorNow := range borders {
			if colorNow != empty {
				if colorFound {
					return // More than one non-empty color found, do not fill
				}
				fillWith = colorNow
				colorFound = true
			}
		}

		if colorFound {
			dotsFilled := 0
			for _, cell := range cluster {
				if g.cellStates[cell.x][cell.y] == empty {
					g.PlaceDot(fillWith, cell.x, cell.y)
					dotsFilled++
				}
			}
			// Update the dot count based on the number of dots filled
			*g.dotCount += dotsFilled
		}
	}
}

// Additional methods for handling grid functionalities
// ...

// Timer struct for managing game timer.
type Timer struct {
	ticker           *time.Ticker
	startTime        time.Time
	timeElapsedLabel *widget.Label
}

// NewTimer creates a new Timer instance with a label for displaying time.
func NewTimer(label *widget.Label) *Timer {
	return &Timer{
		timeElapsedLabel: label,
	}
}

// Start begins or resumes the timer.
func (t *Timer) Start() {
	if t.ticker != nil {
		return // Timer is already running
	}
	t.startTime = time.Now()
	t.ticker = time.NewTicker(time.Second)

	go func() {
		for range t.ticker.C {
			elapsed := time.Since(t.startTime)
			t.timeElapsedLabel.SetText(fmt.Sprintf("Time: %v", elapsed.Round(time.Second)))
		}
	}()
}

// Stop halts the timer.
func (t *Timer) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
		t.ticker = nil
	}
}

// Reset stops the current timer and starts it anew.
func (t *Timer) Reset() {
	t.Stop()
	t.Start()
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

	//  the timer and dot counters
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

// GameApp manages the application setup and main loop.
type GameApp struct {
	app    fyne.App
	window *GameWindow
}

// NewGameApp creates and initializes a new GameApp.
func NewGameApp() *GameApp {
	fyneApp := app.New()

	gameWindow := NewGameWindow(fyneApp)

	return &GameApp{
		app:    fyneApp,
		window: gameWindow,
	}
}

// Run starts the main loop of the GameApp.
func (app *GameApp) Run() {
	app.window.Configure()
	app.window.Show()
	app.app.Run()
	app.window.Cleanup()
}

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

// Show makes the GameWindow visible.
func (gw *GameWindow) Show() {
	gw.window.Show()
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

// other methods for GameWindow

// MusicPlayer handles background music playing.
type MusicPlayer struct {
	streamer beep.StreamSeekCloser // The audio stream
	format   beep.Format           // The audio format
	ctrl     *beep.Ctrl            // Control for pausing/resuming
	done     chan bool             // Channel to signal when playback is finished
}

// NewMusicPlayer creates a new music player instance.
func NewMusicPlayer(filename string) *MusicPlayer {
	f, err := os.Open(filename)
	if err != nil {
		// Handle error (e.g., file not found)
		return nil
	}

	streamer, format, err := beepmp3.Decode(f)
	if err != nil {
		// Handle error (e.g., decoding error)
		return nil
	}

	return &MusicPlayer{
		streamer: streamer,
		format:   format,
		ctrl:     &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false},
		done:     make(chan bool),
	}
}

// Play starts playing the background music.
func (mp *MusicPlayer) Play() {
	speaker.Init(mp.format.SampleRate, mp.format.SampleRate.N(time.Second/10))
	speaker.Play(mp.ctrl)

	go func() {
		<-mp.done
	}()
}

// Stop halts the music playback.
func (mp *MusicPlayer) Stop() {
	if mp.streamer != nil {
		mp.ctrl.Paused = true
		mp.streamer.Close()
	}
}

// Resume resumes the music playback if it was paused.
func (mp *MusicPlayer) Resume() {
	if mp.streamer != nil {
		mp.ctrl.Paused = false
	}
}

func main() {
	gameApp := NewGameApp()
	gameApp.Run()
}
