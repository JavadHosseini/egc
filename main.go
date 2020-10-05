package main

import (
	"encoding/xml"
	"io/ioutil"
	"xml-parser/calculation"
	"xml-parser/models"
)

func main() {

	var mdc models.MDC
	xmlByte, err := ioutil.ReadFile("data.xml")
	if err != nil {
		println("Error in loading data file...")
	}

	xml.Unmarshal(xmlByte, &mdc)
	calculation.Calculate(mdc)
}
