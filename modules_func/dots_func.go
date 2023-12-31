package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image/color"
)

func placeDot(cellColor cellState, x, y int, cellStates [][]cellState, dotsContainer *fyne.Container, cellSize, gridOffsetX, gridOffsetY float32, dotCount *int) {
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

func fillClusterIfEnclosed(cellStates [][]cellState, cluster []struct{ x, y int }, borders map[cellState]bool, dotsContainer *fyne.Container, cellSize int, gridOffsetX, gridOffsetY float32) {
	if len(borders) == 1 || (len(borders) == 2 && borders[empty]) {
		var fillWith cellState
		colorFound := false

		for colorNow := range borders {
			if colorNow != empty {
				if colorFound {
					return // More than one color found, do not fill
				}
				fillWith = colorNow
				colorFound = true
			}
		}

		if colorFound {
			for _, cell := range cluster {
				placeDot(fillWith, cell.x, cell.y, cellStates, dotsContainer, float32(cellSize), gridOffsetX, gridOffsetY, new(int)) // Dummy counter
			}
		}
	}
}

func checkAndFillClusters(cellStates [][]cellState, dotsContainer *fyne.Container, cellSize int, gridOffsetX, gridOffsetY float32) {
	visited := make([][]bool, gridSize)
	for i := range visited {
		visited[i] = make([]bool, gridSize)
	}

	for x := 0; x < gridSize; x++ {
		for y := 0; y < gridSize; y++ {
			if cellStates[x][y] == empty && !visited[x][y] {
				cluster := findCluster(cellStates, x, y, visited)
				borders := determineClusterBorders(cellStates, cluster)
				fillClusterIfEnclosed(cellStates, cluster, borders, dotsContainer, cellSize, gridOffsetX, gridOffsetY)
			}
		}
	}
}

// Additional methods for handling dot placing functionalities
// ...
