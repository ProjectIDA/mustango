package mustang

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	WebSvcProtocol    = "https"
	WebSvcHost        = "service.iris.edu"
	WebSvcUrlTemplate = "%s://%s/mustang/%s/%s/query"
	ThrottleDelay     = time.Millisecond * 330
)

type ParamMapper interface {
	ParamMap() map[string]string
}
type Request struct {
	Service   string
	Version   string
	Net       string
	Sta       string
	Loc       string
	Chn       string
	Qua       string
	Starttime string
	Endtime   string
	//Format    string

	Opts *map[string]string
}

func (r *Request) String() string {
	rurl, _ := r.Url()
	return fmt.Sprintf("%s\n%s/%s %s %s %s %s %s %s",
		r.Service, r.Version,
		r.Net, r.Sta, r.Loc, r.Chn,
		r.Starttime, r.Endtime,
		rurl)
}

func (r *Request) serviceUrl() string {
	return fmt.Sprintf(WebSvcUrlTemplate,
		WebSvcProtocol,
		WebSvcHost,
		r.Service,
		r.Version)
}

func (r *Request) NSLCQ() string {

	return fmt.Sprintf("%s.%s.%s.%s.%s", r.Net, r.Sta, r.Loc, r.Chn, r.Qua)

}

func (r *Request) queryParams() url.Values {

	params := url.Values{}
	params.Add("target", r.NSLCQ())
	params.Add("starttime", r.Starttime)
	params.Add("endtime", r.Endtime)
	//params.Add("format", r.Format)

	for k, v := range *r.Opts {
		params.Add(k, v)
	}

	return params
}

func (r Request) Url() (string, error) {

	serviceURL := r.serviceUrl()

	baseUrl, err := url.Parse(serviceURL)
	if err != nil {
		return "", err
	}

	params := r.queryParams()

	baseUrl.RawQuery = params.Encode()

	return baseUrl.String(), err
}

func RequestList(svcname, apiversion, nets, stas, locs, chns, quals, begdt, enddt string, opts map[string]string) []*Request {

	reqlist := []*Request{}

	netlist := strings.Split(nets, ",")
	if len(netlist) == 0 {
		return reqlist
	}
	stalist := strings.Split(stas, ",")
	if len(stalist) == 0 {
		return reqlist
	}
	loclist := strings.Split(locs, ",")
	if len(loclist) == 0 {
		return reqlist
	}
	chnlist := strings.Split(chns, ",")
	if len(chnlist) == 0 {
		return reqlist
	}
	quallist := strings.Split(quals, ",")
	if len(chnlist) == 0 {
		return reqlist
	}

	for _, net := range netlist {
		for _, sta := range stalist {
			for _, loc := range loclist {
				for _, chn := range chnlist {
					for _, qua := range quallist {
						areq := &Request{
							Service:   svcname,
							Version:   apiversion,
							Net:       net,
							Sta:       sta,
							Loc:       loc,
							Chn:       chn,
							Qua:       qua,
							Starttime: begdt,
							Endtime:   enddt,
							//Format:    format,
							Opts: &opts,
						}
						reqlist = append(reqlist, areq)
					}
				}
			}
		}
	}
	return reqlist
}

func getResponseBody(res *http.Response) ([]byte, error) {

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func makeRequest(wg *sync.WaitGroup, areq *Request, reschn chan *Result) {

	defer wg.Done()

	httpResp, err := Get(areq)
	if err != nil {
		fmt.Printf("Error executing MUSTANG HTTP request: %s\n", err)
		return
	}

	resbuf, err := getResponseBody(httpResp)
	if err != nil {
		fmt.Printf("Error retrieving Body from HTTP Response: %s\n", err)
	}

	// contruct Result with result buffer and original request
	epres := &Result{Req: areq, HTTPResp: httpResp, Resbuf: resbuf}
	reschn <- epres

}

func MakeRequests(requests []*Request, opts map[string]string, reschn chan *Result) {

	targetCount := len(requests)

	var wg sync.WaitGroup
	wg.Add(targetCount)

	for _, request := range requests {

		fmt.Printf("Querying for: %v\n", request.NSLCQ())

		go makeRequest(&wg, request, reschn)

	}

	go func(wg *sync.WaitGroup, results chan *Result) {
		wg.Wait()
		close(results)
	}(&wg, reschn)

}
