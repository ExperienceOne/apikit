package xhttp

import (
	"net/http"
)

var statusCode = map[int]string{
	http.StatusContinue:           "http.StatusContinue",
	http.StatusSwitchingProtocols: "http.StatusSwitchingProtocols",
	http.StatusProcessing:         "http.StatusProcessing",

	http.StatusOK:                   "http.StatusOK",
	http.StatusCreated:              "http.StatusCreated",
	http.StatusAccepted:             "http.StatusAccepted",
	http.StatusNonAuthoritativeInfo: "http.StatusNonAuthoritativeInfo",
	http.StatusNoContent:            "http.StatusNoContent",
	http.StatusResetContent:         "http.StatusResetContent",
	http.StatusPartialContent:       "http.StatusPartialContent",
	http.StatusMultiStatus:          "http.StatusMultiStatus",
	http.StatusAlreadyReported:      "http.StatusAlreadyReported",
	http.StatusIMUsed:               "http.StatusIMUsed",

	http.StatusMultipleChoices:   "http.StatusMultipleChoices",
	http.StatusMovedPermanently:  "http.StatusMovedPermanently",
	http.StatusFound:             "http.StatusFound",
	http.StatusSeeOther:          "http.StatusSeeOther",
	http.StatusNotModified:       "http.StatusNotModified",
	http.StatusUseProxy:          "http.StatusUseProxy",
	http.StatusTemporaryRedirect: "http.StatusTemporaryRedirect",
	http.StatusPermanentRedirect: "http.StatusPermanentRedirect",

	http.StatusBadRequest:                   "http.StatusBadRequest",
	http.StatusUnauthorized:                 "http.StatusUnauthorized",
	http.StatusPaymentRequired:              "http.StatusPaymentRequired",
	http.StatusForbidden:                    "http.StatusForbidden",
	http.StatusNotFound:                     "http.StatusNotFound",
	http.StatusMethodNotAllowed:             "http.StatusMethodNotAllowed",
	http.StatusNotAcceptable:                "http.StatusNotAcceptable",
	http.StatusProxyAuthRequired:            "http.StatusProxyAuthRequired",
	http.StatusRequestTimeout:               "http.StatusRequestTimeout",
	http.StatusConflict:                     "http.StatusConflict",
	http.StatusGone:                         "http.StatusGone",
	http.StatusLengthRequired:               "http.StatusLengthRequired:",
	http.StatusPreconditionFailed:           "http.StatusPreconditionFailed",
	http.StatusRequestEntityTooLarge:        "http.StatusRequestEntityTooLarge",
	http.StatusRequestURITooLong:            "http.StatusRequestURITooLong",
	http.StatusUnsupportedMediaType:         "http.StatusUnsupportedMediaType",
	http.StatusRequestedRangeNotSatisfiable: "http.StatusRequestedRangeNotSatisfiable",
	http.StatusExpectationFailed:            "http.StatusExpectationFailed",
	http.StatusTeapot:                       "http.StatusTeapot",
	http.StatusUnprocessableEntity:          "http.StatusUnprocessableEntity",
	http.StatusLocked:                       "http.StatusLocked",
	http.StatusFailedDependency:             "http.StatusFailedDependency",
	http.StatusUpgradeRequired:              "http.StatusUpgradeRequired",
	http.StatusPreconditionRequired:         "http.StatusPreconditionRequired",
	http.StatusTooManyRequests:              "http.StatusTooManyRequests",
	http.StatusRequestHeaderFieldsTooLarge:  "http.StatusRequestHeaderFieldsTooLarge",
	http.StatusUnavailableForLegalReasons:   "http.StatusUnavailableForLegalReasons",

	http.StatusInternalServerError:           "http.StatusInternalServerError",
	http.StatusNotImplemented:                "http.StatusNotImplemented",
	http.StatusBadGateway:                    "http.StatusBadGateway",
	http.StatusServiceUnavailable:            "http.StatusServiceUnavailable",
	http.StatusGatewayTimeout:                "http.StatusGatewayTimeout",
	http.StatusHTTPVersionNotSupported:       "http.StatusHTTPVersionNotSupported",
	http.StatusVariantAlsoNegotiates:         "http.StatusVariantAlsoNegotiates",
	http.StatusInsufficientStorage:           "http.StatusInsufficientStorage",
	http.StatusLoopDetected:                  "http.StatusLoopDetected",
	http.StatusNotExtended:                   "http.StatusNotExtended",
	http.StatusNetworkAuthenticationRequired: "http.StatusNetworkAuthenticationRequired",
}

// resolveStatusCode returns a var for the HTTP status code
func ResolveStatusCode(code int) string {
	return statusCode[code]
}
