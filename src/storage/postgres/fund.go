package postgres

import "github.com/erp/api/src/models"

func (q DBQuery) InsertFund(fund models.Fund) error {
	_, err := q.Model(&fund).Insert()
	if err != nil {
		return toServiceError(err)
	}

	if len(fund.Holdings) != 0 {
		_, err = q.Model(&fund.Holdings).Insert()
		if err != nil {
			return toServiceError(err)
		}
	}

	if len(fund.Sectors) != 0 {
		_, err = q.Model(&fund.Sectors).Insert()
		if err != nil {
			return toServiceError(err)
		}
	}

	if len(fund.Countries) != 0 {
		_, err = q.Model(&fund.Countries).Insert()
		if err != nil {
			return toServiceError(err)
		}
	}

	return toServiceError(err)
}

func (q DBQuery) FetchFund(ticker string) (models.Fund, error) {
	var fund models.Fund

	err := q.Model(&fund).
		Where("ticker = ?", ticker).
		Select()
	if err != nil {
		return fund, toServiceError(err)
	}

	err = q.Model(&fund.Holdings).
		Where("fund_ticker = ?", ticker).
		OrderExpr("length(weight) DESC").
		Order("weight DESC").
		Select()
	if err != nil {
		return fund, toServiceError(err)
	}

	err = q.Model(&fund.Sectors).
		Where("fund_ticker = ?", ticker).
		OrderExpr("length(weight) DESC").
		Order("weight DESC").
		Select()
	if err != nil {
		return fund, toServiceError(err)
	}

	err = q.Model(&fund.Countries).
		Where("fund_ticker = ?", ticker).
		OrderExpr("length(weight) DESC").
		Order("weight DESC").
		Select()
	if err != nil {
		return fund, toServiceError(err)
	}

	return fund, toServiceError(err)
}

func (q DBQuery) DeleteFund(ticker string) error {
	var fund models.Fund

	_, err := q.Model(&fund).
		Where("ticker = ?", ticker).
		Delete()
	if err != nil {
		return toServiceError(err)
	}

	_, err = q.Model(&fund.Holdings).
		Where("fund_ticker = ?", ticker).
		Delete()
	if err != nil {
		return toServiceError(err)
	}

	_, err = q.Model(&fund.Sectors).
		Where("fund_ticker = ?", ticker).
		Delete()
	if err != nil {
		return toServiceError(err)
	}

	_, err = q.Model(&fund.Countries).
		Where("fund_ticker = ?", ticker).
		Delete()
	if err != nil {
		return toServiceError(err)
	}

	return toServiceError(err)
}
