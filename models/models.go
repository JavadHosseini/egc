package models

import (
	"encoding/xml"
)

type MFH struct {
	XMLName xml.Name `xml:"mfh"`
	Ffv     string   `xml:"ffv"`
	Sn      string   `xml:"sn"`
	St      string   `xml:"st"`
	Vn      string   `xml:"vn"`
	Cbt     string   `xml:"cbt"`
}

type NEID struct {
	XMLName xml.Name `xml:"neid"`
	Neun    string   `xml:"neun"`
	Nedn    string   `xml:"nedn"`
	Nesw    string   `xml:"nesw"`
}

type MV struct {
	XMLName xml.Name `xml:"mv"`
	Moid    string   `xml:"moid"`
	R       []string `xml:"r"`
}

type MI struct {
	XMLName xml.Name `xml:"mi"`
	Mts     string   `xml:"mts"`
	Gp      string   `xml:"gp"`
	Mt      []string `xml:"mt"`
	Mv      []MV     `xml:"mv"`
}

type MD struct {
	XMLName xml.Name `xml:"md"`
	Neid    NEID     `xml:"neid"`
	Mi      MI       `xml:"mi"`
}

type MDC struct {
	XMLName xml.Name `xml:"mdc"`
	Mfh     MFH      `xml:"mfh"`
	Md      []MD     `xml:"md"`
}
