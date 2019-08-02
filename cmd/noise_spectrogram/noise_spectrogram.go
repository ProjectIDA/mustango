package noise_spectrogram

import (
	"fmt"
	"log"
	"mustango/pkg/mustang"
	"mustango/pkg/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var NoiseSpectrogramCmd = &cobra.Command{
	Version: "0.9.1",
	Use:   "noise-spectrogram net sta loclist chnlist name_tail",
	Short: "Retrieve noise spectrograms from IRIS MUSTANG for the specified net, sta, chnlist and loclist.", ///cmdName + "net sta loclist chanlist"
	Long:  `
Retrieve Noise-Spectrogram plots from the IRIS MUSTANG 'noise-spectrogram' Web Service.

See the IRIS MUSTANG Service documentation at http://service.iris.edu/mustang/noise-spectrogram/1/
for more information on all options of the options available for this service`,
	Args:  cobra.ExactArgs(4), // net sta loc[,loc,...] chn[,chn,...]
	Run: func(cmd *cobra.Command, args []string) {

		processCmd(cmd, "noise-spectrogram", "1", args)
	},
}

func init() {
	NoiseSpectrogramCmd.Flags().String("name-tail", "", "tail of output image file name")
	NoiseSpectrogramCmd.Flags().String("format", "plot", "Output format: 'plot'")
	NoiseSpectrogramCmd.Flags().String("output", "power", "Plot output ['power'|'powerdhnm'|'powerdlnm'|'powerdnm'|'powerdmedian']")
	NoiseSpectrogramCmd.Flags().String("plot.height", "", "Height in pixels (default '1000').\nNote: If only Height or Width specified, the other scales proportionately.")
	NoiseSpectrogramCmd.Flags().String("plot.width", "", "Height in pixels (default '2000').\nNote: If only Height or Width specified, the other scales proportionately.")
	NoiseSpectrogramCmd.Flags().String("plot.horzaxis", "time", "Horizontal Axis ['time'|'freq']")
	NoiseSpectrogramCmd.Flags().String("plot.title", "", "Plot Title text or 'hide' (default SNLCQ target)")
	NoiseSpectrogramCmd.Flags().String("plot.titlefont.size", "", "Plot Title font size in points (default 'auto')")
	NoiseSpectrogramCmd.Flags().String("plot.subtitle", "", "Plot Subtitle text or 'hide' (default 'hide')")
	NoiseSpectrogramCmd.Flags().String("plot.subtitlefont.size", "", "Plot Subtitle font size in points (default 'auto')")
	NoiseSpectrogramCmd.Flags().String("plot.labelfont.size", "", "Plot Label font size in points (default 'auto')")
	NoiseSpectrogramCmd.Flags().String("plot.frequency.label", "", "Plot Frequency Axis label text or 'hide' (default 'Frequency (Hz)')")
	NoiseSpectrogramCmd.Flags().String("plot.frequency.invert", "", "Plot Frequency Axis Invert ['true'|'false'] (default 'false')")
	NoiseSpectrogramCmd.Flags().String("plot.frequency.range", "", "Plot Frequency Axis Range <min,max> (default full freq of data)")
	NoiseSpectrogramCmd.Flags().String("plot.time.format", "", "Plot Time Axis format (default 'yyyy-MM-dd')")
	NoiseSpectrogramCmd.Flags().String("plot.time.label", "", "Plot Time Axis label text or hide (default <DATE RANGE>)")
	NoiseSpectrogramCmd.Flags().String("plot.time.matchrequest", "true", "Plot Time Match Request (plot time axis as requested or as available in the data ['true'|'false']")
	NoiseSpectrogramCmd.Flags().String("plot.time.tickunit", "", "Plot Time Tick Unit ['auto'|'day'|'month'|'year'] (default 'auto')")
	NoiseSpectrogramCmd.Flags().String("plot.time.invert", "", "Plot Time Axis Invert. ['true'|'false'] (default 'false')")
	NoiseSpectrogramCmd.Flags().String("plot.powerscale.autorange", "", "Plot Power Scale Autorange (float: 0.0 - 1.0) (default 0.95)")
	NoiseSpectrogramCmd.Flags().String("plot.powerscale.range", "", "Plot Power Scale Range <min,max> (Autorange and Range or mutually exclusive options)")
	NoiseSpectrogramCmd.Flags().String("plot.powerscale.show", "", "Plot Power Scale Show ['true'|'false'] (default 'true')")
	NoiseSpectrogramCmd.Flags().String("plot.powerscale.x", "", "Plot Power Scale X position in pixels (default 5)")
	NoiseSpectrogramCmd.Flags().String("plot.powerscale.y", "", "Plot Power Scale Y position in pixels (default 5)")
	NoiseSpectrogramCmd.Flags().String("plot.powerscale.height", "", "Plot Power Scale Height in pixels (default 12)")
	NoiseSpectrogramCmd.Flags().String("plot.powerscale.width", "", "Plot Power Scale Width in pixels (default 150)")
	NoiseSpectrogramCmd.Flags().String("plot.powerscale.orientation", "", "Plot Powerscale Orientation ['horz'|'vert'] (default 'horz')")
	NoiseSpectrogramCmd.Flags().String("noisemode.byperiod", "", "noisemode.byperiod (pipe-seprated list)")
	NoiseSpectrogramCmd.Flags().String("noisemode.byfrequency", "", "noisemode.byfrequency (pipe-seprated list)")
}

