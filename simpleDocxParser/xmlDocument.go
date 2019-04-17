package simpleDocxParser

import (
	"encoding/xml"
	"fmt"
)

type xmlDocument struct {
	XMLName xml.Name `xml:"document"`
	Xbody xmlBody `xml:"body"`
}

type xmlBody struct {
	XMLName xml.Name `xml:"body"`
	Xparagraphs []xmlParagraph `xml:"p"`
}

type xmlParagraph struct {
	XMLName xml.Name `xml:"p"`
	XpPr xmlpPr `xml:"pPr"`
	Xr []xmlr `xml:"r"`
	Xhyper xmlHyperlink `xml:"hyperlink"`
}

type xmlHyperlink struct {
	XMLName xml.Name `xml:"hyperlink"`
	Xr []xmlr `xml:"r"`
}

type xmlpPr struct {
	XpStyle xmlpStyle `xml:"pStyle"`
	Xspacing xmlspacing `xml:"spacing"`
	XrPr xmlrPr `xml:"rPr"`
}

type xmlpStyle struct {
	Xval string `xml:"val"`
}

type xmlspacing struct {
	Xbefore int `xml:"before,attr"`
	Xafter int `xml:"after,attr"`
}

type xmlr struct {
	XrPr xmlrPr `xml:"rPr"`
	Xt string `xml:"t"`
}

type xmlrPr struct {
	XrStyle xmlrStyle `xml:"rStyle"`
}

type xmlrStyle struct {
	Xval string `xml:"val,attr"`
}

// Custom Marshal to keep order.

func (hl *xmlHyperlink/*p *xmlParagraph*/) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	fmt.Println(start)

	return nil

}