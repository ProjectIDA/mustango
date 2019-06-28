package mustang

import (
	"strings"
)

type Target struct {
	Net string
	Sta string
	Loc string
	Chn string
	Qua string
}

func TargetListFromCSV(nets, stas, locs, chns, quals string) ([]*Target) {

	targets := []*Target{}

	netlist := strings.Split(nets, ",")
	if len(netlist) == 0 {
		return targets
	}
	stalist := strings.Split(stas, ",")
	if len(stalist) == 0 {
		return targets
	}
	loclist := strings.Split(locs, ",")
	if len(loclist) == 0 {
		return targets
	}
	chnlist := strings.Split(chns, ",")
	if len(chnlist) == 0 {
		return targets
	}
	quallist := strings.Split(quals, ",")
	if len(chnlist) == 0 {
		return targets
	}

	for _, net := range netlist {
		for _, sta := range stalist {
			for _, loc := range loclist {
				for _, chn := range chnlist {
					for _, qua := range quallist {
						targets = append(targets, &Target{
							strings.ToUpper(net),
							strings.ToUpper(sta),
							strings.ToUpper(loc),
							strings.ToUpper(chn),
							strings.ToUpper(qua),
						})
					}
				}
			}
		}
	}
	return targets
}


//func (t Target) NSLCQ() string {
//
//	return fmt.Sprintf("%s.%s.%s.%s", t.Net, t.Sta, t.Loc, t.Chn, t.Qua)
//
//}
//
//func(t Target) UrlValues() url.Values {
//
//	params := url.Values{}
//	params.Add("net", t.Net)
//	params.Add("sta", t.Sta)
//	params.Add("loc", t.Loc)
//	params.Add("chn", t.Chn)
//
//	return params
//}

