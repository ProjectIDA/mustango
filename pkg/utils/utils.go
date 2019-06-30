package utils

import "strings"

func ProcessArgs(args []string) (string, string, string, string, string, string) {
	networks := strings.ToUpper(args[0])
	stations := strings.ToUpper(args[1])
	locations := strings.ToUpper(args[2])
	channels := strings.ToUpper(args[3])

	startdate := ""
	enddate := ""

	return networks, stations, locations, channels, startdate, enddate
}

