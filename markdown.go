package hydrawebcomponents

import (
	"bytes"
	"html/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.Table),
)

// RenderMarkdown converts markdown source to safe HTML.
func RenderMarkdown(src []byte) (template.HTML, error) {
	var buf bytes.Buffer
	if err := md.Convert(src, &buf); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}
