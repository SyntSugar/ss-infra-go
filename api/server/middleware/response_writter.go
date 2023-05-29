package middleware

import (
	"bufio"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	noWritten     = -1            // Constant indicating no bytes written.
	defaultStatus = http.StatusOK // Default HTTP status.
)

// ResponseWriter is an interface combining multiple http related interfaces,
// including http.ResponseWriter, http.Hijacker, and http.Flusher.
// It also provides additional methods to retrieve status, size and first byte time.
type ResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher

	Status() int              // Retrieve the status.
	Size() int                // Retrieve the size.
	FirstByteTime() time.Time // Retrieve the time of the first byte written.
}

type responseWriterImpl struct {
	http.ResponseWriter
	size          int
	status        int
	firstByteTime time.Time
}

// Ensure responseWriterImpl implements ResponseWriter.
var _ ResponseWriter = &responseWriterImpl{}

func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	writer := &responseWriterImpl{ResponseWriter: w}
	writer.reset(w)
	return writer
}

// reset the responseWriterImpl with the provided http.ResponseWriter.
func (w *responseWriterImpl) reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.size = noWritten
	w.status = defaultStatus
}

// WriteHeader writes the status code and updates the first byte time if not already set.
func (w *responseWriterImpl) WriteHeader(code int) {
	if code > 0 && w.status != code {
		w.ResponseWriter.WriteHeader(code)
		if w.firstByteTime.IsZero() {
			w.status = code
			w.firstByteTime = time.Now()
		}
	}
}

// Write writes the data and updates the size and first byte time if not already set.
func (w *responseWriterImpl) Write(data []byte) (n int, err error) {
	if w.firstByteTime.IsZero() {
		w.WriteHeader(http.StatusOK)
	}
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
}

func (w *responseWriterImpl) Status() int {
	return w.status
}

func (w *responseWriterImpl) Size() int {
	if w.size == noWritten {
		return 0
	}
	return w.size
}

// FirstByteTime returns the time when the first byte is written.
func (w *responseWriterImpl) FirstByteTime() time.Time {
	return w.firstByteTime
}

// Hijack hijacks the connection.
func (w *responseWriterImpl) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.size < 0 {
		w.size = 0
	}
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

// Flush sends any buffered data to the client.
func (w *responseWriterImpl) Flush() {
	w.ResponseWriter.(http.Flusher).Flush()
}

type ginResponseWriter struct {
	gin.ResponseWriter
	fbt time.Time
}

func (responseWriter *ginResponseWriter) FirstByteTime() time.Time {
	return responseWriter.fbt
}

func (responseWriter *ginResponseWriter) WriteHeaderNow() {
	responseWriter.ResponseWriter.WriteHeaderNow()
	if responseWriter.fbt.IsZero() {
		responseWriter.fbt = time.Now()
	}
}

func (responseWriter *ginResponseWriter) Write(data []byte) (n int, err error) {
	responseWriter.WriteHeaderNow()
	return responseWriter.ResponseWriter.Write(data)
}

// WriteString writes a string and calls WriteHeaderNow before writing the data```
func (responseWriter *ginResponseWriter) WriteString(s string) (n int, err error) {
	responseWriter.WriteHeaderNow()
	return responseWriter.ResponseWriter.WriteString(s)
}

// newResponseWriter creates a new ResponseWriter from a gin.ResponseWriter.
func newResponseWriter(writer gin.ResponseWriter) ResponseWriter {
	return &ginResponseWriter{ResponseWriter: writer}
}
