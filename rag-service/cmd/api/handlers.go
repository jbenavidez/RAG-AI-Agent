package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Config) TestEndpoint(w http.ResponseWriter, r *http.Request) {

	// test query
	searchQuery := "Interceptor Sewers at Various Locations in the Boroughs of M"
	result, err := c.WDBRepo.GetDocuments(searchQuery)
	if err != nil {
		fmt.Println("unable to get data", err)
		return
	}
	textContext := c.DocsToContext(result)
	//set test resposne
	jsonResponse := make(map[string]string)
	jsonResponse["message"] = textContext
	//set response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)

}
