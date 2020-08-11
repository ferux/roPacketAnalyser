package rpa

type GeneralError string

func (err GeneralError) Error() string { return string(err) }

type NotFoundError string

func (err NotFoundError) Error() string { return "not found: " + string(err) }

const (
	ErrHandlerNotFound NotFoundError = "handler"
)
