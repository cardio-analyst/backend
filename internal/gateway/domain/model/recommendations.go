package model

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
	Risk
)

var recommendationTypeNames = []string{
	"Smoking",
	"SBPLevel",
	"BMI",
	"CholesterolLevel",
	"Risk",
}

func (t RecommendationType) String() string {
	if Smoking <= t && t <= Risk {
		return recommendationTypeNames[t]
	}
	buf := make([]byte, 20)
	n := fmtInt(buf, uint64(t))
	return "%!RecommendationType(" + string(buf[n:]) + ")"
}
