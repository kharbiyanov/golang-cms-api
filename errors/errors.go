package errors

const (
	ForbiddenCode           = "FORBIDDEN"
	InternalServerErrorCode = "INTERNAL_SERVER_ERROR"
	InvalidParamsCode       = "INVALID_PARAMS"
)

const (
	ForbiddenCodeMessage          = "Forbidden"
	InternalServerErrorMessage    = "Internal Server Error"
	InvalidTokenErrorCodeMessage  = "Invalid Token"
	InvalidLoginErrorMessage      = "Invalid Login"
	WrongPasswordErrorMessage     = "Wrong Password"
	PostSlugExistMessage          = "Post Slug Exist"
	TermSlugExistMessage          = "Term Slug Exist"
	TermNotFoundMessage           = "Term Not Found"
	TermParentIDNotFoundMessage   = "Term Parent ID Not Found"
	LangCodeExistMessage          = "Lang Code Exist"
	LangNotFoundMessage           = "Lang Not Found"
	MenuItemNotFoundMessage       = "Menu Item Not Found"
	MenuNotFoundMessage           = "Menu Not Found"
	UserNameExistMessage          = "Username Exist"
	EmailExistMessage             = "Email Exist"
	UserNotActivatedMessage       = "User Not Activated"
	InvalidActivationCodeMessage  = "Invalid Activation Code"
	ActivationCodeNotFoundMessage = "Activation Code Not Found"
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