func processFlags(cc *cobra.Command, opts map[string]string) {

	if v, err := cc.Flags().GetString("outDIR"); err == nil {
		opts["outDIR"] = v
	}
	if v, err := cc.Flags().GetString("format"); err == nil {
		opts["format"] = v
	}
	if v, err := cc.Flags().GetString("output"); (err == nil) && (v != "") {
		opts["output"] = v
	}
	if v, err := cc.Flags().GetString("plot.height"); (err == nil) && (v != "") {
		opts["plot.height"] = v
	}
	if v, err := cc.Flags().GetString("plot.width"); (err == nil) && (v != "") {
		opts["plot.width"] = v
	}
	if v, err := cc.Flags().GetString("plot.horzaxis"); (err == nil) && (v != "") {
		opts["plot.horzaxis"] = v
	}
	if v, err := cc.Flags().GetString("plot.title"); (err == nil) && (v != "") {
		opts["plot.title"] = v
	}
	if v, err := cc.Flags().GetString("plot.titlefont.size"); (err == nil) && (v != "") {
		opts["plot.titlefont.size"] = v
	}
	if v, err := cc.Flags().GetString("plot.subtitle"); (err == nil) && (v != "") {
		opts["plot.subtitle"] = v
	}
	if v, err := cc.Flags().GetString("plot.subtitlefont.size"); (err == nil) && (v != "") {
		opts["plot.subtitlefont.size"] = v
	}
	if v, err := cc.Flags().GetString("plot.axisfont.size"); (err == nil) && (v != "") {
		opts["plot.axisfont.size"] = v
	}
	if v, err := cc.Flags().GetString("plot.labelfont.size"); (err == nil) && (v != "") {
		opts["plot.labelfont.size"] = v
	}
	if v, err := cc.Flags().GetString("plot.frequency.label"); (err == nil) && (v != "") {
		opts["plot.frequency.label"] = v
	}
	if v, err := cc.Flags().GetString("plot.frequency.invert"); (err == nil) && (v != "") {
		opts["plot.frequency.invert"] = v
	}
	if v, err := cc.Flags().GetString("plot.frequency.range"); (err == nil) && (v != "") {
		opts["plot.frequency.range"] = v
	}
	if v, err := cc.Flags().GetString("plot.time.format"); (err == nil) && (v != "") {
		opts["plot.time.format"] = v
	}
	if v, err := cc.Flags().GetString("plot.time.label"); (err == nil) && (v != "") {
		opts["plot.time.label"] = v
	}
	if v, err := cc.Flags().GetString("plot.time.matchrequest"); (err == nil) && (v != "") {
		opts["plot.time.matchrequest"] = v
	}
	if v, err := cc.Flags().GetString("plot.time.tickunit"); (err == nil) && (v != "") {
		opts["plot.time.tickunit"] = v
	}
	if v, err := cc.Flags().GetString("plot.time.invert"); (err == nil) && (v != "") {
		opts["plot.time.invert"] = v
	}
	if v, err := cc.Flags().GetString("plot.powerscale.autorange"); (err == nil) && (v != "") {
		opts["plot.powerscale.autorange"] = v
	}
	if v, err := cc.Flags().GetString("plot.powerscale.range"); (err == nil) && (v != "") {
		opts["plot.powerscale.range"] = v
	}
	if v, err := cc.Flags().GetString("plot.powerscale.show"); (err == nil) && (v != "") {
		opts["plot.powerscale.show"] = v
	}
	if v, err := cc.Flags().GetString("plot.powerscale.x"); (err == nil) && (v != "") {
		opts["plot.powerscale.x"] = v
	}
	if v, err := cc.Flags().GetString("plot.powerscale.y"); (err == nil) && (v != "") {
		opts["plot.powerscale.y"] = v
	}
	if v, err := cc.Flags().GetString("plot.powerscale.height"); (err == nil) && (v != "") {
		opts["plot.powerscale.height"] = v
	}
	if v, err := cc.Flags().GetString("plot.powerscale.width"); (err == nil) && (v != "") {
		opts["plot.powerscale.width"] = v
	}
	if v, err := cc.Flags().GetString("plot.powerscale.orientation"); (err == nil) && (v != "") {
		opts["plot.powerscale.orientation"] = v
	}
	if v, err := cc.Flags().GetString("noisemode.byperiod"); (err == nil) && (v != "") {
		opts["noisemode.byperiod"] = v
	}
	if v, err := cc.Flags().GetString("noisemode.byfrequency"); (err == nil) && (v != "") {
		opts["noisemode.byfrequency"] = v
	}
	opts["nodata"] = "404"
}

