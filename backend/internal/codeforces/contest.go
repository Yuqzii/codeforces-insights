package codeforces

type Contestant struct {
	Rank    int     `json:"rank"`
	Points  float64 `json:"points"`
	Penalty int     `json:"penalty"`
	Rating  int
}
