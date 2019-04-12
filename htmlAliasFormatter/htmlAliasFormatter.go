package htmlAliasFormatter

import (
	"fmt"
	"strings"
)

type HTMLElementReplacer struct {
	htmlElementAliases map[string]string
}

// Factory Functions

func New(userHtmlElementAliases map[string]string) HTMLElementReplacer {
	var htmlElementAliases map[string]string
	if userHtmlElementAliases != nil {
		htmlElementAliases = userHtmlElementAliases
	} else {
		htmlElementAliases = make(map[string]string)
	}
	return HTMLElementReplacer{ htmlElementAliases }
}

// Worker Functions
func (hep *HTMLElementReplacer) AddAlias(alias string, htmlTag string) {
	hep.htmlElementAliases[alias] = htmlTag
}

func (hep *HTMLElementReplacer) DelAlias(alias string) {
	delete(hep.htmlElementAliases, alias)
}

func (hep *HTMLElementReplacer) ConvertToHtml(content string, alias string) string {
	htmlAlias, ok := hep.htmlElementAliases[alias]
	if len(strings.Replace(content, " ", "", -1)) > 0 {
		if ok {
			return fmt.Sprintf(htmlAlias, content)
		} else {
			return fmt.Sprintf("<p>%s</p>", content)
		}
	} else {
		return "<br />"
	}

}