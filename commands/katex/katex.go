package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"unsafe"
)

const (
	TEMP_FILE = "/tmp/%x.png"
	URL       = "https://render-katex.herokuapp.com/tex/"
)

func katex(math string) (*os.File, error) {
	sum := sha1.Sum([]byte(math))           // sha1sum of math expression
	tempFile := fmt.Sprintf(TEMP_FILE, sum) // tempFile of png
	u, err := url.Parse(URL)
	if err != nil {
		return nil, err
	}
	u, err = u.Parse(math)
	url := u.String() // url of request
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	switch res.StatusCode {
	case 400:
		errmsg, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(*(*string)(unsafe.Pointer(&errmsg)))
	case 200:
		break
	default:
		return nil, errors.New(res.Status)
	}
	if _, err := os.Stat(tempFile); err != nil { // generate png if necessary
		exec := exec.Command("xvfb-run", "--", "wkhtmltoimage", "--width", "512", url, tempFile)
		err := exec.Run()
		if err != nil {
			defer os.Remove(tempFile)
			return nil, err
		}
	}
	return os.OpenFile(tempFile, os.O_RDWR, 0664)
}
