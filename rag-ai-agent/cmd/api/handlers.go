package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *RagConfig) TestEndpoint(w http.ResponseWriter, r *http.Request) {

	// test query
	searchQuery := "List the project names along with their budget forecast and total budget changes for each interceptor sewer project."
	// get doscs
	result, err := c.WDBRepo.GetDocuments(searchQuery)
	if err != nil {
		fmt.Println("unable to get data", err)
		return
	}
	// generate resp from slide of docs
	resp, err := c.GenerateAnswerFromSlides(r.Context(), searchQuery, result)
	if err != nil {
		fmt.Println("error generate response", err)
		return
	}
	// generate json respon
	jsonResponse := make(map[string]any)
	jsonResponse["message"] = resp
	//set response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)

}
