package simpleDocxParser

import "encoding/xml"

type xmlStyles struct {
	XMLName xml.Name `xml:"styles"`
	XlatentStyles xmlLatentStyles `xml:"latentStyles"`
	Xstyles []xmlStyle `xml:"style"`
}

type xmlLatentStyles struct {
	XMLName xml.Name `xml:"latentStyles"`
	XLsdException []xmlLsdException `xml:"lsdException"`
}

type xmlLsdException struct {
	XMLName xml.Name `xml:"lsdException"`
	Xname string `xml:"name,attr"`
	/*
		XuiPriority int `xml:"uiPriority,attr"`
		XsemiHidden int `xml:"semiHidden,attr"`
		XunhideWhenUsed int `xml:"unhideWhenUsed,attr"`
		XqFormat int `xml:"qFormat,attr"`
	*/
}

type xmlStyle struct {
	XMLName xml.Name `xml:"style"`
	Xtype string `xml:"type,attr"`
	XstyleId string `xml:"styleId,attr"`
	Xname xmlName `xml:"name"`
	XbasedOn xmlBasedOn `xml:"basedOn"`
}
type xmlName struct {
	XMLName xml.Name `xml:"name"`
	Xval string `xml:"val,attr"`
}

type xmlBasedOn struct {
	XMLName xml.Name `xml:"basedOn"`
	Xval string `xml:"val,attr"`
}

