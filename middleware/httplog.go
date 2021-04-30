package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"time"

	"github.com/ExperienceOne/apikit/internal/framework/xserver"
	"github.com/go-ozzo/ozzo-routing"
	"github.com/pkg/errors"
)

// Log creates a middleware to log HTTP requests and responses.
func Log(log LogFunc) xserver.Middleware {
	return xserver.Middleware{
		Handler: LogHandler(log),
	}
}

// LogFunc represents the func header for logging HTTP requests
// and responses on a middleware.
type LogFunc func(r *http.Request, w *LogResponseWriter, elapsed time.Duration)

// LogHandler creates a middleware handler for logging HTTP requests
// and responses
func LogHandler(log LogFunc) routing.Handler {

	return func(c *routing.Context) error {

		start := time.Now()

		if err := storeReqInCtx(c); err != nil {
			return errors.Wrap(err, "failed to store request in context")
		}

		rw := NewLogResponseWriter(c.Response)
		c.Response = rw

		err := c.Next()

		if data := c.Get("requestData"); data != nil {
			c.Request.Body = ioutil.NopCloser(data.(*bytes.Buffer))
		}

		log(c.Request, rw, time.Since(start))

		return err
	}
}

func storeReqInCtx(c *routing.Context) error {

	if c.Request.ContentLength <= 0 {
		return nil
	}

	buf, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	// TODO: maybe we should add a kind of restriction for big data blobs.
	c.Set("requestData", bytes.NewBuffer(buf))

	return nil
}

type LogResponseStatusWriter struct {
	http.ResponseWriter
	Status int
}

func NewLogResponseStatusWriter(w http.ResponseWriter) *LogResponseStatusWriter {

	return &LogResponseStatusWriter{ResponseWriter: w}
}

func (r *LogResponseStatusWriter) WriteHeader(status int) {

	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

// LogResponseWriter represents a http.ResponseWriter with additional
// functionality to store response related data for a later use.
type LogResponseWriter struct {
	LogResponseStatusWriter
	BytesWritten int64
	Body         []byte
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {

	return &LogResponseWriter{LogResponseStatusWriter: *NewLogResponseStatusWriter(w)}
}

func (r *LogResponseWriter) Write(p []byte) (int, error) {

	// TODO: maybe we should add a kind of restriction for big data blobs.
	r.Body = make([]byte, len(p))
	copy(r.Body, p)

	return r.ResponseWriter.Write(p)
}
