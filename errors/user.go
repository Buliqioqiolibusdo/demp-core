package errors

func NewUserError(msg string) (err error) {
	return NewError(ErrorPrefixUser, msg)
}

var (
	ErrorUserInvalidType           = NewUserError("invalid type")
	ErrorUserInvalidToken          = NewUserError("invalid token")
	ErrorUserNotExists             = NewUserError("not exists")
	ErrorUserAlreadyExists         = NewUserError("already exists")
	ErrorUserMismatch              = NewUserError("mismatch")
	ErrorUserMissingRequiredFields = NewUserError("missing required fields")
	ErrorUserUnauthorized          = NewUserError("unauthorized")
	ErrorUserNotExistsInContext    = NewUserError("not exists in context")
	ErrorUserInvalidPassword       = NewUserError("invalid password (length must be no less than 5)")
)
