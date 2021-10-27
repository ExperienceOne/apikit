package requestid

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"os"
	"sync/atomic"
)

const RequestIdHttpHeader = "X-Request-ID"

// incremental fallback ID in case of UUID generation failure
var requestIdFallback uint64

// use hostname as prefix for the request
// Hostname can be used to detect from where the request comes and decrease the risk of collisions (very marginal)
var hostname string

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		fmt.Println(err)
	}
}

// Generate generates a new request id
// Format: hostname-UUID or hostname-int64
func Generate() string {
	var reqID string
	u2, err := uuid.NewV4()
	if err != nil {
		newFallbackId := atomic.AddUint64(&requestIdFallback, 1)
		reqID = fmt.Sprintf("%s-%d", hostname, newFallbackId)
	} else {
		reqID = fmt.Sprintf("%s-%s", hostname, u2)
	}

	return reqID
}

// requestIdCxtKey key type used to store a request ID.
type requestIdCxtKey int

// RequestIdKey is the key to store a unique request ID in a request context.
const RequestIdKey requestIdCxtKey = 0

// Get returns a request ID for the given context; otherwise an empty string
// will be returned if the context is nil or the request ID can not be found.
func Get(ctx context.Context) string {

	if ctx == nil {
		return ""
	}

	if reqID, ok := ctx.Value(RequestIdKey).(string); ok {
		return reqID
	}

	return ""
}
