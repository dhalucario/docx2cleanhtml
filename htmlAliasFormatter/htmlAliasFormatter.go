package htmlAliasFormatter

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
	return HTMLElementReplacer{htmlElementAliases}
}

// Worker Functions
func (hep *HTMLElementReplacer) AddAlias(alias string, htmlTag string) {
	hep.htmlElementAliases[alias] = htmlTag
}

func (hep *HTMLElementReplacer) DelAlias(alias string) {
	delete(hep.htmlElementAliases, alias)
}