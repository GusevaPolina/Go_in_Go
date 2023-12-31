package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"image/color"
)

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

// Additional methods for handling dot placing functionalities
// ...
