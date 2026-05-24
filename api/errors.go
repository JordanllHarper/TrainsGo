package main

import (
	"errors"
)

var (
	errorNotFound      = errors.New("doesn't exist")
	errorAlreadyExists = errors.New("already exists")
)
