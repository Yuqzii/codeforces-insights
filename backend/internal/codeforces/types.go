package codeforces

type User struct {
	Handle     string `json:"handle"`
	Rating     int    `json:"rating"`
	MaxRating  int    `json:"maxRating"`
	Rank       string `json:"rank"`
	MaxRank    string `json:"maxRank"`
	TitlePhoto string `json:"titlePhoto"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	Country    string `json:"country,omitempty"`
}

type Problem struct {
	Name      string   `json:"name"`
	ContestID int      `json:"contestId,omitempty"`
	Index     string   `json:"index"`
	Rating    int      `json:"rating"`
	Tags      []string `json:"tags"`
}

type Submission struct {
	ID                  int     `json:"id"`
	Verdict             string  `json:"verdict"`
	Problem             Problem `json:"problem"`
	ProgrammingLanguage string  `json:"programmingLanguage"`
}
