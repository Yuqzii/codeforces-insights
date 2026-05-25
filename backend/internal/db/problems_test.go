package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuqzii/codeforces-insights/internal/codeforces"
)

func TestCorrectProblemIndices(t *testing.T) {
	tests := []struct {
		name  string
		input []probWithDiv
		wants []string
	}{
		{
			name: "single contest",
			input: []probWithDiv{{
				Problem: codeforces.Problem{Index: "A"},
				Div:     2,
			}, {
				Problem: codeforces.Problem{Index: "B"},
				Div:     2,
			}, {
				Problem: codeforces.Problem{Index: "C"},
				Div:     2,
			}},
			wants: []string{"A", "B", "C"},
		}, {
			name: "multiple contests",
			input: []probWithDiv{{
				Problem: codeforces.Problem{Index: "A"},
				Div:     2,
			}, {

				Problem: codeforces.Problem{Index: "B"},
				Div:     2,
			}, {
				Problem: codeforces.Problem{Index: "A"},
				Div:     1,
			}, {
				Problem: codeforces.Problem{Index: "B"},
				Div:     1,
			}},
			wants: []string{"A", "B", "C", "D"},
		}, {
			name: "multiple contests with problem numbers",
			input: []probWithDiv{{
				Problem: codeforces.Problem{Index: "A1"},
				Div:     3,
			}, {
				Problem: codeforces.Problem{Index: "A2"},
				Div:     3,
			}, {
				Problem: codeforces.Problem{Index: "A"},
				Div:     2,
			}, {
				Problem: codeforces.Problem{Index: "B1"},
				Div:     2,
			}, {
				Problem: codeforces.Problem{Index: "B2"},
				Div:     2,
			}},
			wants: []string{"A1", "A2", "B", "C1", "C2"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mp := make(map[int][]codeforces.Problem)
			_, got, err := correctProblemIndices(test.input, mp)
			require.Nil(t, err)

			for _, probs := range got {
				for i, p := range probs {
					assert.Equal(t, test.wants[i], p.Index, "problem %d has wrong index", i)
				}
			}
		})
	}
}
