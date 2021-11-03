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

	fund, err := s.db.QueryContext(ctx).FetchFund(ticker)
	if errors.Is(err, models.ErrNotFound) {
		// Fetch and store fund if there no data in DB
		fund, err = s.updateFund(ctx, ticker)
		if err != nil {
			return fund, err
		}
	}
	if err != nil {
		return models.Fund{}, err
	}

	return fund, err
}

func (s Service) updateFund(ctx context.Context, ticker string) (models.Fund, error) {
	fundURL, ok := mapFundTickerURL[ticker]
	if !ok {
		return models.Fund{}, models.ErrNotFound
	}

	fund, err := s.ssgaRepo.FetchFundPage(fundURL)
	if err != nil {
		return models.Fund{}, err
	}

	tx, err := s.db.NewTXContext(ctx)
	if err != nil {
		return models.Fund{}, err
	}
	defer tx.Rollback()

	err = tx.DeleteFund(ticker)
	if err != nil {
		return models.Fund{}, err
	}

	err = tx.InsertFund(fund)
	if err != nil {
		return models.Fund{}, err
	}

	if err = tx.Commit(); err != nil {
		return models.Fund{}, err
	}

	return fund, err
}

func (s Service) fetchTickerListIfRequired() error {
	var err error

	// Fetch list on request only once in hour
	if len(mapFundTickerURL) == 0 || time.Now().Sub(lastFetched) > time.Hour {
		mapFundTickerURL, err = s.ssgaRepo.FetchFundFinderPage()
		if err != nil {
			return err
		}

		lastFetched = time.Now()
	}

	return nil
}
