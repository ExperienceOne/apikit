package generator

import (
	"math"
	"regexp"
	"sort"

	"github.com/go-openapi/spec"
)

func getDefinitionName(schema spec.Schema) string {

	if !schema.Ref.GetPointer().IsEmpty() {
		ref := schema.Ref.GetPointer().DecodedTokens()
		if len(ref) == 2 && ref[0] == "definitions" {
			return ref[1]
		}
	}
	return ""
}

func walkResponses(operation *Operation, f func(statusCode int, response spec.Response)) {
	responses := operation.Responses.StatusCodeResponses

	type Priority struct {
		statusCode int
		value      int
	}

	statusCodesOrder := make([]Priority, 0, len(responses))
	for statusCode := range responses {
		// change priority of response types with file content
		p := Priority{
			statusCode: statusCode,
			value:      statusCode,
		}

		if responses[statusCode].Schema != nil && responses[statusCode].Schema.Type.Contains("file") && operation.HasProduces(ContentTypesForFiles...) {
			p.value = int(math.Inf(1))
		}

		statusCodesOrder = append(statusCodesOrder, p)
	}

	sort.Slice(statusCodesOrder, func(i, j int) bool {
		return statusCodesOrder[i].value < statusCodesOrder[j].value
	})

	for _, v := range statusCodesOrder {
		f(v.statusCode, responses[v.statusCode])
	}
}

func hasFileEndpointValidProduce(operation *Operation) (bool, int) {

	var match bool
	var counter int

	walkResponses(operation, func(statusCode int, response spec.Response) {
		if response.Schema != nil && response.Schema.Type.Contains("file") {
			match = true
			if operation.HasProduces(ContentTypesForFiles...) {
				counter++
			}
		}
	})

	return match, counter
}

var pathItem = regexp.MustCompile(`\{([^}]*)\}`)

func makeMatcherForRoute(route string) string {

	return pathItem.ReplaceAllString(route, "<$1>")
}

func regexContains(strings []string, s string) bool {

	regex := regexp.MustCompile(s)
	for _, s := range strings {
		if regex.MatchString(s) {
			return true
		}
	}
	return false
}

func filterContentTypes(consumes []string, supportedConsumes []string) []string {

	validConsumes := make([]string, 0)
outerLoop:
	for _, consume := range consumes {
		for _, supportedConsume := range supportedConsumes {
			if string(supportedConsume[0]) == "^" {
				regex := regexp.MustCompile(supportedConsume)
				if ok := regex.MatchString(consume); ok {
					validConsumes = append(validConsumes, consume)
					continue outerLoop
				}
			} else if consume == supportedConsume {
				validConsumes = append(validConsumes, supportedConsume)
				continue outerLoop
			}
		}
	}
	return validConsumes
}
