package main

import (
	"fmt"
	"github.com/faiface/beep"
	beepmp3 "github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"os"
	"time"
)

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
		defer func(streamer beep.StreamSeekCloser) {
			err := streamer.Close()
			if err != nil {

			}
		}(streamer)

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		done := make(chan bool)
		speaker.Play(beep.Seq(streamer, beep.Callback(func() {
			done <- true
		})))

		// This will block until the streamer is done playing.
		<-done
	}()
}

// Additional methods for handling music player functionalities
// ...
