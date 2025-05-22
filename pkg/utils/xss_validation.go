package utils

import "github.com/microcosm-cc/bluemonday"

func IsXSSSafe(input string) bool {
	p := bluemonday.UGCPolicy()
	sanitized := p.Sanitize(input)

	return input == sanitized
}
