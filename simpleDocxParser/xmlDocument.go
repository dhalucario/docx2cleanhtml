package simpleDocxParser

import (
	"encoding/xml"
	"unicode/utf8"
)

type xmlDocument struct {
	XMLName xml.Name `xml:"document"`
	Xbody   xmlBody  `xml:"body"`
}

type xmlBody struct {
	XMLName     xml.Name       `xml:"body"`
	Xparagraphs []xmlParagraph `xml:"p"`
}

// From here on we need to take over because otherwise we discard the order of the elements.
type xmlParagraph struct {
	paragraph Fparagraph
}

// Custom Style Sorted Classes
type Fparagraph struct {
	style string
	run   []Frun
}

type Frun struct {
	text  string
	urlId string
}

func (ph *xmlParagraph) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {

	var token interface{}
	var err error
	done := false

	ph.paragraph.run = make([]Frun, 1)

	for !done {
		token, err = d.Token()
		if err != nil {
			if err.Error() == "EOF" {
				done = true
			}
		} else {
			switch token.(type) {
			case xml.StartElement:
				switch (token.(xml.StartElement)).Name.Local {
				case "t":
					ph.paragraph.UnmarshalText(d)
				case "pStyle":
					// This should only happen once per paragraph/hyperlink
					ph.paragraph.UnmarshalStyle((token.(xml.StartElement)).Attr)
				case "hyperlink":
					ph.paragraph.UnmarshalHyperlink((token.(xml.StartElement)).Attr)
				}
			}
		}

	}

	ph.paragraph.TrimLast()

	return nil

}

func (fpg *Fparagraph) UnmarshalText(d *xml.Decoder) {
	token, _ := d.Token()
	out := ""

	// We just want to check for chardata but golang doesn't like typechecks in if branches
	switch token.(type) {
	case xml.CharData:
		out = CharData2string(token.(xml.CharData))
	}

	fpg.run[len(fpg.run)-1].text = out
	fpg.run = append(fpg.run, Frun{
		text:  "",
		urlId: "",
	})
}

func (fpg *Fparagraph) UnmarshalStyle(attrs []xml.Attr) {
	for _, attr := range attrs {
		if attr.Name.Local == "val" {
			fpg.style = attr.Value
		}
	}
}

func (fpg *Fparagraph) UnmarshalHyperlink(attrs []xml.Attr) {
	for _, attr := range attrs {
		if attr.Name.Local == "id" {
			fpg.run[len(fpg.run)-1].urlId = attr.Value
		}
	}
}

func (fpg *Fparagraph) TrimLast() {
	fpg.run = fpg.run[0 : len(fpg.run)-1]
}

func CharData2string(data xml.CharData) string {
	stringBuffer := ""
	for i := 0; i < len(data); {
		singleRune, width := utf8.DecodeRune(data[i:len(data)])

		if singleRune != utf8.RuneError {
			stringBuffer = stringBuffer + string(singleRune)
			i = i + width
		}
	}
	return stringBuffer
}
