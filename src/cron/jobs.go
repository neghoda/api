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

func (cw *Wrapper) syncFundsData() {
	start := time.Now()

	mapFundTickerURL, err := cw.ssgaRepo.FetchFundFinderPage()
	if err != nil {
		log.Errorf(errMsg, err)
	}

	for _, url := range mapFundTickerURL {
		cw.syncFundData(context.Background(), url)
	}

	log.Printf(finishedMsg, time.Since(start))
}

func (cw *Wrapper) syncFundData(ctx context.Context, url string) {
	fund, err := cw.ssgaRepo.FetchFundPage(url)
	if err != nil {
		log.Errorf(errMsg, err)

		return
	}

	err = cw.fundRepo.ReplaceFund(ctx, &fund)
	if err != nil {
		log.Errorf(errMsg, err)

		return
	}
}
