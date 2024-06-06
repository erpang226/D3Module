package tools

import (
	"bytes"
	"strings"
	"unicode"
)

func ToSnakeCase(str string) string {
	var buf bytes.Buffer
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				buf.WriteRune('_')
			}
			buf.WriteString(strings.ToLower(string(r)))
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}
