package simpleDocxParser

import "encoding/xml"

type xmlRelationships struct {
	XMLName xml.Name `xml:"Relationships"`
	Xrelationships []xmlRelationship `xml:"Relationship"`
}

type xmlRelationship struct {
	XMLName xml.Name `xml:"Relationship"`
	Xid string `xml:"Id,attr"`
	Xtype string `xml:"Type,attr"`
	Xtarget string `xml:"Target,attr"`
}