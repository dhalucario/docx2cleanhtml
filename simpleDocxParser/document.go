package simpleDocxParser


type xmlDocument struct {
	body xmlBody `xml:"w:body"`
}

type xmlBody struct {
	paragraph []xmlParagraph `xml:"w:p"`
}

type xmlParagraph struct {
	pPr xmlpPr `xml:"w:pPr"`
	r xmlr `xml:"r"`
}

type xmlHyperlink struct {
	r xmlr `xml:""`
}

type xmlpPr struct {
	pStyle xmlpStyle `xml:"w:pStyle"`
	spacing xmlspacing `xml:"w:spacing"`
	rPr xmlrPr `xml:"w:rPr"`
}

type xmlpStyle struct {
	val string `xml:"val"`
}

type xmlspacing struct {
	before int `xml:"w:before"`
	after int `xml:"w:after"`
}

type xmlr struct {
	rPr xmlrPr `xml:"w:rPr"`
	t string `xml:"w:t"`
}

type xmlrPr struct {
	rStyle xmlrStyle `xml:"rStyle"`
}

type xmlrStyle struct {
	val string `xml:"val"`
}