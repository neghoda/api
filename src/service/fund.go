package service

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/neghoda/api/src/models"
)

var (
	mapFundTickerURL map[string]string
	lastFetched      time.Time
)

func (s Service) TickerList() ([]string, error) {
	err := s.fetchTickerListIfRequired()
	if err != nil {
		return nil, err
	}

	fundTickerList := make([]string, 0, len(mapFundTickerURL))
	for k := range mapFundTickerURL {
		fundTickerList = append(fundTickerList, k)
	}

	sort.Strings(fundTickerList)

	return fundTickerList, err
}

func (s Service) FundByTicker(ctx context.Context, ticker string) (models.Fund, error) {
	err := s.fetchTickerListIfRequired()
	if err != nil {
		return models.Fund{}, err
	}

	fund, err := s.fundRepo.FetchFund(ctx, ticker)
	if errors.Is(err, models.ErrNotFound) {
		// Fetch and store fund if there no data in DB
		fund, err = s.replaceFund(ctx, ticker)
		if err != nil {
			return fund, err
		}
	}
	if err != nil {
		return models.Fund{}, err
	}

	return fund, err
}

func (s Service) replaceFund(ctx context.Context, ticker string) (models.Fund, error) {
	fundURL, ok := mapFundTickerURL[ticker]
	if !ok {
		return models.Fund{}, models.ErrNotFound
	}

	fund, err := s.ssgaRepo.FetchFundPage(fundURL)
	if err != nil {
		return models.Fund{}, err
	}

	return fund, s.fundRepo.ReplaceFund(ctx, fund)
}

func (s Service) fetchTickerListIfRequired() error {
	var err error

	// Fetch list on request only once in hour
	if len(mapFundTickerURL) == 0 || time.Since(lastFetched) > time.Hour {
		mapFundTickerURL, err = s.ssgaRepo.FetchFundFinderPage()
		if err != nil {
			return err
		}

		lastFetched = time.Now()
	}

	return nil
}
