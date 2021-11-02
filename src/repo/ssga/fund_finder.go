package ssga

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/erp/api/src/models"
)

// Fund finder page half client-side rendered. List of ETFs fetched separately as JSON
func (c *Client) FetchFundFinderPage() (map[string]string, error) {
	fundFinderPageURL := c.fundFinderPageURL()

	log.Printf(fetchingPageMsg, fundFinderPageURL)

	resp, err := c.httpClient.Get(fundFinderPageURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		logRespStatus(resp)

		return nil, errBadResponse
	}
	defer resp.Body.Close()

	return parseFundFinderResponse(resp.Body)
}

func parseFundFinderResponse(r io.Reader) (map[string]string, error) {
	var (
		fundFinderInfo models.FundFinderInfo

		mapTickerURL = make(map[string]string)
		dec          = json.NewDecoder(r)
	)
	err := dec.Decode(&fundFinderInfo)
	if err != nil {
		return nil, err
	}

	for _, v := range fundFinderInfo.Data.Funds.Etfs.Datas {
		mapTickerURL[v.FundTicker] = v.FundURI
	}

	return mapTickerURL, nil
}

func (c *Client) fundFinderPageURL() string {
	u := url.URL{
		Scheme: c.shema,
		Host:   c.host,
		Path:   fundFinderPagePath,
	}

	u.RawQuery = fmt.Sprintf(fundFinderPageQuery, c.contrySlug, c.languageSlug)

	return u.String()
}
