package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Tag struct {
	Color     string
	Id        int
	Name      string
	ProjectId int
}

type Response struct {
	TimeEntries []struct {
		HoursDecimal float64   `json:"hoursDecimal"`
		ProjectID    int       `json:"projectId"`
		Tags         []Tag     `json:"tags"`
		Date         time.Time `json:"date"`
	} `json:"timeEntries"`
}

type LogResponse struct {
	Status    string `json:"STATUS"`
	TimeLogId string `json:"timeLogId"`
}

type CroatianNoneWorkingDays []struct {
	Date string `json:"date"`
}

func prepareAuthHeader(token string) string {
	// Convert string to bytes
	dataBytes := []byte(token)

	// Encode to base64
	encodedData := base64.StdEncoding.EncodeToString(dataBytes)
	header := fmt.Sprintf("Basic %s", encodedData)

	return header
}

func handler(url string, requestMethod string, apiKey string, requestBody interface{}) (*[]byte, error) {
	var req *http.Request
	var err error

	switch body := requestBody.(type) {
	// Create a new request depending on the request body
	case *bytes.Buffer:
		req, err = http.NewRequest(requestMethod, url, body)
	case nil:
		req, err = http.NewRequest(requestMethod, url, nil)
	default:
		err = fmt.Errorf("body type not supported")
	}

	if err != nil {
		return nil, err
	}

	if apiKey != "" {
		token := prepareAuthHeader(apiKey)

		// Add headers to the request
		req.Header.Add("Authorization", token)
		req.Header.Add("Content-Type", "application/json")
	}

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &responseBody, nil
}

func getTimeLogs(startDate string, endDate string, configuration Config) (*Response, error) {

	// URL
	url := fmt.Sprintf("%s/projects/api/v2/time.json?page=1&pageSize=50&getTotals=true&userId=%s&fromDate=%s&toDate=%s&sortBy=date&sortOrder=desc&matchAllTags=true", configuration.Url, configuration.UserId, startDate, endDate)

	responseBody, err := handler(url, "GET", configuration.ApiKey, nil)
	if err != nil {
		fmt.Printf("Error in request handler: %s", err)
	}

	var result Response
	if err := json.Unmarshal(*responseBody, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
		os.Exit(1)
	}

	return &result, nil
}

func postTimeLogs(timeLog TimeLog, projectMode bool, configuration Config) (bool, error) {
	urlReference := "tasks"
	if projectMode {
		urlReference = "projects"
	}

	// URL
	url := fmt.Sprintf("%s/%s/%s/time_entries.json", configuration.Url, urlReference, timeLog.logTimeMetaData.taskId)
	// JSON body
	data := struct {
		TimeLog struct {
			Hours       int16  `json:"hours"`
			Minutes     int16  `json:"minutes"`
			UserID      int    `json:"person-id"`
			Date        string `json:"date"`
			Time        string `json:"time"`
			Description string `json:"description"`
			IsBillable  bool   `json:"isBillable"`
			Tag         string `json:"tags,omitempty"`
		} `json:"time-entry"`
	}{
		TimeLog: struct {
			Hours       int16  `json:"hours"`
			Minutes     int16  `json:"minutes"`
			UserID      int    `json:"person-id"`
			Date        string `json:"date"`
			Time        string `json:"time"`
			Description string `json:"description"`
			IsBillable  bool   `json:"isBillable"`
			Tag         string `json:"tags,omitempty"`
		}{
			Hours:       timeLog.logTimeMetaData.hours,
			Minutes:     timeLog.logTimeMetaData.minutes,
			UserID:      timeLog.userId,
			Date:        timeLog.date,
			Time:        timeLog.time,
			Description: timeLog.logTimeMetaData.description,
			IsBillable:  timeLog.isBillable,
			Tag:         timeLog.logTimeMetaData.tag,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

	res, err := handler(url, "POST", configuration.ApiKey, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error in request handler: %s", err)
		return false, err
	}

	var result LogResponse
	if err := json.Unmarshal(*res, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
		os.Exit(1)
	}
	if result.TimeLogId == "0" {
		return false, errors.New("Did not log time for selected date")
	}

	return true, nil
}

func getCroatianNoneWorkingDays(year int) (*CroatianNoneWorkingDays, error) {
	url := fmt.Sprintf("https://date.nager.at/api/v3/publicholidays/%d/HR", year)

	responseBody, err := handler(url, "GET", "", nil)
	if err != nil {
		fmt.Printf("Error in request handler: %s", err)
	}

	var result CroatianNoneWorkingDays
	if err := json.Unmarshal(*responseBody, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
		os.Exit(1)
	}

	return &result, nil
}
