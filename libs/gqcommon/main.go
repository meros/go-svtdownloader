package gqcommon

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

// GetDoc will download and parse a html document
func GetDoc(url url.URL) (doc *goquery.Document, err error) {
	res, err := http.Get(url.String())
	if err != nil {
		return
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		err = errors.New("Bad status code, expected 200")
		return
	}

	// Load the HTML document
	return goquery.NewDocumentFromReader(res.Body)
}
