package main

import (
	"mustango/cmd"
	"mustango/cmd/noise_spectrogram"
)

func main() {
	cmd.RootCmd.AddCommand(noise_spectrogram.NoiseSpectrogramCmd)
	cmd.Execute()
}
