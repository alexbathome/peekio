// Package peekio provides a common interface for Peeking into a stream of data.
//
// Additionally, it provides a convenience io.Reader implementation such that
// a peekio.Peeker can be adapted for use where an io.Reader is required.
package peekio

// Peeker is an interface that allows peeking into a stream of data without consuming it.
type Peeker interface {
	// Peek returns the next n bytes from the stream without advancing the read position.
	// If fewer than n bytes are available, it returns all available bytes.
	// If the end of the stream is reached, it returns an error.
	//
	// The returned slice is only valid until the next call to Peek or Read.
	// The caller should not modify the contents of the returned slice.
	Peek(n int) ([]byte, error)
}

type PeekReader struct {
	offset int
	peeker Peeker
}

// NewPeekReader creates a new PeekReader that wraps the provided Peeker.
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
