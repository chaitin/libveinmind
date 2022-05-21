package behaviour

import (
	"github.com/chaitin/libveinmind/go/pkg/binding"
)

// Closer annotates handle to have io.Closer behaviour.
type Closer struct {
	h *binding.Handle
}

func (c *Closer) Close() error {
	if c.h.IsNil() {
		return nil
	}
	// Specially, the underlying handle will be freed
	// whenever the caller is attempting to close it.
	defer c.h.Free()
	return c.h.Close()
}

func NewCloser(h *binding.Handle) Closer {
	return Closer{h: h}
}

// Reader annotates handle to have io.Reader behaviour.
type Reader struct {
	h *binding.Handle
}

func (r *Reader) Read(b []byte) (int, error) {
	buf := binding.NewBytes(b)
	defer buf.Free()
	size, err := r.h.Read(buf)
	copy(b, buf.Bytes()[:])
	return int(size), err
}

func NewReader(h *binding.Handle) Reader {
	return Reader{h: h}
}

// ReaderAt annotates handle to have io.ReaderAt behaviour.
type ReaderAt struct {
	h *binding.Handle
}

func (r *ReaderAt) ReadAt(b []byte, off int64) (int, error) {
	buf := binding.NewBytes(b)
	defer buf.Free()
	size, err := r.h.ReadAt(buf, off)
	copy(b, buf.Bytes()[:])
	return int(size), err
}

func NewReaderAt(h *binding.Handle) ReaderAt {
	return ReaderAt{h: h}
}

// Writer annotates handle to have io.Writer behaviour.
type Writer struct {
	h *binding.Handle
}

func (w *Writer) Write(b []byte) (int, error) {
	buf := binding.NewBytes(b)
	defer buf.Free()
	size, err := w.h.Write(buf)
	return int(size), err
}

func NewWriter(h *binding.Handle) Writer {
	return Writer{h: h}
}

// Writer annotates handle to have io.WriterAt behaviour.
type WriterAt struct {
	h *binding.Handle
}

func (w *WriterAt) WriteAt(b []byte, off int64) (int, error) {
	buf := binding.NewBytes(b)
	defer buf.Free()
	size, err := w.h.WriteAt(buf, off)
	return int(size), err
}

func NewWriterAt(h *binding.Handle) WriterAt {
	return WriterAt{h: h}
}

// Seeker annotates handle to have io.Seeker behaviour.
type Seeker struct {
	h *binding.Handle
}

func (s *Seeker) Seek(off int64, whence int) (int64, error) {
	return s.h.Seek(off, whence)
}

func NewSeeker(h *binding.Handle) Seeker {
	return Seeker{h: h}
}
