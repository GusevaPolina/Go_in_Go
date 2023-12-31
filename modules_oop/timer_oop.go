package main

import (
	"fmt"
	"fyne.io/fyne/v2/widget"
	"time"
)

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

// Additional methods for handling timer functionalities
// ...
