package htmlAliasFormatter

import (
	"fmt"
)

type htmlElementReplacer struct {
	htmlElementAliases map[string]string
}

// Factory Functions

func New(userHtmlElementAliases map[string]string) htmlElementReplacer {
	var htmlElementAliases map[string]string
	if userHtmlElementAliases != nil {
		htmlElementAliases = userHtmlElementAliases
	} else {
		htmlElementAliases = make(map[string]string)
	}
	return htmlElementReplacer{ htmlElementAliases }
}

// Worker Functions
func (hep *htmlElementReplacer) AddAlias(alias string, htmlTag string) {
	hep.htmlElementAliases[alias] = htmlTag
}

func (hep *htmlElementReplacer) DelAlias(alias string) {
	delete(hep.htmlElementAliases, alias)
}

func (hep *htmlElementReplacer) ConvertToHtml(content string, alias string) string {
	htmlAlias, ok := hep.htmlElementAliases[alias]
	if ok {
		return fmt.Sprintf(htmlAlias, content)
	} else {
		return fmt.Sprintf("<p>%s</p>", content)
	}
}