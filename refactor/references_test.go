package refactor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReferences(t *testing.T) {
	cases := []struct {
		name  string
		addr  string
		src   string
		count int
	}{
		{
			name: "one_ref",
			addr: "a.b.c",
			src: `
			b {
				a = a.b.c
				b = a.b.d
			}
			`,
			count: 1,
		},
		{
			name: "two_refs",
			addr: "a.b.c",
			src: `
			b {
				a = a.b.c
			}
			d {
				z = a.b.c
			}
			`,
			count: 2,
		},
		{
			name: "ref_in_nested_block",
			addr: "a.b.c",
			src: `
			b {
				d {
					z = a.b.c
				}
			}
			`,
			count: 1,
		},
		{
			name: "addr_shorter",
			addr: "a.b",
			src: `
			b {
				y = a.b.d
				z = a.b.c
			}
			`,
			count: 2,
		},
		{
			name: "addr_longer",
			addr: "a.b.c.d",
			src: `
			b {
				z = a.b.c
			}
			`,
			count: 0,
		},
		{
			name: "interpolation",
			addr: "a.b",
			src: `
			b {
				z = "prefix${a.b.c}"
			}
			`,
			count: 1,
		},
	}
	for _, tc := range cases {
		hclFile, err := ParseHclBytes([]byte(tc.src), "test.tf")
		if err != nil {
			t.Fatalf("unexpected err = %s", err)
		}
		found := findReferencingExpresssions(hclFile.Body(), ParseAddress(tc.addr))
		assert.Equal(t, tc.count, len(found), "Case %v", tc.name)
	}
}
