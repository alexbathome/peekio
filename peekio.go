// Package peekio provides a common interface for Peeking into a stream of data.
//
// Additionally, it provides a convenience io.Reader implementation such that
// a peekio.Peeker can be adapted for use where an io.Reader is required.
package peekio

// Peeker is an interface that allows peeking into a stream of data without consuming it.
type Peeker interface {
	Peek(n int) ([]byte, error)
}

// PeekReader is an io.Reader that adapts a Peeker to be passed into a caller
// that expects an io.Reader.
//
// It behaves the same as an io.Reader, where the underlying buffer appears to
// be consumed, however in actuality, PeekReader is just advancing an offset.
type PeekReader struct {
	offset int
	peeker Peeker
}

// NewPeekReader creates a new PeekReader that wraps the provided Peeker.
//
// Example:
//
//	r := bufio.NewReader(strings.NewReader("Hello, World!"))
//	pr := peekio.NewPeekReader(r)
//
//	buf := make([]byte, 5)
//	_, _ = pr.Read(buf)
//	println("buf:", string(buf)) // "buf: Hello"
//
//	_, _ = pr.Read(buf)
//	println("buf:", string(buf)) // "buf: , Wor"
//
//	// the underlying reader has not been consumed
//	_, _ = r.Read(buf)
//	println("buf:", string(buf)) // "buf: Hello"
//	// Output:
//	//	buf: Hello
//	//	buf: , Wor
//	//	buf: Hello
func NewPeekReader(peeker Peeker) *PeekReader {
	return &PeekReader{
		offset: 0,
		peeker: peeker,
	}
}

// Read implements the io.Reader interface for PeekReader.
// It reads up to len(p) bytes into p from the underlying Peeker.
// It returns the number of bytes read and any error encountered.
func (pr *PeekReader) Read(p []byte) (int, error) {
	var (
		peeked, err = pr.peeker.Peek(pr.offset + len(p))
		next        = peeked[pr.offset:]
		n           = copy(p, next)
	)
	pr.offset += n
	if n > 0 {
		return n, nil
	}
	return n, err
}
