package diskinfo

import (
	"bufio"
	"io"
)

// LineReader reads a reader by line
type LineReader struct {
	r    *bufio.Reader
	line string
}

// NewLineReader makes a new LineReader of an io.Reader
func NewLineReader(r io.Reader) *LineReader {
	return &LineReader{r: bufio.NewReader(r)}
}

// ReadLine returns the next line, if line is empty, you should skip the iteration.
//you shuold check errors on every line.
func (l *LineReader) ReadLine() (string, error) {

	line, isPrefix, err := l.r.ReadLine()

	ret := ""

	if isPrefix {
		l.line += string(line)
	} else if l.line != "" {
		ret = l.line
		l.line = ""
	} else {
		ret = string(line)
	}

	return ret, err
}
