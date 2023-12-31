package main

import (
	"github.com/faiface/beep"
	beepmp3 "github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"os"
	"time"
)

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

// Additional methods for handling music player functionalities
// ...
