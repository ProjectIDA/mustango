package mustang

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type EPResult struct {
	Req    *EPRequest
	Resbuf []byte
}

func (resptr *EPResult) SaveToFile(fn string) (error) {
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

