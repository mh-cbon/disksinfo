package diskinfo

import (
	"bufio"
	"bytes"
	"testing"
)

func TestWmicParser(t *testing.T) {

	testsTable := []parseTable{
		parseTable{
			//todo need to fix this output with a real one from windows.
			in: `
Caption  Description       FileSystem
C:       Local Fixed Disk  NTFS
`,
			expectErr: nil,
			expectOut: []*Properties{
				&Properties{
					Path:      "C:",
					MountPath: "C:",
				},
			},
		},
	}

	for i, testTable := range testsTable {

		var b bytes.Buffer
		r := NewWmicReader(bufio.NewReader(&b))
		b.WriteString(testTable.in)

		res, err := r.Read()
		if err != nil && testTable.expectErr != err {
			t.Fatalf("Test(%v): Unexpected error %v", i, err)
		}

		for _, p := range res {
			found := PropertiesList(testTable.expectOut).FindByPath(p.Path)
			if found == nil {
				t.Errorf("Test(%v): Unexpected property %q not found\n%#v\ntestTable.in=\n%v", i, p.Path, p, testTable.in)
			}
		}

		for _, p := range testTable.expectOut {
			found := PropertiesList(res).FindByPath(p.Path)
			if found == nil {
				t.Errorf("Test(%v): Property %q not found\n%#v\ntestTable.in=\n%v", i, p.Path, p, testTable.in)
			}
		}
	}
}
