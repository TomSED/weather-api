package tagging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}

type FactSheetResponse struct {
	Data   FactSheetData   `json:"data"`
	Errors []ErrorResponse `json:"errors,omitempty"`
}

type FactSheetData struct {
	FactSheet FactSheet `json:"factSheet"`
}

type FactSheet struct {
	ID                 string `json:"id"`
	DisplayName        string `json:"displayName"`
	Description        string `json:"description"`
	Name               string `json:"name,omitempty"`
	Alias              string `json:"alias,omitempty"`
	AST                string `json:"AST,omitempty"`
	Partner            string `json:"Partner,omitempty"`
	CT                 string `json:"CT,omitempty"`
	CTO                string `json:"CTO,omitempty"`
	DataClassification string `json:"DataClassification,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func authLeanix(config Config) (string, error) {
	authData := `grant_type=client_credentials`
	req, err := http.NewRequest("POST", config.LeanIXAuthURL, bytes.NewBuffer([]byte(authData)))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth("apitoken", config.LeanIXAPIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // You can log the error if needed
		return "", fmt.Errorf("LeanIX: authorization error %s", body)
	}

	var authRespData AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authRespData); err != nil {
		return "", err
	}

	return ("Bearer " + authRespData.AccessToken), nil
}

func getFactsheetByID(config Config, leanixID string) (FactSheet, error) {
	query := `
        query ($leanIX_id: ID!) {
            factSheet(id: $leanIX_id) {
                id
                displayName
                description
                ... on ITComponent {
                    name
                    alias
                    AST
                    Partner
                    CT
                    CTO
                    DataClassification
                }
            }
        }`

	authHeader, err := authLeanix(config)
	if err != nil {
		fmt.Println("Failed to authenticate with Leanix")
		return FactSheet{}, err
	}

	requestData := map[string]interface{}{
		"query":     query,
		"variables": map[string]string{"leanIX_id": leanixID},
	}

	requestBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshalling request JSON:", err)
		return FactSheet{}, err
	}

	req, err := http.NewRequest("POST", config.LeanIXRequestURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return FactSheet{}, err
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("LeanIX get_id error:", err)
		return FactSheet{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("LeanIX get_id error:", resp.Status)
		return FactSheet{}, fmt.Errorf("failed to get id from LeanIX: %s", resp.Status)
	}

	var factSheetResp FactSheetResponse
	if err := json.NewDecoder(resp.Body).Decode(&factSheetResp); err != nil {
		fmt.Println("Error decoding response:", err)
		return FactSheet{}, err
	}

	if len(factSheetResp.Errors) > 0 {
		fmt.Printf("No factsheet found for id: %s. %v\n", leanixID, factSheetResp.Errors)
		return FactSheet{}, fmt.Errorf("no factsheet returned")
	}

	return factSheetResp.Data.FactSheet, nil
}
