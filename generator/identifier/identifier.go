package identifier

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/exp/utf8string"
)

var (
	underscoreHyphenLiteralChars = regexp.MustCompile("([^a-zA-Z0-9])")
	assignIdentifierLiteralChars = regexp.MustCompile(`([^\w(),".*&\[\]!]*)`)
	invalidOperationIDChars      = regexp.MustCompile(`(\s|-)`)
)

func MakeIdentifier(name string) string {

	if len(underscoreHyphenLiteralChars.FindAllStringSubmatch(name, -1)) > 0 {
		partsOfName := underscoreHyphenLiteralChars.Split(name, -1)
		name = ""
		first := true
		for _, part := range partsOfName {
			if first {
				name += strings.TrimPrefix(part, "_")
				first = false
			} else {
				name += strings.Title(part)
			}
		}
	}

	ident := assignIdentifierLiteralChars.ReplaceAllString(name, "")
	return ident
}

func ValidateAndCleanOperationsID(s string) (string, error) {

	if s == "" {
		return "", errors.New("operations ID is empty")
	}

	if strings.HasPrefix(s, "_") {
		return "", fmt.Errorf("operations ID cannot start with _ (operation ID: %s)", s)
	}

	if invalidOperationIDChars.MatchString(s) {
		return "", fmt.Errorf("operations ID isn't a camel case string (operation ID: %s)", s)
	}

	utf8s := utf8string.NewString(s)
	if !unicode.IsUpper(utf8s.At(0)) {
		s = string(unicode.ToUpper(utf8s.At(0))) + utf8s.Slice(1, utf8s.RuneCount())
	}

	s = assignIdentifierLiteralChars.ReplaceAllString(s, "")
	return s, nil
}
