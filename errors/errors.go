package errors

const (
	StatusUnauthorized        = "UNAUTHORIZED"
	StatusInternalServerError = "INTERNAL_SERVER_ERROR"
)

const (
	StatusUnauthorizedText        = "Unauthorized"
	StatusWrongPasswordText       = "Wrong Password"
	StatusInternalServerErrorText = "Internal Server Error"
)

type ErrorWithCode struct {
	Message string
	Code    string
}

func (err *ErrorWithCode) Error() string {
	return err.Message
}

func (err *ErrorWithCode) Extensions() map[string]interface{} {
	return map[string]interface{}{
		"code": err.Code,
	}
}
