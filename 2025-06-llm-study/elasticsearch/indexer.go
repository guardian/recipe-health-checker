package elasticsearch

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/guardian/recipe-health-checker/models"
	"io"
	"log"
	"net/http"
)

func WriteDoc(baseUrl *string, indexName *string, rec *models.Record) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	urlStr := fmt.Sprintf("%s/%s/_doc", *baseUrl, *indexName)
	//esUrl, _ := url.Parse(urlStr)
	jsonData, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	buffer := bytes.NewReader(jsonData)
	resp, err := client.Post(urlStr, "application/json", buffer)
	if err != nil {
		return err
	}

	maybeContent, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("WARNING unable to read full response from server: %s", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 201 {
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
