package noise_spectrograms

import (
	"fmt"
	"mustango/pkg/mustang"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var NoiseSpectrogramsCmd = &cobra.Command{
	Use:   "noise-spectrogram",
	Short: "Retrieve Noise-Spectrogram plots (png) from IRIS MUSTANG 'noise-spectrograms' Web Service",
	Run: func(cmd *cobra.Command, args []string) {
		ProcessQuery("noise-spectrogram", args)
	},
}

func ProcessQuery(svcname string, args []string) {

		now := time.Now().UTC()

		endt := now.Truncate(time.Duration(time.Hour * 24))
		endY, endM, endD := endt.Date()

		startt := endt.Add( -time.Duration(365 * 24 * time.Hour))
		startY, startM, startD := startt.Date()

		beginDate := fmt.Sprintf("%4d-%02d-%02d", startY, startM, startD)
		endDate := fmt.Sprintf("%4d-%02d-%02d", endY, endM, endD)

		ii_network := []string{"AAK", }  // "ABPO", "ALE", "ARTI" ,
			//"ASCN", "BFO", "BORG", }
		//"BRVK", "CMLA", "COCO", "DGAR", "EFI", "ERM", "ESK",
		//"FFC", "HOPE", "JTS", "KAPI", "KDAK", "KIV", "KURK",
		//"KWJN", "LVZ", "MBAR", "MSEY", "MSVF", "NNA",
		//"OBN", "PALK", "PFO", "RAYN", "RPN", "SACV", "SHEL",
		//"SIMI", "SUR", "TAU", "TLY", "UOSS", "WRAB"}
		locs := []string{"00", "10"}
		chns := []string{"LHZ", "LH1", "LH2", "BHZ", "BH1", "BH2"}

		var reschn = make(chan *mustang.EPResult, len(ii_network) * len(locs) * len(chns))
		go func (reschn chan *mustang.EPResult) {

			var wg sync.WaitGroup

			for _, sta := range ii_network {
				for _, loc := range locs {
					for _, chn := range chns {

						wg.Add(1)
						time.Sleep(time.Millisecond * 330)
						fmt.Printf("Querying for: %s:%s:%s\n", sta, loc, chn)
						areq := &mustang.EPRequest{
							Service:   svcname,
							Version:   "1",
							Netw:      "II",
							Sta:       sta,
							Loc:       loc,
							Chn:       chn,
							Starttime: beginDate,
							Endtime:   endDate,
							Format:    "plot",
							Output:    "power", //fn,
						}

						go func(wg *sync.WaitGroup, areq *mustang.EPRequest) {
							defer wg.Done()
							resbuf, err := mustang.Get(areq)
							if err != nil {
								fmt.Printf("Error executing MUSTANG request: %s, %s\n", areq, err)
							}

							// contruct EPResult with result buffer and original request
							epres := &mustang.EPResult{Req: areq, Resbuf: resbuf}
							reschn <- epres

						}(&wg, areq)
					}
				}

			}

			go func (wg *sync.WaitGroup, results chan *mustang.EPResult) {
				wg.Wait()
				close(results)
			}(&wg, reschn)

		}(reschn)

		func(results chan *mustang.EPResult) {

			for resptr := range results {

				if len(resptr.Resbuf) > 0 {
					pngFilename := fmt.Sprintf("results/%s/%s/%s",
						resptr.Req.Service,
						strings.ToLower(resptr.Req.Sta),
						resptr.Req.FileName())
					err := resptr.SaveToFile(pngFilename)
					if err != nil {
						fmt.Printf("Error saving image to %s: %s\n", pngFilename, err)
					}
				}
			}

		}(reschn)
}


