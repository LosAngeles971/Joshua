package knowledge

import (
	"fmt"
	"strings"
)

// clean removes carriage returns from the given string and trims the result
func clean(source string) string {
	source = strings.Replace(source, "\n", "", -1)
	source = strings.Replace(source, "\t", "", -1)
	return strings.TrimSpace(source)
}

type Source struct {
	source string
}

func newSource(source string) *Source {
	return &Source{
		source: source,
	}
}

// readUntil get the first chars from the given string till to the stop sequence
func (s *Source) readUntil(source string, stop string) (string, error) {
	token := ""
	s.source = clean(s.source)
	for {
		if len(s.source) < len(stop) {
			return "", fmt.Errorf("string [%v] is shorter than stop char %v", source, stop)
		}
		if string(source[0:len(stop)]) == stop {
			s.source = source[len(stop):]
			return clean(token), nil
		}
		token = token + source[0:1]
		source = source[1:]
	}
}