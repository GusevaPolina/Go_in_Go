package main

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	"time"
)

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

// Additional methods for handling timer functionalities
// ...
