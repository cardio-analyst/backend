package model

import (
	"errors"
	"fmt"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const DatetimeLayout = "15:04:05 02.01.2006"

// Datetime represents the datetime with DatetimeLayout layout.
type Datetime struct {
	time.Time
}

func (d *Datetime) Validate(value interface{}) error {
	datetime, ok := value.(Datetime)
	if !ok {
		return errors.New("cannot cast to datetime")
	}
	return validation.Validate(datetime.String(), validation.Required, validation.Date(DatetimeLayout))
}

func (d *Datetime) String() string {
	return fmt.Sprintf("%v", d.Format(DatetimeLayout))
}

func (d *Datetime) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf("%q", d.String())), nil
}

func (d *Datetime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return
	}
	d.Time, err = time.Parse(DatetimeLayout, s)
	return
}
