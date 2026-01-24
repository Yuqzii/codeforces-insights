package codeforces

type Problem struct {
	Name      string   `json:"name" db:"name"`
	ContestID int      `json:"contestId,omitempty" db:"contest_id"`
	Index     string   `json:"index" db:"index"`
	Rating    int      `json:"rating" db:"rating"`
	Tags      []string `json:"tags" db:"tags"`
}

func (p *Problem) Hash() int64 {
	return HashProblemIndex(p.ContestID, p.Index)
}

func HashProblemIndex(contestID int, index string) int64 {
	var res int64
	res = int64(contestID) << 32
	res |= 1 << (index[0] - 'A' + 2) // One bit for indicating the character

	if len(index) > 1 {
		res |= int64(index[1] - '1') // Last two bits indicating number
	}

	return res
}
