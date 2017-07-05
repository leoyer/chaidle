package types

import (
	"encoding/xml"
	"io"

)

type SaveRequest struct{
	XMLName xml.Name `xml:"REQUEST"`
	Id string  `xml:"id,attr"`
	Comment string  `xml:"comment,attr"`
	Version string `xml:"VERSION"`
	Table string  `xml:"TABLE"`
	ROWKEY string `xml:"ROWKEY"`
	Column string `xml:"COLUMN"`
	Qualifiers  StringMap `xml:"QUALIFIERS"`
}


// StringMap is a map[string]string.
type StringMap map[string]string

// StringMap marshals into XML.
func (s StringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {

	tokens := []xml.Token{start}

	for key, value := range s {
		t := xml.StartElement{Name: xml.Name{"", key}}
		tokens = append(tokens, t, xml.CharData(value), xml.EndElement{t.Name})
	}

	tokens = append(tokens, xml.EndElement{start.Name})

	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	// flush to ensure tokens are written
	err := e.Flush()
	if err != nil {
		return err
	}


	return nil
}
type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}
func (m *StringMap) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*m = StringMap{}
	for {
		var e xmlMapEntry

		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		(*m)[e.XMLName.Local] = e.Value
	}
	return nil
}
