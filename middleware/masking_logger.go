package middleware

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
)

type (
	logRequest struct {
		Body      interface{} `json:",omitempty"`
		Params    interface{} `json:",omitempty"`
		RequestID string      `json:"requestId,omitempty"`
		Endpoint  string
		Method    string
	}

	logResponse struct {
		Body       interface{} `json:",omitempty"`
		StatusCode int         `json:"responseCode"`
	}

	LogEntry struct {
		Request   *logRequest  `json:",omitempty"`
		Response  *logResponse `json:",omitempty"`
		Duration  float64      `json:"durationMs"`
		Timestamp string
	}
)

// NewLogger returns a new Logger based on APIKit middleware.LogFunc by
// recognizing the given time format.
func NewLogger(timeFormat string, pathsToIgnore []string, logFunc func(entity LogEntry, values ...interface{})) LogFunc {
	return NewMaskingLogger(timeFormat, nil, nil, pathsToIgnore, logFunc)
}

// NewMaskingLogger returns a new Logger based on APIKit middleware.LogFunc that
// will mask the body content of HTTP requests and responses for a given array
// of field names and types.
func NewMaskingLogger(timeFormat string, fieldsToMask []string, typesToMask []reflect.Type, pathsToIgnore []string, logFunc func(logEntry LogEntry, values ...interface{})) LogFunc {

	return func(r *http.Request, w *LogResponseWriter, elapsed time.Duration) {

		entry := &LogEntry{
			Duration:  float64(elapsed.Nanoseconds()) / 1e6,
			Timestamp: time.Now().UTC().Format(timeFormat),
		}

		req := &logRequest{
			Endpoint:  r.URL.Path,
			Method:    r.Method,
			Params:    r.URL.Query(),
			RequestID: GetRequestID(r.Context()),
		}

		if !isIgnoredPath(pathsToIgnore, req.Endpoint) {
			if buf, err := ioutil.ReadAll(r.Body); err == nil {
				if strings.Contains(r.Header.Get("Content-Type"), "application/json") {
					if len(buf) > 0 {
						switch buf[0] {
						case '[':
							var body []interface{}
							if err := json.Unmarshal(buf, &body); err == nil {
								if err := maskBody(&body, fieldsToMask, typesToMask); err == nil {
									req.Body = body
								}
							}
						case '{':
							var body map[string]interface{}
							if err := json.Unmarshal(buf, &body); err == nil {
								if err := maskBody(&body, fieldsToMask, typesToMask); err == nil {
									req.Body = body
								}
							}
						default:
							req.Body = string(w.Body)
						}
					}
				} else {
					req.Body = string(buf)
				}
			}
		}
		entry.Request = req

		resp := &logResponse{StatusCode: w.Status}

		if !isIgnoredPath(pathsToIgnore, req.Endpoint) {
			if len(w.Body) > 0 {
				switch w.Body[0] {
				case '[':
					var body []interface{}
					if err := json.Unmarshal(w.Body, &body); err == nil {
						if err := maskBody(&body, fieldsToMask, typesToMask); err == nil {
							resp.Body = body
						}
					}
				case '{':
					var body map[string]interface{}
					if err := json.Unmarshal(w.Body, &body); err == nil {
						if err := maskBody(&body, fieldsToMask, typesToMask); err == nil {
							resp.Body = body
						}
					}
				default:
					resp.Body = string(w.Body)
				}
			}
		}

		entry.Response = resp
		logFunc(*entry, "request handled")
	}
}

func maskBody(body interface{}, fieldsToMask []string, typesToMask []reflect.Type) error {

	if len(fieldsToMask) > 0 {
		for _, f := range fieldsToMask {
			if err := ClearFieldByName(&body, f); err != nil {
				return err
			}
		}
	}

	if len(typesToMask) > 0 {
		for _, t := range typesToMask {
			if err := ClearFieldByType(&body, t); err != nil {
				return err
			}
		}
	}

	return nil
}

func isIgnoredPath(ignoredPaths []string, path string) bool {
	for _, ip := range ignoredPaths {
		if strings.Contains(path, ip) {
			return true
		}
	}
	return false
}
