package ssga

import (
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/neghoda/api/src/models"
	"golang.org/x/net/html"
)

func (c *Client) FetchFundPage(fundPath string) (models.Fund, error) {
	fundURL := c.fundURL(fundPath)

	log.Printf(fetchingPageMsg, fundURL)

	resp, err := c.httpClient.Get(fundURL)
	if err != nil {
		return models.Fund{}, err
	}

	if resp.StatusCode != http.StatusOK {
		logRespStatus(resp)

		return models.Fund{}, errBadResponse
	}
	defer resp.Body.Close()

	fund, err := parseFundPage(resp.Body)
	if err != nil {
		return models.Fund{}, err
	}

	return fund, err
}

func parseFundPage(r io.Reader) (models.Fund, error) {
	document, err := html.Parse(r)
	if err != nil {
		return models.Fund{}, err
	}

	fund := models.Fund{
		ID: uuid.New(),
	}

	extractFundInfo(document, &fund)

	return fund, nil
}

func extractFundInfo(n *html.Node, fund *models.Fund) {
	if n.Type == html.ElementNode && n.Data == "h1" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			lookForFundName(c, fund)
			lookForTicker(c, fund)
		}
	}

	lookForFundDescription(n, fund)

	if fund.Ticker != "" {
		lookForHoldings(n, fund)
		lookForSectors(n, fund)
		lookForCountries(n, fund)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractFundInfo(c, fund)
	}
}

func lookForTicker(n *html.Node, fund *models.Fund) {
	if n.Type == html.ElementNode && n.Data == "span" {
		class := extractAttributeValue(n, "class")

		if strings.EqualFold(class, "ticker") {
			if n.FirstChild != nil {
				fund.Ticker = strings.TrimSpace(n.FirstChild.Data)
			}
		}
	}
}

func lookForFundName(n *html.Node, fund *models.Fund) {
	if n.Type == html.ElementNode && n.Data == "span" {
		class := extractAttributeValue(n, "class")

		if strings.EqualFold(class, "") {
			if n.FirstChild != nil {
				fund.Name = strings.TrimSpace(n.FirstChild.Data)
			}
		}
	}
}

func lookForFundDescription(n *html.Node, fund *models.Fund) {
	if n.Type == html.ElementNode && n.Data == "meta" {
		class := extractAttributeValue(n, "name")

		if strings.EqualFold(class, "description") {
			content := extractAttributeValue(n, "content")
			fund.Description = strings.TrimSpace(content)
		}
	}
}

func (c *Client) fundURL(fundPath string) string {
	u := url.URL{
		Scheme: c.shema,
		Host:   c.host,
		Path:   fundPath,
	}

	return u.String()
}
