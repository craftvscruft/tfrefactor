package refactor

import (
	"testing"
)

func TestExtract(t *testing.T) {
	cases := []struct {
		name string
		args []string
		ok   bool
		from string
		to   string
	}{
		{
			name: "case_data_single_block",
			args: []string{"a.a", "mymodule"},
			ok:   true,
			from: "test_data/extract_module/case_one_resource/from",
			to:   "test_data/extract_module/case_one_resource/to",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan, err := Extract(tc.args[0:], tc.args[len(tc.args)-1], tc.from)
			if (err == nil) != tc.ok {

			}
			assertMockCmdFileOutput(t, tc.name, tc.from, tc.to, plan)
		})
	}
}
