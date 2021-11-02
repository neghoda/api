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

// swagger:operation GET /ticker funds get_ticker_list
//   list of all avalible tickers
// ---
// parameters:
// responses:
//   '201':
//     description: Created
//     schema:
//       "$ref": "#/definitions/TokenPair"
//   '400':
//     description: Bad Request
//     schema:
//       "$ref": "#/definitions/ValidationErr"
//   '500':
//     description: Internal Server Error
//     schema:
//       "$ref": "#/definitions/CommonError"
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

// swagger:operation GET /funds funds get_fund_data
//   list of all avalible tickers
// ---
// parameters:
// - name: ticker
//   in: query
//   required: true
//   type: string
// responses:
//   '201':
//     description: Created
//     schema:
//       "$ref": "#/definitions/TokenPair"
//   '400':
//     description: Bad Request
//     schema:
//       "$ref": "#/definitions/ValidationErr"
//   '500':
//     description: Internal Server Error
//     schema:
//       "$ref": "#/definitions/CommonError"
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
