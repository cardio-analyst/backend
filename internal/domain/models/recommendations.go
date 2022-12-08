package models

type Recommendation struct {
	What string `json:"what"`
	Why  string `json:"why"`
	How  string `json:"how"`
}

type RecommendationType int

// possible recommendation types
const (
	Smoking RecommendationType = iota
	SBPLevel
	BMI
	CholesterolLevel
)

var recommendationTypeNames = []string{
	"None",
	"Smoking",
	"SBPLevel",
	"BMI",
	"CholesterolLevel",
}

func (t RecommendationType) String() string {
	if Smoking <= t && t <= CholesterolLevel {
		return recommendationTypeNames[t]
	}
	buf := make([]byte, 20)
	n := fmtInt(buf, uint64(t))
	return "%!RecommendationType(" + string(buf[n:]) + ")"
}
