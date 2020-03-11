package errors

const (
	ForbiddenCode           = "FORBIDDEN"
	InternalServerErrorCode = "INTERNAL_SERVER_ERROR"
	InvalidParamsCode       = "INVALID_PARAMS"
)

const (
	ForbiddenCodeMessage         = "Forbidden"
	InternalServerErrorMessage   = "Internal Server Error"
	InvalidTokenErrorCodeMessage = "Invalid Token"
	InvalidLoginErrorMessage     = "Invalid Login"
	WrongPasswordErrorMessage    = "Wrong Password"
	PostSlugExistMessage         = "Post Slug Exist"
	LangCodeExistMessage         = "Lang Code Exist"
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
