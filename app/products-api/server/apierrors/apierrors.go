package apierrors

type ApiError struct {
	Code    int    `json:"code"`
	Err     error  `json:"error"`
	Message string `json:"message"`
}

func (err ApiError) Error() string {
	if err.Err != nil {
		return err.Err.Error()
	}
	return err.Message
}

func NewApiError(message string, code int, err error) error {
	return &ApiError{
		Message: message,
		Code:    code,
		Err:     err,
	}
}