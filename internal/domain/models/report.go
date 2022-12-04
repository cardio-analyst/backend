package models

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

func fmtInt(buf []byte, v uint64) int {
	w := len(buf)
	if v == 0 {
		w--
		buf[w] = '0'
	} else {
		for v > 0 {
			w--
			buf[w] = byte(v%10) + '0'
			v /= 10
		}
	}
	return w
}
