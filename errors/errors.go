package errors

import "errors"

var ErrCarCreationMessage = errors.New("unable to create a car without id")
var ErrCarUpdatingMessage = errors.New("unable to update a car without id")
