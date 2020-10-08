package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"xml-parser/calculation"
	"xml-parser/models"
)

func main() {
	http.HandleFunc("/parse", parser)

	log.Fatal(http.ListenAndServe(":8082", nil))
}

func parser(w http.ResponseWriter, req *http.Request) {
	var mdc models.MDC
	xmlByte, err := ioutil.ReadFile("data.xml")
	if err != nil {
		println("Error in loading data file...")
	}
	xml.Unmarshal(xmlByte, &mdc)
	output := calculation.Calculate(mdc)

	jsOutput, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsOutput)
}
