package cron

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	errMsg      = "syncFundData faild with err: %v"
	finishedMsg = "job finished in %v"
)

func (cw *CronWrapper) syncFundsData() {
	start := time.Now()

	mapFundTickerURL, err := cw.ssgaRepo.FetchFundFinderPage()
	if err != nil {
		log.Errorf(errMsg, err)
	}

	for ticker, url := range mapFundTickerURL {
		cw.syncFundData(context.Background(), ticker, url)
	}

	log.Printf(finishedMsg, time.Since(start))
}

func (cw *CronWrapper) syncFundData(ctx context.Context, ticker, url string) {
	fund, err := cw.ssgaRepo.FetchFundPage(url)
	if err != nil {
		log.Errorf(errMsg, err)

		return
	}

	tx, err := cw.db.NewTXContext(ctx)
	if err != nil {
		log.Errorf(errMsg, err)

		return
	}
	defer tx.Rollback()

	err = tx.DeleteFund(ticker)
	if err != nil {
		log.Errorf(errMsg, err)

		return
	}

	err = tx.InsertFund(fund)
	if err != nil {
		log.Errorf(errMsg, err)

		return
	}

	if err = tx.Commit(); err != nil {
		log.Errorf(errMsg, err)

		return
	}

	return
}
