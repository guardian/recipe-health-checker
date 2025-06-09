package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/guardian/recipe-health-checker/models"
)

func WriteDoc(baseUrl *string, rec *models.Record) error {
	url := fmt.Sprintf("%s/_doc", *baseUrl)
	jsonData, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	buffer := bytes.NewReader(jsonData)
	resp, err := http.Post(url, "application/json", buffer)
	if err != nil {
		return err
	}

	maybeContent, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("WARNING unable to read full response from server: %s", err)
	}
	if resp.StatusCode != 200 {
		var outputString string
		if maybeContent == nil {
			outputString = "(no output)"
		} else {
			outputString = string(maybeContent)
		}

		return fmt.Errorf("Server responded %d: %s", resp.StatusCode, outputString)
	}
	return nil
}
