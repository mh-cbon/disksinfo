package diskinfo

type parseTable struct {
	in        string
	expectErr error
	expectOut []*Properties
	path      string // for ls test
}
