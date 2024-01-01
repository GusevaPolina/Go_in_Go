package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	beepmp3 "github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"image/color"
	"os"
	"time"
)

const (
	gridSize = 9
	cellSize = 40
)

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

func newTappableArea(x, y int, gridOffsetX, gridOffsetY float32, onTap func(x, y int)) *tappableArea {
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

func placeDot(cellColor cellState, x, y int, cellStates [][]cellState, dotsContainer *fyne.Container, gridOffsetX, gridOffsetY float32, dotCount *int) {
	cellStates[x][y] = cellColor
	var dotColor color.Color

	if cellColor == blue {
		dotColor = color.NRGBA{B: 255, A: 255} // Blue dot
	} else {
		dotColor = color.NRGBA{R: 255, A: 255} // Red dot
	}

	// Create a new circle with the dot color
	dot := canvas.NewCircle(dotColor)
	// Resize the dot to be twice as big as before
	dot.Resize(fyne.NewSize(cellSize/5, cellSize/5))
	// Move the dot to the correct position, adjusting for the new size
	dot.Move(fyne.NewPos(float32(x)*cellSize+cellSize/2-cellSize/10+gridOffsetX, float32(y)*cellSize+cellSize/2-cellSize/10+gridOffsetY))
	// Add the dot to the dotsContainer
	dotsContainer.Add(dot)
	// Refresh the container to update the display
	dotsContainer.Refresh()

	if cellColor == blue {
		blueDotCount++
	} else {
		redDotCount++
	}

	*dotCount++
}

func findCluster(cellStates [][]cellState, x, y int, visited [][]bool) []struct{ x, y int } {
	if x < 0 || x >= gridSize || y < 0 || y >= gridSize || visited[x][y] || cellStates[x][y] != empty {
		return nil
	}

	visited[x][y] = true
	cluster := []struct{ x, y int }{{x, y}}

	cluster = append(cluster, findCluster(cellStates, x+1, y, visited)...)
	cluster = append(cluster, findCluster(cellStates, x-1, y, visited)...)
	cluster = append(cluster, findCluster(cellStates, x, y+1, visited)...)
	cluster = append(cluster, findCluster(cellStates, x, y-1, visited)...)

	return cluster
}

func determineClusterBorders(cellStates [][]cellState, cluster []struct{ x, y int }) map[cellState]bool {
	borders := make(map[cellState]bool)

	for _, cell := range cluster {
		for _, dir := range []struct{ dx, dy int }{{0, -1}, {1, 0}, {0, 1}, {-1, 0}} {
			nx, ny := cell.x+dir.dx, cell.y+dir.dy

			if nx < 0 || ny < 0 || nx >= gridSize || ny >= gridSize {
				borders[empty] = true // Mark grid edge as a border
				continue
			}

			if cellStates[nx][ny] != empty {
				borders[cellStates[nx][ny]] = true
			}
		}
	}

	return borders
}

func fillClusterIfEnclosed(cellStates [][]cellState, cluster []struct{ x, y int }, borders map[cellState]bool, dotsContainer *fyne.Container, gridOffsetX, gridOffsetY float32) {
	if len(borders) == 1 || (len(borders) == 2 && borders[empty]) {
		var fillWith cellState
		colorFound := false

		for color := range borders {
			if color != empty {
				if colorFound {
					return // More than one color found, do not fill
				}
				fillWith = color
				colorFound = true
			}
		}

		if colorFound {
			for _, cell := range cluster {
				placeDot(fillWith, cell.x, cell.y, cellStates, dotsContainer, gridOffsetX, gridOffsetY, new(int)) // Dummy counter
			}
		}
	}
}

func checkAndFillClusters(cellStates [][]cellState, dotsContainer *fyne.Container, gridOffsetX, gridOffsetY float32) {
	visited := make([][]bool, gridSize)
	for i := range visited {
		visited[i] = make([]bool, gridSize)
	}

	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			if cellStates[x][y] == empty && !visited[x][y] {
				cluster := findCluster(cellStates, x, y, visited)
				borders := determineClusterBorders(cellStates, cluster)
				fillClusterIfEnclosed(cellStates, cluster, borders, dotsContainer, gridOffsetX, gridOffsetY)
			}
		}
	}
}

func drawGrid(container *fyne.Container, offsetX, offsetY float32, cellSize, gridSize int) {
	// Adjust the starting point based on the radius of the dots
	adjustedOffsetX := offsetX + float32(cellSize)/2
	adjustedOffsetY := offsetY + float32(cellSize)/2

	for i := 0; i < gridSize; i++ {
		// Draw vertical lines
		vLine := canvas.NewLine(color.Black)
		// vLine.StrokeWidth = 2
		vLine.Move(fyne.NewPos(adjustedOffsetX+float32(i)*float32(cellSize), adjustedOffsetY))
		vLine.Resize(fyne.NewSize(2, float32(cellSize)*(float32(gridSize)-1)))
		container.Add(vLine)

		// Draw horizontal lines
		hLine := canvas.NewLine(color.Black)
		// hLine.StrokeWidth = 2
		hLine.Move(fyne.NewPos(adjustedOffsetX, adjustedOffsetY+float32(i)*float32(cellSize)))
		hLine.Resize(fyne.NewSize(float32(cellSize)*(float32(gridSize)-1), 2))
		container.Add(hLine)
	}
}

