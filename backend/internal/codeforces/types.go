package codeforces

type User struct {
	Handle    string `json:"handle"`
	Rating    int    `json:"rating"`
	MaxRating int    `json:"maxRating"`
	Rank      string `json:"rank"`
	MaxRank   string `json:"maxRank"`
	Avatar    string `json:"avatar"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Country   string `json:"country,omitempty"`
}