func processCmd(cc *cobra.Command, svcname, apiversion string, args []string) {

	//BEGIN these params to come from args or caller
	networks, stations, locations, channels, _, _ := utils.ProcessArgs(args)
	quals := "*"

	now := time.Now().UTC()
	endt := now.Truncate(time.Duration(time.Hour * 24))
	endY, endM, endD := endt.Date()
	startt := endt.Add(-time.Duration(365 * 24 * time.Hour))
	startY, startM, startD := startt.Date()
	beginDate := fmt.Sprintf("%4d-%02d-%02d", startY, startM, startD)
	endDate := fmt.Sprintf("%4d-%02d-%02d", endY, endM, endD)

	nameTail, _ := cc.Flags().GetString("name-tail")

	opts := map[string]string{}
	processFlags(cc, opts)

	outLocation, _ := cc.Flags().GetString("outdir")

	if _, err := os.Stat(outLocation); os.IsNotExist(err) {
		err := os.MkdirAll(outLocation, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	outAbsPath, _ := filepath.Abs(outLocation)
	outdir := filepath.Join(outAbsPath, svcname)

	requests := mustang.RequestList(svcname, apiversion, networks, stations, locations, channels, quals, beginDate, endDate, opts)

	results := make(chan *mustang.Result, len(requests))

	go mustang.MakeRequests(requests, opts, results)

	for result := range results {

		filename := filepath.Join(
			outdir,
			strings.ToLower(result.Req.Sta),
			result.FileName(nameTail))
		err := result.SaveToFile(filename)
		if err != nil {
			fmt.Printf("Error saving image to %s: %s\n", filename, err)
		}
	}

}
