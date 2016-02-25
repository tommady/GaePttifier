package gaepttifer

import (
	"fmt"
	"time"
)

type GaePttiferError struct {
	When time.Time
	What string
	Err  error
}

func (e GaePttiferError) Error() string {
	return fmt.Sprintf("%v -> %v -> %v", e.When, e.What, e.Err)
}

func ReportError(what string, err error) error {
	return GaePttiferError{
		time.Now(),
		what,
		err,
	}
}
