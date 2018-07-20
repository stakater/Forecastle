package handlers

import (
	"io/ioutil"
	"net/http"
)

// FileHandler func responsible for accessing files and returning their data
// Useful for getting namespace from file in case of Forecastle
func FileHandler(responseWriter http.ResponseWriter, request *http.Request) {
	if path := request.FormValue("path"); path != "" {
		fileContents, err := ioutil.ReadFile(path)

		if err != nil {
			logger.Error("File not found", err)
			http.Error(responseWriter, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = responseWriter.Write(fileContents)
		if err != nil {
			logger.Error("An error occurred while rendering contents to output: ", err)
			http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		}
	} else {
		errorString := "No file path specified"

		logger.Error(errorString)
		http.Error(responseWriter, errorString, http.StatusBadRequest)
	}

}
