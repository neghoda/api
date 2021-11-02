package ssga

import (
	"strings"

	"github.com/erp/api/src/models"
	"golang.org/x/net/html"
)

func lookForHoldings(n *html.Node, fund *models.Fund) {
	if n.Type == html.ElementNode && n.Data == "tbody" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			lookForHolding(c, fund)
		}
	}
}

func lookForHolding(n *html.Node, fund *models.Fund) {
	holding := models.Holding{
		FundTicker: fund.Ticker,
	}

	if n.Type == html.ElementNode && n.Data == "tr" {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			lookForHoldingName(c, &holding)
			lookForHoldingShares(c, &holding)
			lookForHoldingWeight(c, &holding)
		}
	}

	if holding.Name != "" && holding.Share != "" && holding.Weight != "" {
		fund.Holdings = append(fund.Holdings, holding)
	}
}

func lookForHoldingName(n *html.Node, holding *models.Holding) {
	if n.Type == html.ElementNode && n.Data == "td" {
		label := extractAttributeValue(n, "data-label")

		if strings.EqualFold(label, "Name:") {
			if n.FirstChild != nil {
				holding.Name = strings.TrimSpace(n.FirstChild.Data)
			}
		}
	}
}

func lookForHoldingShares(n *html.Node, holding *models.Holding) {
	if n.Type == html.ElementNode && n.Data == "td" {
		label := extractAttributeValue(n, "data-label")

		if strings.EqualFold(label, "Shares Held:") ||
			strings.EqualFold(label, "Par Value:") ||
			strings.EqualFold(label, "Market Value:") {
			if n.FirstChild != nil {
				holding.Share = strings.TrimSpace(n.FirstChild.Data)
			}
		}
	}
}

func lookForHoldingWeight(n *html.Node, holding *models.Holding) {
	if n.Type == html.ElementNode && n.Data == "td" {
		label := extractAttributeValue(n, "data-label")

		if strings.EqualFold(label, "Weight:") {
			if n.FirstChild != nil {
				holding.Weight = strings.TrimSpace(n.FirstChild.Data)
			}
		}
	}
}
