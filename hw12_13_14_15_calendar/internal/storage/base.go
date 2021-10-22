package storage

import "errors"

var ErrNotFound = errors.New("not found")

var ErrDateBusy = errors.New("date busy")

var ErrStartDateMoreThanEndDate = errors.New("start date more than end date")

var ErrInvalidInitialization = errors.New("invalid initialization struct")
