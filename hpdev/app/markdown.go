package app

import (
	"bytes"
	"strings"

	"github.com/PuerkitoBio/goquery"
	bf "github.com/russross/blackfriday"
	"github.com/sourcegraph/syntaxhighlight"
)

// Markdown is markdown processor
type Markdown struct{}

// Parse -
func (m *Markdown) ConvertHTML(md []byte) []byte {
	return bf.Run(md)
}

// SyntaxHighlight replace code literal.
// https://zupzup.org/go-markdown-syntax-highlight/
func SyntaxHighlight(htm []byte) ([]byte, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(htm))
	if err != nil {
		return nil, err
	}

	var parseErrs []error
	doc.Find("code[class*=\"language-\"]").Each(func(i int, s *goquery.Selection) {
		oldCode := s.Text()
		formatted, err := syntaxhighlight.AsHTML([]byte(oldCode))
		if err != nil {
			parseErrs = append(parseErrs, err)
			return
		}
		s.SetHtml(string(formatted))
	})

	// replace unnecessarily added html tags
	new, err := doc.Html()
	if err != nil {
		return nil, err
	}
	new = strings.Replace(new, "<html><head></head><body>", "", 1)
	new = strings.Replace(new, "</body></html>", "", 1)
	return []byte(new), nil
}
