package ssga

import (
	"strings"

	"github.com/neghoda/api/src/models"
	"golang.org/x/net/html"
)

func lookForSectors(n *html.Node, fund *models.Fund) {
	if n.Type == html.ElementNode && n.Data == "div" {
		class := extractAttributeValue(n, "class")

		if strings.EqualFold(class, "fund-sector-breakdown ") {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				walkThroughSectorsChildren(c, fund)
			}
		}
	}
}

// For every child of "fund-sector-breakdown" check everything below
func walkThroughSectorsChildren(n *html.Node, fund *models.Fund) {
	if n.Type == html.ElementNode && n.Data == "tbody" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			lookForSector(c, fund)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkThroughSectorsChildren(c, fund)
	}
}

func lookForSector(n *html.Node, fund *models.Fund) {
	sector := models.Sector{
		FundTicker: fund.Ticker,
	}

	if n.Type == html.ElementNode && n.Data == "tr" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			lookForSectorName(c, &sector)
			lookForSectorWeight(c, &sector)
		}
	}

	if sector.Name != "" && sector.Weight != "" {
		fund.Sectors = append(fund.Sectors, sector)
	}
}

func lookForSectorName(n *html.Node, sector *models.Sector) {
	if n.Type == html.ElementNode && n.Data == "td" {
		class := extractAttributeValue(n, "class")

		if strings.EqualFold(class, "label") {
			if n.FirstChild != nil {
				sector.Name = strings.TrimSpace(n.FirstChild.Data)
			}
		}
	}
}

func lookForSectorWeight(n *html.Node, sector *models.Sector) {
	if n.Type == html.ElementNode && n.Data == "td" {
		class := extractAttributeValue(n, "class")

		if strings.EqualFold(class, "data") {
			if n.FirstChild != nil {
				sector.Weight = strings.TrimSpace(n.FirstChild.Data)
			}
		}
	}
}
