package recovery

import "testing"

func TestClassWorkspaceInitializeKnownAnswers(t *testing.T) {
	cases := []struct {
		word, seed uint16
		product    [9]uint16
	}{
		{2, 2, [9]uint16{27, 64, 125, 216, 343, 512, 729, 1000, 1331}},
		{0, 0, [9]uint16{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{0x7ffe, 0x7ffe, [9]uint16{0, 1, 128, 30180, 16384, 6732, 29301, 18473, 64}},
		{65535, 1, [9]uint16{4, 9, 16, 25, 36, 49, 64, 81, 100}},
	}
	for _, tc := range cases {
		got := classWorkspaceInitialize(tc.word)
		if got.Seed != tc.seed || got.Product != tc.product {
			t.Fatalf("word %d: initializer differs", tc.word)
		}
		for _, v := range got.Current {
			if v != tc.seed {
				t.Fatalf("word %d: current plane differs", tc.word)
			}
		}
	}
}
