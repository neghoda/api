package ssga

import (
	"errors"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"golang.org/x/net/html"
)

const (
	shema         = "https"
	host          = "www.ssga.com"
	usContrySlug  = "us"
	enCountrySlug = "en"

	fundFinderPagePath  = "bin/v1/ssmp/fund/fundfinder"
	fundFinderPageQuery = "country=%v&language=%v&role=individual&product=etfs&ui=fund-finder"

	fetchingPageMsg = "Fetching page %v\r"
)

var (
	errBadResponse = errors.New("bad response")
)

// Client is an SSGA client used to parse search results from ssga.com.
type Client struct {
	shema        string
	host         string
	contrySlug   string
	languageSlug string

	httpClient http.Client
}

func NewClient(
	httpClient http.Client,
) *Client {
	return &Client{
		shema:        shema,
		host:         host,
		contrySlug:   usContrySlug,
		languageSlug: enCountrySlug,
		httpClient:   httpClient,
	}
}

func extractAttributeValue(n *html.Node, attr string) string {
	for _, v := range n.Attr {
		if v.Key == attr {
			return v.Val
		}
	}

	return ""
}

func logRespStatus(resp *http.Response) {
	log.Printf("Response for %v is %v not 200", resp.Request.URL, resp.StatusCode)

	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Response body %v", string(resBody))
}
