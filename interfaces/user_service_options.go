package interfaces

type UserCreateOptions struct {
	Username string
	Password string
	Email    string
	Role     string
	Name     string
}

type UserLoginOptions struct {
	Username string
	Password string
}
