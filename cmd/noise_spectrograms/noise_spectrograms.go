package noise_spectrograms

//net = <network>
//sta = <station>
//loc = <location>
//cha = <channel>  OR  target = <nslcq>
//
//=  =  =  =  =  =  =  =  =  =  =  =  =  =  =  =
//OPTIONAL
//
//start = <day yyyy-mm-dd>
//end = <day yyyy-mm-dd>
//
//output = <power|powerdhnm|powerdlnm|powerdnm|powerdmedian>
//noisemodel_byperiod = <pipe-separated-list>
//noisemodel_byfrequency = <pipe-separated-list>
//
//plot_height = <pixel-height>
//plot_width = <pixel-width>
//plot_horzaxis = <time|freq|frequency>
//
//plot_powerscale_autorange = <value 0 to 1>|plot_powerscale_range = <min,max>
//plot_powerscale_height = <value>
// plot_powerscale_width = <value>
//plot_powerscale_show = <true|false>
// plot_powerscale_x = <value>
//plot_powerscale_y = <value>
//plot_powerscale_orientation = <horizontal|horz|vertical|vert>
//
//plot_title = <title|hide>
//plot_subtitle = <title|hide>
//plot_titlefont_size = <value>
//plot_subtitlefont_size = <value>
//plot_axisfont_size = <value>
//plot_labelfont_size = <value>
//
//plot_time_matchrequest = <true|false>
// plot_time_format = <time-format>
// plot_time_label = <text|hide>
//plot_time_tickunit = <auto|day|month|year>
// plot_time_invert = <true|false>
//
//plot_frequency_label = <text|hide>
// plot_frequency_invert = <true|false>
//plot_frequency_range = <min,max>

import (
	"fmt"
	"mustango/pkg/mustang"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var NoiseSpectrogramsCmd = &cobra.Command{
	Use:   "noise-spectrogram",
	Short: "Retrieve Noise-Spectrogram plots (png) from IRIS MUSTANG 'noise-spectrograms' Web Service",
	Run: func(cmd *cobra.Command, args []string) {
		processCmd("noise-spectrogram", "1", args)
	},
}

//type noiseSpectrogramOptions struct {
//	//BeginDate                   string // empty or <YYYY-MM-DD> [oldest data]
//	//EndDate                     string // empty or <YYYY-MM-DD> [newest data]
//	//Format                      string // <plot> [plot]
//	Output                      string // <power|powerdhnm|powerdlnm|powerdnm|powerdmedian> [power]
//	Plot_height                 string // <value> [500]
//	Plot_width                  string // <value> [1000]
//	Plot_horzaxis               string // <time|freq|frequency> [time]
//	Plot_title                  string // <title|hide> [SNLCQ Target]
//	Plot_titlefont_size         string // <value> [auto]
//	Plot_subtitle               string // <title|hide> [hide]
//	Plot_subtitlefont_size      string // <value> [auto]
//	Plot_axisfont_size          string // <value>
//	Plot_labelfont_size         string // <value>
//	Plot_frequency_label        string // <text|hide> [“Frequency (Hz)”]
//	Plot_frequency_invert       string // <true|false> [false]
//	Plot_frequency_range        string // <min,max> [ful freq of data]
//	Plot_time_format            string // <time-format> [yyyy-MM-dd]
//	Plot_time_label             string // <text|hide> [DATE + RANGE]
//	Plot_time_matchrequest      string // <true|false> [true]
//	Plot_time_tickunit          string // <auto|day|month|year [auto]
//	Plot_time_invert            string // <true|false>
//	Plot_powerscale_autorange   string // <float value 0 to 1> [0.95] | Plot_powerscale_range
//	Plot_powerscale_range		string // <min,max>
//	Plot_powerscale_show        string // <true|false> [true]
//	Plot_powerscale_x           string // <value> [5]
//	Plot_powerscale_y           string // <value> [5]
//	Plot_powerscale_height      string // <value> [12]
//	Plot_powerscale_width       string // <value> [150]
//	Plot_powerscale_orientation string // <horizontal|horz|vertical|vert> [horz]
//	Noisemode_byperiod          string // <pipe-separated-list> []
//	Noisemode_byfrequency       string // <pipe-separated-list> []
//	Nodata                      string // <404|204> [404]
//}
//
//func newOptions() noiseSpectrogramOptions {
//
//	return noiseSpectrogramOptions{
//		//Output: "power",
//		//Format: "plot",
//
//		//Plot_horzaxis:               "time",
//		//Plot_time_matchrequest:      "true",
//		//Plot_time_tickunit:          "auto",
//		//Plot_time_invert:            "false",
//		//Plot_powerscale_show:        "true",
//		//Plot_powerscale_orientation: "horz",
//	}
//}
//

func processCmd(svcname, apiversion string, args []string) {

	//BEGIN these params to come from args or caller
	networks := "II"
	stations := "BFO" //"AAK,ABPO,ALE" //,ARTI,ASCN,BORG"
	//"BRVK", "CMLA", "COCO", "DGAR", "EFI", "ERM", "ESK", "FFC", "HOPE", "JTS", "KAPI", "KDAK", "KIV", "KURK",
	//"KWJN", "LVZ", "MBAR", "MSEY", "MSVF", "NNA", "OBN", "PALK", "PFO", "RAYN", "RPN", "SACV", "SHEL",
	//"SIMI", "SUR", "TAU", "TLY", "UOSS", "WRAB"}
	locations := "00,10"
	channels := "LHZ,LH1,LH2" //,BHZ,BH1,BH2"
	quals := "*"

	now := time.Now().UTC()
	endt := now.Truncate(time.Duration(time.Hour * 24))
	endY, endM, endD := endt.Date()
	startt := endt.Add(-time.Duration(365 * 24 * time.Hour))
	startY, startM, startD := startt.Date()

	beginDate := fmt.Sprintf("%4d-%02d-%02d", startY, startM, startD)
	endDate := fmt.Sprintf("%4d-%02d-%02d", endY, endM, endD)

	format := "plot"

	opts := map[string]string{}
	opts["format"] = format
	opts["plot.powerscale.orientation"] = "vert"

	outdir := fmt.Sprintf("%s/%s", "results", svcname)


	requests := mustang.RequestList(svcname, apiversion, networks, stations, locations, channels, quals, beginDate, endDate, opts)

	results := make(chan *mustang.Result, len(requests))

	go mustang.MakeRequests(requests, opts, results)

	for result := range results {

		filename := fmt.Sprintf("%s/%s/%s",
			outdir,
			strings.ToLower(result.Req.Sta),
			result.FileName())
		err := result.SaveToFile(filename)
		if err != nil {
			fmt.Printf("Error saving image to %s: %s\n", filename, err)
		}
	}

}
