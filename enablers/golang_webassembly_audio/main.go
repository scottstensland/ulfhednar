package main

import (
	"math"
	"syscall/js"
)

func playAudio(this js.Value, args []js.Value) interface{} {
	// Example: Generate a sine wave (replace with your own audio generation logic)
	const sampleRate = 44100
	const duration = 10.0 // in seconds
	// const frequency = 440.0 // Hz (A4)
	const frequency = 40.0 // Hz

	data := make([]byte, int(sampleRate*duration))
	for i := 0; i < len(data); i++ {
		sample := 127.5 * (1 + math.Sin(2*math.Pi*frequency*float64(i)/sampleRate))
		data[i] = byte(sample)
	}

	// Create a JavaScript Uint8Array and copy data
	uint8Array := js.Global().Get("Uint8Array").New(len(data))
	js.CopyBytesToJS(uint8Array, data)

	// Directly call the playAudioData function on the global window object
	js.Global().Call("playAudioData", uint8Array)

	return nil
}

func main() {
	// Export the Go function to JavaScript
	js.Global().Set("goPlayAudio", js.FuncOf(playAudio))

	// Prevent the Go program from exiting
	select {}
}
