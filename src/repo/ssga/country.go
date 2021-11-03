package ssga

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/neghoda/api/src/models"
	"golang.org/x/net/html"
)

var (
	badCountriesListMsg = "Counties list cannot be parsed for fund with ID %v\n"
)

func lookForCountries(n *html.Node, fund *models.Fund) {
	if n.Type == html.ElementNode && n.Data == "input" {
		id := extractAttributeValue(n, "id")

		if strings.EqualFold(id, "fund-geographical-breakdown") {
			if err := parseForCountriesJSON(extractAttributeValue(n, "value"), fund); err != nil {
				log.Printf(badCountriesListMsg, fund.ID)
			}
		}
	}
}

// Countries list seams to be client-side rendered and comes as JSON inserted into HTML
func parseForCountriesJSON(countryData string, fund *models.Fund) error {
	var countriesInfo models.CountriesInfo

	dec := json.NewDecoder(strings.NewReader(countryData))
	err := dec.Decode(&countriesInfo)
	if err != nil {
		return err
	}

	countries := make([]models.Country, 0, len(countriesInfo.AttrArray))
	for i := range countriesInfo.AttrArray {
		countries = append(countries, models.Country{
			FundTicker: fund.Ticker,
			Name:       countriesInfo.AttrArray[i].Name.Value,
			Weight:     countriesInfo.AttrArray[i].Weight.Value,
		})
	}

	fund.Countries = countries

	return nil
}
