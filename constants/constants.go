package constants

type ErrorResponse struct {
	Error string `json:"error"`
}

type CustomErrorType string

const (
	NoSuchCurrency CustomErrorType = "NoSuchCurrency"
)
