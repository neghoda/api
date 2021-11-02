package models

import "github.com/google/uuid"

type Fund struct {
	ID          uuid.UUID `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Ticker      string    `json:"ticker"`
	Holdings    []Holding `json:"holdings" sql:"-"`
	Sectors     []Sector  `json:"sectors" sql:"-"`
	Countries   []Country `json:"countries" sql:"-"`
}

type Holding struct {
	FundTicker string `json:"-"`
	Name       string `json:"name"`
	Share      string `json:"share"`
	Weight     string `json:"weight"`
}

type Sector struct {
	FundTicker string `json:"-"`
	Name       string `json:"name"`
	Weight     string `json:"weight"`
}

type Country struct {
	FundTicker string `json:"-"`
	Name       string `json:"name"`
	Weight     string `json:"weight"`
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
