package main

import (
	"mustango/cmd"
	"mustango/cmd/noise_spectrograms"
)

func main() {
	cmd.RootCmd.AddCommand(noise_spectrograms.NoiseSpectrogramsCmd)
	cmd.Execute()
}
