package utils

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func ToJson(message any, w http.ResponseWriter, r *http.Request) {

	response, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Failed to convert to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func AppHTTPClient(req *http.Request) ([]byte, error) {
	// Send the request
	client := &http.Client{
		Timeout: time.Duration(15) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("error reading response body: %v", err)
	}
	return body, nil
}
