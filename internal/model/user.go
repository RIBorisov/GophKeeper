package model

// UserCredentials is a struct for the user register and/or login calls.
type UserCredentials struct {
	Login    string
	Password string
}

type key int

const CtxUserIDKey key = iota
