package templates

import (
	"html/template"
	"strings"
)

func Indent(s template.HTML, prefix string) template.HTML {
	lines := strings.Split(string(s), "\n")
	for i, line := range lines {
		lines[i] = prefix + line
	}
	return template.HTML(strings.Join(lines, "\n"))
}
