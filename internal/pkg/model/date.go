package model

import (
	"fmt"
	"strings"
	"time"
)

const DateLayout = "02.01.2006"

// Date represents the date with DateLayout layout.
type Date struct {
	time.Time
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