// Timer variables
var startTime time.Time
var timer *time.Ticker
var timeElapsedLabel *widget.Label

// Starts or resets the timer
func resetTimer() {
	startTime = time.Now()
	if timer != nil {
		timer.Stop()
	}
	timer = time.NewTicker(time.Second)
	go func() {
		for range timer.C {
			elapsed := time.Since(startTime)
			timeElapsedLabel.SetText(fmt.Sprintf("Time: %v", elapsed.Round(time.Second)))
		}
	}()
}

// Stops the timer
func stopTimer() {
	if timer != nil {
		timer.Stop()
		timer = nil
	}
}

func playBackgroundMusic(filename string) {
	go func() {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}

		streamer, format, err := beepmp3.Decode(f)
		if err != nil {
			fmt.Println("Error decoding MP3:", err)
			return
		}
		defer streamer.Close()

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		done := make(chan bool)
		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			done <- true
		})))

		// This will block until the streamer is done playing.
		<-done
	}()
}

func main() {
	myApp := app.New()
	myApp.Settings().SetTheme(myTheme{theme.LightTheme()})

	myWindow := myApp.NewWindow("9x9 Grid with Advanced Placement Rules")

	gridContainer := container.NewWithoutLayout()
	dotsContainer := container.NewWithoutLayout()

	windowWidth := float32(gridSize*cellSize) / 0.8
	windowHeight := float32(gridSize*cellSize) / 0.8

	backgroundImage := canvas.NewImageFromFile("../image.png")
	backgroundImage.FillMode = canvas.ImageFillContain
	backgroundImage.Resize(fyne.NewSize(windowWidth, windowHeight))
	backgroundImage.Move(fyne.NewPos(0, 0))
	gridContainer.Add(backgroundImage)

	gridOffsetX := windowWidth * 0.1
	gridOffsetY := windowHeight * 0.1

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
			area := newTappableArea(x, y, gridOffsetX, gridOffsetY, func(x, y int) {
				if cellStates[x][y] == empty {
					currentColor := red
					if dotCount%2 == 0 {
						currentColor = blue
					}
					placeDot(currentColor, x, y, cellStates, dotsContainer, gridOffsetX, gridOffsetY, &dotCount)
					if dotCount > 1 { // Skip the first step
						checkAndFillClusters(cellStates, dotsContainer, gridOffsetX, gridOffsetY)
					}
					updateDotCountLabels()
				}
				checkAndUpdateGameEnd() // Update banner visibility
			})
			area.Resize(fyne.NewSize(cellSize, cellSize))
			area.Move(fyne.NewPos(float32(x)*cellSize+gridOffsetX, float32(y)*cellSize+gridOffsetY))
			gridContainer.Add(area)
		}
	}

	// Modify deleteButton's click handler to also hide the banner
	deleteButton := widget.NewButton("Delete All Dots", func() {
		for x := range cellStates {
			for y := range cellStates[x] {
				cellStates[x][y] = empty
			}
		}
		dotsContainer.RemoveAll()
		dotsContainer.Refresh()
		dotCount = 0
		blueDotCount = 0
		redDotCount = 0
		updateDotCountLabels()
		gameEndBanner.Hide()                 // Hide the banner when all dots are deleted
		timeElapsedLabel.SetText("Time: 0s") // Reset the timer label text
		resetTimer()
	})

	// Reset the timer when the app starts
	resetTimer()

	// Layout with the timer, banner, and dot counters
	topBar := container.NewHBox(deleteButton, blueDotCountLabel, redDotCountLabel, layout.NewSpacer(), timeElapsedLabel)

	// Draw the grid
	drawGrid(gridContainer, gridOffsetX, gridOffsetY, cellSize, gridSize)

	// Start playing background music
	playBackgroundMusic("../background.mp3")

	// Layout for the banner to span the entire width
	bannerContainer := container.NewHBox(layout.NewSpacer(), gameEndBanner, layout.NewSpacer())

	// Adjust the content layout to include the banner
	// Use a VBox layout to position the banner in the middle vertically
	content := container.NewVBox(
		container.NewBorder(topBar, nil, nil, nil, container.NewMax(gridContainer, dotsContainer)),
		bannerContainer,
	)

	myWindow.SetOnClosed(func() {
		stopTimer() // Stop the timer when the window is closed
	})
	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(windowWidth, windowHeight))
	myWindow.ShowAndRun()
}
