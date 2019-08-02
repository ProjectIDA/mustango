package mustang

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type Result struct {
	Req      *Request
	HTTPResp *http.Response
	Resbuf   []byte
}

func (resptr *Result) SaveToFile(fn string) (error) {
	//TODO change printf's top log statements or remove

	err := os.MkdirAll(filepath.Dir(fn), 0775)
	if err != nil {
		fmt.Printf("Error creating directory %s: %s\n", filepath.Dir(fn), err)
		return err
	}

	err = ioutil.WriteFile(fn, resptr.Resbuf, 0644)
	if err != nil {
		fmt.Printf("Error saving image to %s: %s\n", fn, err)
	}
	return err
}

func (r *Result) FileName(name_tail string) (filename string) {

	if name_tail == "" {
		name_tail = r.Req.Starttime + "." + r.Req.Endtime
	}
	fn := fmt.Sprintf("%s.%s.%s.%s.%s.%s",
		r.Req.Net,
		r.Req.Sta,
		r.Req.Loc,
		r.Req.Chn,
		name_tail,
		r.FormatFileExt())

	return fn
}

func (r *Result) FormatFileExt() string {

	ctype := r.HTTPResp.Header.Get("Content-Type")
	//fmt.Printf("%s: %s\n", ctype, r.HTTPResp.Request.URL)

	switch ctype {
	case "image/png":
		return "png"
	case "csv":
		return "csv"
	case "xml":
		return "xml"
	case "text/plain":
		return "txt"
	case "json":
		return "json"
	default:
		return "UNK"
	}
}

