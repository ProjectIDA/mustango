package mustang

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Urler interface {
	Url() (string, error)
}

func Get(req Urler) ([]byte, error) {

	url, err := req.Url()
	if err != nil {
		return []byte{}, err
	}

	var res *http.Response

	res, err = http.Get(url)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
	case 404:
		return []byte{}, fmt.Errorf("http status code %d. No data found for request", res.StatusCode)
	default:
		return []byte{}, fmt.Errorf("http status code %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil

}
