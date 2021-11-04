package postgres

import (
	"context"

	"github.com/neghoda/api/src/models"
)

type FundRepo struct {
	*Connector
}

func (c *Connector) NewFundRepo() *FundRepo {
	return &FundRepo{c}
}

func (fr *FundRepo) FetchFund(ctx context.Context, ticker string) (models.Fund, error) {
	var fund models.Fund

	err := fr.QueryContext(ctx).Model(&fund).
		Where("ticker = ?", ticker).
		Select()
	if err != nil {
		return fund, toServiceError(err)
	}

	err = fr.QueryContext(ctx).Model(&fund.Holdings).
		Where("fund_ticker = ?", ticker).
		OrderExpr("length(weight) DESC").
		Order("weight DESC").
		Select()
	if err != nil {
		return fund, toServiceError(err)
	}

	err = fr.QueryContext(ctx).Model(&fund.Sectors).
		Where("fund_ticker = ?", ticker).
		OrderExpr("length(weight) DESC").
		Order("weight DESC").
		Select()
	if err != nil {
		return fund, toServiceError(err)
	}

	err = fr.QueryContext(ctx).Model(&fund.Countries).
		Where("fund_ticker = ?", ticker).
		OrderExpr("length(weight) DESC").
		Order("weight DESC").
		Select()
	if err != nil {
		return fund, toServiceError(err)
	}

	return fund, toServiceError(err)
}

func (fr *FundRepo) ReplaceFund(ctx context.Context, fund *models.Fund) error {
	err := fr.WithTXContext(ctx, func(query DBQuery) error {
		err := deleteFund(query, fund.Ticker)
		if err != nil {
			return err
		}

		return insertFund(query, fund)
	})

	return toServiceError(err)
}

func deleteFund(query DBQuery, ticker string) error {
	var fund models.Fund

	_, err := query.Model(&fund).
		Where("ticker = ?", ticker).
		Delete()
	if err != nil {
		return toServiceError(err)
	}

	_, err = query.Model(&fund.Holdings).
		Where("fund_ticker = ?", ticker).
		Delete()
	if err != nil {
		return toServiceError(err)
	}

	_, err = query.Model(&fund.Sectors).
		Where("fund_ticker = ?", ticker).
		Delete()
	if err != nil {
		return toServiceError(err)
	}

	_, err = query.Model(&fund.Countries).
		Where("fund_ticker = ?", ticker).
		Delete()
	if err != nil {
		return toServiceError(err)
	}

	return toServiceError(err)
}

func insertFund(query DBQuery, fund *models.Fund) error {
	_, err := query.Model(fund).Insert()
	if err != nil {
		return toServiceError(err)
	}

	if len(fund.Holdings) != 0 {
		_, err = query.Model(&fund.Holdings).Insert()
		if err != nil {
			return toServiceError(err)
		}
	}

	if len(fund.Sectors) != 0 {
		_, err = query.Model(&fund.Sectors).Insert()
		if err != nil {
			return toServiceError(err)
		}
	}

	if len(fund.Countries) != 0 {
		_, err = query.Model(&fund.Countries).Insert()
		if err != nil {
			return toServiceError(err)
		}
	}

	return toServiceError(err)
}
