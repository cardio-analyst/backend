package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

const DateLayout = "02.01.2006"

// Date represents the date with DateLayout layout.
type Date struct {
	time.Time
}

func (d *Date) Validate(value interface{}) error {
	date, ok := value.(Date)
	if !ok {
		return errors.New("cannot cast to date")
	}
	return validation.Validate(date.String(), validation.Required, validation.Date(DateLayout))
}

func (d *Date) String() string {
	return fmt.Sprintf("%v", d.Format(DateLayout))
}

func (d *Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf("%q", d.String())), nil
}

func (d *Date) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return
	}
	d.Time, err = time.Parse(DateLayout, s)
	return
}
