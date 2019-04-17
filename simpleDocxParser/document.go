package simpleDocxParser

import "encoding/xml"

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
}

type xmlHyperlink struct {
	Xr xmlr `xml:""`
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
	Xbefore int `xml:"before"`
	Xafter int `xml:"after"`
}

type xmlr struct {
	XrPr xmlrPr `xml:"rPr"`
	Xt string `xml:"t"`
}

type xmlrPr struct {
	XrStyle xmlrStyle `xml:"rStyle"`
}

type xmlrStyle struct {
	Xval string `xml:"val"`
}