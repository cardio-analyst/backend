package model

import "errors"

var ErrNotEnoughDataToCompileReport = errors.New("not enough data to compile a report")

type ReportType int

// possible report types
const (
	PDF ReportType = iota
	Excel
	Word
)

var reportTypeNames = []string{
	"PDF",
	"Excel",
	"Word",
}

func (t ReportType) String() string {
	if PDF <= t && t <= Word {
		return reportTypeNames[t]
	}
	buf := make([]byte, 20)
	n := fmtInt(buf, uint64(t))
	return "%!ReportType(" + string(buf[n:]) + ")"
}
