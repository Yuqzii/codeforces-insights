package codeforces

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalContest(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Contest
	}{
		{
			name: "no div",
			input: `{
				"id": 2067,
				"name": "Test Contest",
				"startTimeSeconds": 1770984000,
				"durationSeconds": 7200,
				"phase": "FINISHED"
			}`,
			want: Contest{
				ID:        2067,
				Name:      "Test Contest",
				StartTime: time.Date(2026, time.February, 13, 12, 0, 0, 0, time.UTC),
				Duration:  7200,
				Phase:     "FINISHED",
			},
		}, {
			name: "single div",
			input: `{
				"id": 1420,
				"name": "Codeforces Round 1000 (Div. 2)"
			}`,
			want: Contest{
				ID:   1420,
				Name: "Codeforces Round 1000 (Div. 2)",
				Div:  2,
			},
		}, {
			name: "multiple divs",
			input: `{
				"id": 2069,
				"name": "Codeforces Round 1069 (Div. 1 + Div. 2)"
			}`,
			want: Contest{
				ID:   2069,
				Name: "Codeforces Round 1069 (Div. 1 + Div. 2)",
				Div:  2,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var got Contest
			err := json.Unmarshal([]byte(test.input), &got)
			require.Nil(t, err)

			assert.EqualExportedValues(t, test.want, got)
		})
	}
}
