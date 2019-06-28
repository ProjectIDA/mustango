package noise_mode_timeseries

//start=<yyyy-mm-dd>
//end=<yyyy-mm-dd>
//
//output = <power|powerdhnm|powerdlnm|powerdnm>
//
//noisemodel.byperiod = <pipe-separated-list>
//noisemodel.byfrequency = <pipe-separated-list>
//
//plot.height = <pixel-height>
//plot.width = <pixel-width>
//plot.title = <title|hide>
//plot.subtitle = <title|hide>
//plot.power.min = <value>
//plot.power.max = <value>
//plot.power.label = <'title'|hide>
//plot.legend = <show|hide>
//plot.titlefont.size = <value>
//plot.subtitlefont.size = <value>
//plot.powerlabelfont.size = <value>
//plot.poweraxisfont.size = <value>
//plot.timeaxisfont.size = <value>
//plot.legendfont.size = <value>
//plot.backgroundcolor = <color|r,g,b>
//plot.gridcolor = <color|r,g,b>
//plot.linewidth = <value>
//
//frequencies = <freq1,freq1...>|<all>|<range-start,range-end> |
//periods = <period1,period2,...>|<all>|<range-start,range-end>
//
//text.style = <list|table>
//text.units = <hertz|seconds>
//xml.style = <byday|byperiod|byfrequency|flat>
//xml.units = <hertz|seconds>
//format = <plot|text|xml>
//nodata = <404|204>
//

import (
	"fmt"
	"mustango/pkg/mustang"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var NoiseModeTimeSeriesCmd = &cobra.Command{
	Use:   "noise-mode-timeseries",
	Short: "Retrieve Noise mode Time Series from IRIS MUSTANG 'noise-mode-timeseries' Web Service",
	Run: func(cmd *cobra.Command, args []string) {
		processQuery("noise-mode-timeseries", args)
	},
}

func processQuery(svcname string, args []string) {

	//BEGIN these params to come from args or caller
	now := time.Now().UTC()

	endt := now.Truncate(time.Duration(time.Hour * 24))
	endY, endM, endD := endt.Date()

	startt := endt.Add( -time.Duration(365 * 24 * time.Hour))
	startY, startM, startD := startt.Date()

	beginDate := fmt.Sprintf("%4d-%02d-%02d", startY, startM, startD)
	endDate := fmt.Sprintf("%4d-%02d-%02d", endY, endM, endD)

	network := "II"
	stations := []string{"AAK", }  // "ABPO", "ALE", "ARTI" ,
	//"ASCN", "BFO", "BORG", }
	//"BRVK", "CMLA", "COCO", "DGAR", "EFI", "ERM", "ESK",
	//"FFC", "HOPE", "JTS", "KAPI", "KDAK", "KIV", "KURK",
	//"KWJN", "LVZ", "MBAR", "MSEY", "MSVF", "NNA",
	//"OBN", "PALK", "PFO", "RAYN", "RPN", "SACV", "SHEL",
	//"SIMI", "SUR", "TAU", "TLY", "UOSS", "WRAB"}
	locs := []string{"00", "10"}
	chns := []string{"LHZ", "LH1", "LH2", "BHZ", "BH1", "BH2"}
	//qua := "M"

	format := "plot"
	//output := "power"
	//
	//beginDate := fmt.Sprintf("%4d-%02d-%02d", startY, startM, startD)
	//endDate := fmt.Sprintf("%4d-%02d-%02d", endY, endM, endD)
	//
	//options := newOptions()
	//options.BeginDate, options.EndDate = beginDate, endDate
	//options.Format, options.Output = "plot", "power"


	//targets := []string{}
	//for _, sta := range stations {
	//	for _, loc := range locs {
	//		for _, chn := range chns {
	//			targets = append(targets, fmt.Sprintf("%s.%s.%s.%s.%s",
	//				network,
	//				sta,
	//				loc,
	//				chn,
	//				qua))
	//		}
	//	}
	//}

	//END

	var reschn = make(chan *mustang.Result, len(stations) * len(locs) * len(chns))
	go func (reschn chan *mustang.Result) {

		var wg sync.WaitGroup

		for _, sta := range stations {
			for _, loc := range locs {
				for _, chn := range chns {

					wg.Add(1)
					time.Sleep(time.Millisecond * 330)
					fmt.Printf("Querying for: %s:%s:%s\n", sta, loc, chn)
					areq := &mustang.Request{
						Service:   svcname,
						Version:   "1",
						Net:       network,
						Sta:       sta,
						Loc:       loc,
						Chn:       chn,
						//Starttime: beginDate,
						//Endtime:   endDate,
						Format:    format,
					}

					go func(wg *sync.WaitGroup, areq *mustang.Request) {
						defer wg.Done()
						resbuf, err := mustang.Get(areq)
						if err != nil {
							fmt.Printf("Error executing MUSTANG request: %s, %s\n", areq, err)
						}

						// contruct Result with result buffer and original request
						epres := &mustang.Result{Req: areq, Resbuf: resbuf}
						reschn <- epres

					}(&wg, areq)
				}
			}

		}

		go func (wg *sync.WaitGroup, results chan *mustang.Result) {
			wg.Wait()
			close(results)
		}(&wg, reschn)

	}(reschn)

	func(results chan *mustang.Result) {

		for resptr := range results {

			if len(resptr.Resbuf) > 0 {
				filename := fmt.Sprintf("results/%s/%s/%s",
					resptr.Req.Service,
					strings.ToLower(resptr.Req.Sta),
					resptr.Req.FileName())
				err := resptr.SaveToFile(filename)
				if err != nil {
					fmt.Printf("Error saving image to %s: %s\n", filename, err)
				}
			}
		}

	}(reschn)
}
