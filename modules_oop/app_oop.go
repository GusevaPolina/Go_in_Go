package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

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

// Additional methods for handling GameApp functionalities
// ...
