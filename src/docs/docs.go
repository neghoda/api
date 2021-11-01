// API Challenge.
//
// API for UI
//
// Schemes: https, http
// BasePath: /api/v1
// Version: 0.0.1
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
//
// swagger:meta
package docs

// swagger:model
type CommonError struct {
	Code string `json:"code"`
}

// swagger:model
type ValidationErr struct {
	Code             string `json:"code"`
	ValidationErrors struct {
		Field string `json:"field"`
		Code  string `json:"code"`
	} `json:"validation_errors"`
}
