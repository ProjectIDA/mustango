package mustang

import (
	"fmt"
	"net/http"
)

type Urler interface {
	Url() (string, error)
}

func Get(req Urler) (*http.Response, error) {

	url, err := req.Url()
	fmt.Printf("URL: %s\n", url)
	if err != nil {
		return &http.Response{}, err
	}

	var res *http.Response
	res, err = http.Get(url)

	return res, err
}
