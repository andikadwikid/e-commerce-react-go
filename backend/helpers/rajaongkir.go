package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"backend-commerce/structs"

)

const RAJAONGKIR_URL = "https://rajaongkir.komerce.id/api/v1"

func RajaOngkirRequest(method, path string, body []byte, contentType string) (any, error) {
	apiKey := os.Getenv("RAJAONGKIR_API_KEY")
	url := RAJAONGKIR_URL + path

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if contentType == "" {
		contentType = "application/json"
	}

	req.Header.Set("Key", apiKey)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	var result structs.RajaOngkirResponseWrapper
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("Status %d: %s", resp.StatusCode, string(respBody))
	}

	if result.Meta.Code != 200 {
		return nil, fmt.Errorf("RajaOngkir API Error: %s (Code: %d)", result.Meta.Message, result.Meta.Code)
	}

	return result.Data, nil
}
