package refactor

import (
	"testing"
)

func TestMv(t *testing.T) {
	cases := []struct {
		name string
		args []string
		ok   bool
		from string
		to   string
	}{
		{
			name: "case_data_single_block",
			args: []string{"data", "data.tf"},
			ok:   true,
			from: "test_data/mv/case_data_single_block/from",
			to:   "test_data/mv/case_data_single_block/to",
		},
		// {
		// 	name: "case_data_single_block_qualified",
		// 	args: []string{"data.a", "data.tf"},
		// 	ok:   true,
		// 	from: "test_data/mv/case_data_single_block/from",
		// 	to:   "test_data/mv/case_data_single_block/to",
		// },
		// {
		// 	name: "case_data_single_block_fully_qualified",
		// 	args: []string{"data.a.b", "data.tf"},
		// 	ok:   true,
		// 	from: "test_data/mv/case_data_single_block/from",
		// 	to:   "test_data/mv/case_data_single_block/to",
		// },
		// {
		// 	name: "case_data_single_block_new_file",
		// 	args: []string{"data", "data.tf"},
		// 	ok:   true,
		// 	from: "test_data/mv/case_data_single_block_new_file/from",
		// 	to:   "test_data/mv/case_data_single_block_new_file/to",
		// },
		// {
		// 	name: "case_all_data_blocks",
		// 	args: []string{"data", "data.tf"},
		// 	ok:   true,
		// 	from: "test_data/mv/case_all_data_blocks/from",
		// 	to:   "test_data/mv/case_all_data_blocks/to",
		// },
		// {
		// 	name: "case_all_vars",
		// 	args: []string{"variable", "variables.tf"},
		// 	ok:   true,
		// 	from: "test_data/mv/case_all_vars/from",
		// 	to:   "test_data/mv/case_all_vars/to",
		// },
		// {
		// 	name: "case_single_var",
		// 	args: []string{"variable.a", "variables.tf"},
		// 	ok:   true,
		// 	from: "test_data/mv/case_single_var/from",
		// 	to:   "test_data/mv/case_single_var/to",
		// },
		// {
		// 	name: "case_all_locals",
		// 	args: []string{"locals", "locals.tf"},
		// 	ok:   true,
		// 	from: "test_data/mv/case_all_locals/from",
		// 	to:   "test_data/mv/case_all_locals/to",
		// },
		// {
		// 	name: "case_single_local",
		// 	args: []string{"local.b", "locals.tf"},
		// 	ok:   true,
		// 	from: "test_data/mv/case_single_local/from",
		// 	to:   "test_data/mv/case_single_local/to",
		// },
		// {
		// 	name: "case_resource_type",
		// 	args: []string{"resource.a", "dest.tf"},
		// 	ok:   true,
		// 	from: "test_data/mv/case_resource_type/from",
		// 	to:   "test_data/mv/case_resource_type/to",
		// },
		// {
		// 	name: "case_all_outputs",
		// 	args: []string{"output", "outputs.tf"},
		// 	ok:   true,
		// 	from: "test_data/mv/case_all_outputs/from",
		// 	to:   "test_data/mv/case_all_outputs/to",
		// },
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan, err := Mv(tc.args[0], tc.args[1], tc.from)
			if (err == nil) != tc.ok {

			}
			assertMockCmdFileOutput(t, tc.name, tc.from, tc.to, plan)
		})
	}
}
