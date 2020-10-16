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
	http.HandleFunc("/kpis", kpis)
	http.HandleFunc("/counters", counters)

	log.Fatal(http.ListenAndServe(":8082", nil))
}

func kpis(w http.ResponseWriter, req *http.Request) {
	var mdc models.MDC

	siteDir, err := ioutil.ReadDir("data/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		println("data directory not found")
		return
	}
	for _, site := range siteDir {
		if site.IsDir() {
			siteFiles, _ := ioutil.ReadDir("data/" + site.Name())

			for _, file := range siteFiles {
				xmlFile, err := ioutil.ReadFile("data/" + site.Name() + "/" + file.Name())
				if err != nil {
					println("Error in loading" + site.Name() + "data file...")
				}
				xml.Unmarshal(xmlFile, &mdc)
				output := calculation.Calculate(mdc)

				jsOutput, err := json.Marshal(output)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(jsOutput)
			}
		}
	}
}

func counters(w http.ResponseWriter, req *http.Request) {
	var mdc models.MDC
	xmlByte, err := ioutil.ReadFile("data.xml")
	if err != nil {
		println("Error in loading data file...")
	}
	xml.Unmarshal(xmlByte, &mdc)
	output := calculation.PopulateCountersKeyValue(mdc, "siteName")
	intOutput := make(map[string]map[string]int)

	for EUtranCellFDDKey, EUtranCellFDDValue := range output {
		counters := calculation.ConvertCounterToInt(EUtranCellFDDValue)

		intOutput[EUtranCellFDDKey] = counters
	}

	jsOutput, err := json.Marshal(intOutput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	print("counters are ok!")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsOutput)
}
