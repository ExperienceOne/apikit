package generator

import (
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

func walkResponses(respones map[int]spec.Response, f func(statusCode int, response spec.Response)) {

	statusCodesOrder := make([]int, 0, len(respones))
	for statusCode := range respones {
		statusCodesOrder = append(statusCodesOrder, statusCode)
	}
	sort.Ints(statusCodesOrder)

	for _, statusCode := range statusCodesOrder {
		f(statusCode, respones[statusCode])
	}
}

func hasFileEndpointValidProduce(respones map[int]spec.Response, operation *Operation) (bool, int) {

	var match bool
	var counter int

	walkResponses(respones, func(statusCode int, response spec.Response) {
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
