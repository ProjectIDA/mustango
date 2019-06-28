package main

import (
	"mustango/cmd"
	"mustango/cmd/noise_spectrograms"
)

func main() {
	cmd.RootCmd.AddCommand(noise_spectrograms.NoiseSpectrogramsCmd)
	//cmd.RootCmd.AddCommand(noise_mode_timeseries.NoiseModeTimeSeriesCmd)
	cmd.Execute()
}
