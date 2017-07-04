package utils

import "testing"

func TestGetWords(t *testing.T) {
	var tests = []struct {
		len, want int
	}{
		{2, 105},
		{3, 1079},
		{4, 4212},
	}

	for _, test := range tests {
		got := GetWordsFromWordGameDict(test.len)
		if len(got) != test.want {
			t.Errorf("GetWordsFromWordGameDict(%d) == %d, want %d ", test.len, len(got), test.want)
		}
	}
}
