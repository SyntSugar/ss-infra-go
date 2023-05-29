package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponseWriterImpl(t *testing.T) {
	w := httptest.NewRecorder()
	rw := NewResponseWriter(w)

	assert.Equal(t, defaultStatus, rw.Status())
	assert.Equal(t, 0, rw.Size())

	rw.WriteHeader(http.StatusNotFound)
	assert.Equal(t, http.StatusNotFound, rw.Status())

	data := []byte("Hello, World!")
	_, err := rw.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data), rw.Size())

	fbt := rw.FirstByteTime()
	assert.NotEqual(t, time.Time{}, fbt)
}
