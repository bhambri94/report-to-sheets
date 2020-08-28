package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type CucumberReport []struct {
	Keyword  string        `json:"keyword"`
	Name     string        `json:"name"`
	Line     int           `json:"line"`
	ID       string        `json:"id"`
	Tags     []interface{} `json:"tags"`
	URI      string        `json:"uri"`
	Elements []struct {
		ID      string        `json:"id"`
		Keyword string        `json:"keyword"`
		Line    int           `json:"line"`
		Name    string        `json:"name"`
		Tags    []interface{} `json:"tags"`
		Type    string        `json:"type"`
		Steps   []struct {
			Arguments []interface{} `json:"arguments,omitempty"`
			Keyword   string        `json:"keyword"`
			Line      int           `json:"line,omitempty"`
			Name      string        `json:"name,omitempty"`
			Match     struct {
				Location string `json:"location"`
			} `json:"match"`
			Result struct {
				Status   string `json:"status"`
				Duration int64  `json:"duration"`
			} `json:"result"`
			Hidden bool `json:"hidden,omitempty"`
		} `json:"steps"`
	} `json:"elements"`
}

func GetReport(jsonName string) [][]interface{} {
	jsonFile, err := os.Open("uploads/latestreport.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened latestreport.json")
	byteValue, _ := ioutil.ReadAll(jsonFile)
	defer jsonFile.Close()
	var cucumberReport CucumberReport
	json.Unmarshal(byteValue, &cucumberReport)
	i := 0
	var finalValues [][]interface{}
	for i < len(cucumberReport) {
		var row []interface{}
		if cucumberReport[i].Keyword == "Feature" {
			row = append(row, jsonName[:len(jsonName)-4]+"html", cucumberReport[i].Name)
		}
		j := 0
		for j < len(cucumberReport[i].Elements) {
			if cucumberReport[i].Elements[j].Keyword == "Scenario" {
				row = append(row, cucumberReport[i].Elements[j].Name)
			}
			k := 0
			Status := "true"
			for k < len(cucumberReport[i].Elements[j].Steps) {
				if cucumberReport[i].Elements[j].Steps[k].Result.Status == "failed" {
					Status = "false"
				}
				if cucumberReport[i].Elements[j].Steps[k].Result.Status == "skipped" {
					if Status == "true" || Status == "skipped" {
						Status = "skipped"
					}
				}
				k++
			}
			if Status == "false" {
				row = append(row, "Fail")
			} else if Status == "true" {
				row = append(row, "Pass")
			} else if Status == "skipped" {
				row = append(row, "Skipped")
			}
			j++
		}
		finalValues = append(finalValues, row)
		i++
	}
	fmt.Println(finalValues)
	return finalValues
}
