package mustang

import (
	"fmt"
	"net/url"
)

const (
	WebSvcProtocol    = "https"
	WebSvcHost        = "service.iris.edu"
	WebSvcUrlTemplate = "%s://%s/mustang/%s/%s/query"
)

//type MustangFormatType int
//
//const (
//	MustangFormatTypePNG MustangFormatType = iota
//	MustangFormatTypeJSON
//	MustangFormatTypeCSV
//	MustangFormatTypeTEXT
//	MustangFormatTypeXML
//)
//
//var mustangFormatTypeMap = map[string]MustangFormatType{
//	"plot": MustangFormatTypePNG,
//	"json": MustangFormatTypeJSON,
//	"csv": MustangFormatTypeCSV,
//	"text": MustangFormatTypeTEXT,
//	"xml": MustangFormatTypeXML,
//}

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
//plot_powerscale_height = <value>plot_powerscale_width = <value>
//plot_powerscale_show = <true|false>plot_powerscale_x = <value>plot_powerscale_y = <value>
//plot_powerscale_orientation = <horizontal|horz|vertical|vert>
//
//plot_title = <title|hide>plot_subtitle = <title|hide>
//plot_titlefont_size = <value>plot_subtitlefont_size = <value>
//plot_axisfont_size = <value>plot_labelfont_size = <value>
//
//plot_time_matchrequest = <true|false>plot_time_format = <time-format>plot_time_label = <text|hide>
//plot_time_tickunit = <auto|day|month|year>plot_time_invert = <true|false>
//
//plot_frequency_label = <text|hide>plot_frequency_invert = <true|false>
//plot_frequency_range = <min,max>


type EPRequest struct {
	Service   string
	Version   string
	Netw      string
	Sta       string
	Loc       string
	Chn       string
	Starttime string
	Endtime   string
	Format    string
	Output    string
}

func (r *EPRequest) String() string {
	rurl, _ := r.Url()
	return fmt.Sprintf("%s\n%s/%s %s %s %s %s %s %s",
		r.Service, r.Version,
		r.Netw, r.Sta, r.Loc, r.Chn,
		r.Starttime, r.Endtime,
		rurl)
}

func (r *EPRequest) serviceUrl() string {
	return fmt.Sprintf(WebSvcUrlTemplate,
		WebSvcProtocol,
		WebSvcHost,
		r.Service,
		r.Version)
}

func (r *EPRequest) queryParams() url.Values {

	params := url.Values{}
	params.Add("net", r.Netw)
	params.Add("sta", r.Sta)
	params.Add("loc", r.Loc)
	params.Add("cha", r.Chn)
	params.Add("starttime", r.Starttime)
	params.Add("endtime", r.Endtime)
	params.Add("output", r.Output)
	params.Add("format", r.Format)
	params.Add("plot.horzaxis", "time")
	params.Add("plot.time.matchrequest", "true")
	params.Add("plot.time.tickunit", "auto")
	params.Add("plot.time.invert", "false")
	params.Add("plot.powerscale.show", "true")
	params.Add("plot.powerscale.orientation", "horz")
	params.Add("nodata", "404")

	return params
}

func (r EPRequest) Url() (string, error) {

	serviceURL := r.serviceUrl()

	baseUrl, err := url.Parse(serviceURL)
	if err != nil {
		return "", err
	}

	params := r.queryParams()

	baseUrl.RawQuery = params.Encode()

	return baseUrl.String(), err
}

func (r *EPRequest) FormatFileExt() string {
	switch r.Format {
	case "plot":
		return "png"
	case "csv":
		return "csv"
	case "xml":
		return "xml"
	case "text":
		return "txt"
	case "json":
		return "json"
	default:
		return "UNK"
	}
}

func (req *EPRequest) FileName() (filename string) {

	fn := fmt.Sprintf("%s.%s.%s.%s.%s.%s.%s",
		req.Netw,
		req.Sta,
		req.Loc,
		req.Chn,
		req.Starttime,
		req.Endtime,
		req.FormatFileExt())

	return fn
}
