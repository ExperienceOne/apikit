package xhttp

import "strings"

func ExtractContentType(header string) string {
	if header == "" {
		return ""
	}
	i := strings.Index(header, ";")
	if i == -1 {
		i = len(header)
	}
	return strings.TrimSpace(strings.ToLower(header[:i]))
}

func ContentTypeInList(types []string, typ string) bool {

	for _, t := range types {
		if t == typ {
			return true
		}
	}
	return false
}
