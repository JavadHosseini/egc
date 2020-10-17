package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
				println("kips are ok!")
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsOutput)
			}
		}
	}
}

func counters(w http.ResponseWriter, req *http.Request) {
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
				siteName := strings.Split(strings.Split(mdc.Mfh.Sn, ",")[2], "=")[1]
				countersModelList := calculation.PopulateCountersKeyValue(mdc, siteName)
				for i, counterModel := range countersModelList {
					outputIntCounter := calculation.ConvertCounterToInt(counterModel.StringCounters)
					//counterModel.Counters = make(map[string]int)
					countersModelList[i].Counters = outputIntCounter
				}

				jsOutput, err := json.Marshal(countersModelList)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				println("counters are ok!")
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsOutput)
			}
		}
	}
}
