package handlers

import (
	"net/http"

	"github.com/erp/api/src/models"
	"github.com/erp/api/src/service"
)

type FundHandler struct {
	service *service.Service
}

func NewFundHandler(s *service.Service) *FundHandler {
	return &FundHandler{
		service: s,
	}
}

func (h *FundHandler) GetTickerList(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.TickerList()
	if err != nil {
		SendHTTPError(w, err)

		return
	}

	SendResponse(w, http.StatusCreated, models.TickersListResponse{
		Tickers: res,
	})
}

func (h *FundHandler) GetFundByTicker(w http.ResponseWriter, r *http.Request) {
	ticker := r.URL.Query().Get("ticker")

	if ticker == "" {
		SendEmptyResponse(w, http.StatusBadRequest)
	}

	res, err := h.service.FundByTicker(r.Context(), ticker)
	if err != nil {
		SendHTTPError(w, err)

		return
	}

	SendResponse(w, http.StatusCreated, res)
}
