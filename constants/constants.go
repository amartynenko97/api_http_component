package constants

type CustomErrorType string

const (
	AccountNotHaveBalance CustomErrorType = "AccountNotHaveBalance"
)

type CustomError struct {
	Type CustomErrorType
}

func (e *CustomError) Error() string {
	return string(e.Type)
}
