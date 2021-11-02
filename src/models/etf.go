package models

import "github.com/google/uuid"

type Fund struct {
	ID          uuid.UUID
	Name        string
	Description string
	Ticker      string
	Holdings    []Holding
	Sectors     []Sector
	Countries   []Country
}

type Holding struct {
	FundTicker string
	Name       string
	Shares     string
	Weight     string
}

type Sector struct {
	FundTicker string
	Name       string
	Weight     string
}

type Country struct {
	FundTicker string
	Name       string
	Weight     string
}

type CountriesInfo struct {
	Label          string `json:"label"`
	AsOfDate       string `json:"asOfDate"`
	AsOfDateSimple string `json:"asOfDateSimple"`
	AttrArray      []struct {
		Name struct {
			Label string `json:"label"`
			Value string `json:"value"`
		} `json:"name"`
		Weight struct {
			Label string `json:"label"`
			Value string `json:"value"`
		} `json:"weight"`
	} `json:"attrArray"`
}

type FundFinderInfo struct {
	Data struct {
		Funds struct {
			Etfs struct {
				Datas []struct {
					Domicile   string `json:"domicile"`
					FundName   string `json:"fundName"`
					FundTicker string `json:"fundTicker"`
					FundURI    string `json:"fundUri"`
				}
			} `json:"etfs"`
		} `json:"funds"`
	} `json:"data"`
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}
