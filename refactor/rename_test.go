package refactor

import (
	"testing"
)

func TestTfRenameWithReferencesFilter(t *testing.T) {
	cases := []struct {
		name string
		src  string
		from string
		to   string
		ok   bool
		want string
	}{
		{
			name: "data",
			src: `a0 = data.a.b
data "a" "b" {
  a2 = v2
}

b2 "l2" {
}
`,
			from: "data.a.b",
			to:   "data.a.c",
			ok:   true,
			want: `a0 = data.a.c
data "a" "c" {
  a2 = v2
}

b2 "l2" {
}
`,
		},
		{
			name: "resource",
			src: `a0 = a.b
resource "a" "b" {
  a2 = v2
}

b2 "l2" {
}
`,
			from: "a.b",
			to:   "a.c",
			ok:   true,
			want: `a0 = a.c
resource "a" "c" {
  a2 = v2
}

b2 "l2" {
}
`,
		},
		{
			name: "resource2",
			src: `a0 = a.b
resource "aws_iam_policy" "other" {
}
resource "aws_iam_policy" "read_only_restrictions" {
  a2 = v2
}

b2 "l2" {
}
`,
			from: "aws_iam_policy.read_only_restrictions",
			to:   "aws_iam_policy.read_only_restrictions2",
			ok:   true,
			want: `a0 = a.b
resource "aws_iam_policy" "other" {
}
resource "aws_iam_policy" "read_only_restrictions2" {
  a2 = v2
}

b2 "l2" {
}
`,
		},
		{
			name: "var",
			src: `a0 = var.a
variable "a" {
  a2 = v2
}

b2 "l2" {
}
`,
			from: "var.a",
			to:   "var.b",
			ok:   true,
			want: `a0 = var.b
variable "b" {
  a2 = v2
}

b2 "l2" {
}
`,
		},
		{
			name: "var_as_resource_addr_value",
			src: `
variable "a" {
  a2 = v2
}

b2 "l2" {
  a0 = var.a
}
`,
			from: "var.a",
			to:   "var.b",
			ok:   true,
			want: `
variable "b" {
  a2 = v2
}

b2 "l2" {
  a0 = var.b
}
`,
		},
		{
			name: "var_in_interpolation",
			src: `
variable "a" {
}

b2 "l2" {
  a0 = "pre${var.a}"
}
`,
			from: "var.a",
			to:   "var.b",
			ok:   true,
			want: `
variable "b" {
}

b2 "l2" {
  a0 = "pre${var.b}"
}
`,
		},
		{
			name: "moved_block_change_to",
			src: `
resource "aws_instance" "a" {
  a0 = "pre${var.a}"
}
moved {
  from = aws_instance.z
  to   = aws_instance.a
}
`,
			from: "aws_instance.a",
			to:   "aws_instance.b",
			ok:   true,
			want: `
resource "aws_instance" "b" {
  a0 = "pre${var.a}"
}
moved {
  from = aws_instance.z
  to   = aws_instance.b
}
`,
		},
		{
			name: "moved_block_dont_change_from",
			src: `
resource "aws_instance" "a" {
  a0 = "pre${var.a}"
}
moved {
  from = aws_instance.z
  to   = aws_instance.a
}
`,
			from: "aws_instance.z",
			to:   "aws_instance.y",
			ok:   true,
			want: `
resource "aws_instance" "a" {
  a0 = "pre${var.a}"
}
moved {
  from = aws_instance.z
  to   = aws_instance.a
}
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			hclFile, err := ParseHclBytes([]byte(tc.src), "test.tf")
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}
			err = RenameInFile("test.tf", hclFile, ParseAddress(tc.from), ParseAddress(tc.to))
			output := hclFile.Bytes()
			if tc.ok && err != nil {
				t.Fatalf("unexpected err = %s", err)
			}

			got := string(output)
			if !tc.ok && err == nil {
				t.Fatalf("expected to return an error, but no error, outStream: \n%s", got)
			}

			if got != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}
