package mustang

import "fmt"

type Target struct {
	Net string
	Sta string
	Loc string
	Chn string
	Qua string
}

func (t Target) nslcq() string {

	return fmt.Sprintf("%s.%s.%s.%s.%s", t.Net, t.Sta, t.Loc, t.Chn, t.Qua)

}
