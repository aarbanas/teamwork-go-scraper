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
		AvatarURL           string    `json:"avatarUrl"`
		CanEdit             bool      `json:"canEdit"`
		CompanyID           int       `json:"companyId"`
		CompanyName         string    `json:"companyName"`
		CreatedAt           time.Time `json:"createdAt"`
		Date                time.Time `json:"date"`
		DateUserPerspective time.Time `json:"dateUserPerspective"`
		Description         string    `json:"description"`
		HasStartTime        bool      `json:"hasStartTime"`
		Hours               int       `json:"hours"`
		HoursDecimal        float64   `json:"hoursDecimal"`
		ID                  int       `json:"id"`
		InvoiceNo           any       `json:"invoiceNo"`
		InvoiceStatus       any       `json:"invoiceStatus"`
		IsBillable          bool      `json:"isBillable"`
		IsBilled            bool      `json:"isBilled"`
		Minutes             int       `json:"minutes"`
		ParentTaskID        int       `json:"parentTaskId"`
		ParentTaskName      any       `json:"parentTaskName"`
		ProjectID           int       `json:"projectId"`
		ProjectName         string    `json:"projectName"`
		ProjectStatus       string    `json:"projectStatus"`
		Tags                []Tag     `json:"tags"`
		TaskEstimatedTime   int       `json:"taskEstimatedTime"`
		TaskID              int       `json:"taskId"`
		TaskIsPrivate       int       `json:"taskIsPrivate"`
		TaskIsSubTask       bool      `json:"taskIsSubTask"`
		TaskName            string    `json:"taskName"`
		TaskTags            []any     `json:"taskTags"`
		TasklistID          int       `json:"tasklistId"`
		TasklistName        string    `json:"tasklistName"`
		TicketID            any       `json:"ticketId"`
		UpdatedDate         time.Time `json:"updatedDate"`
		UserDeleted         bool      `json:"userDeleted"`
		UserFirstName       string    `json:"userFirstName"`
		UserID              int       `json:"userId"`
		UserLastName        string    `json:"userLastName"`
	} `json:"timeEntries"`
}

type Environment struct {
	UserId      string
	ApiKey      string
	TeamworkUrl string
}

func getEnvVariables() *Environment {
	userId := os.Getenv("USER_ID")
	apiKey := os.Getenv("API_KEY")
	teamworkUrl := os.Getenv("URL")

	return &Environment{UserId: userId, ApiKey: apiKey, TeamworkUrl: teamworkUrl}
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
	var httpNewRequestErr error

	switch body := requestBody.(type) {
	// Create a new request depending on the request body
	case *bytes.Buffer:
		req, httpNewRequestErr = http.NewRequest(requestMethod, url, body)
	case nil:
		req, httpNewRequestErr = http.NewRequest(requestMethod, url, nil)
	default:
		httpNewRequestErr = fmt.Errorf("body type not supported")
	}

	if httpNewRequestErr != nil {
		return nil, httpNewRequestErr
	}

	token := prepareAuthHeader(apiKey)

	// Add headers to the request
	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	resp, clientErr := client.Do(req)
	if clientErr != nil {
		return nil, clientErr
	}

	defer resp.Body.Close()

	// Read the response body
	responseBody, ioReadAllErr := io.ReadAll(resp.Body)
	if ioReadAllErr != nil {
		return nil, ioReadAllErr
	}

	return &responseBody, nil
}

func getTimeLogs(startDate *string, endDate *string) (*Response, error) {

	envVariables := getEnvVariables()
	if envVariables == nil {
		return nil, errors.New("Can't load environment variables")
	}

	// URL
	url := fmt.Sprintf("%s/v2/time.json?page=1&pageSize=50&getTotals=true&userId=%s&fromDate=%s&toDate=%s&sortBy=date&sortOrder=desc&matchAllTags=true", envVariables.TeamworkUrl, envVariables.UserId, *startDate, *endDate)

	responseBody, handlerErr := handler(url, "GET", envVariables.ApiKey, nil)
	if handlerErr != nil {
		fmt.Printf("Error in request handler: %s", handlerErr)
	}

	var result Response
	if err := json.Unmarshal(*responseBody, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
		os.Exit(1)
	}

	return &result, nil
}

func postTimeLogs(timeLog *TimeLog) (bool, error) {
	envVariables := getEnvVariables()
	if envVariables == nil {
		return false, errors.New("Can't load environment variables")
	}

	// URL
	url := fmt.Sprintf("%s/v3/tasks/%s/time.json", envVariables.TeamworkUrl, timeLog.logTimeMetaData.taskId)

	// JSON body
	data := struct {
		TimeLog struct {
			Hours       int16  `json:"hours"`
			Minutes     int16  `json:"minutes"`
			UserID      int    `json:"userId"`
			Date        string `json:"date"`
			Time        string `json:"time"`
			Description string `json:"description"`
			IsBillable  bool   `json:"isBillable"`
		} `json:"timelog"`
	}{
		TimeLog: struct {
			Hours       int16  `json:"hours"`
			Minutes     int16  `json:"minutes"`
			UserID      int    `json:"userId"`
			Date        string `json:"date"`
			Time        string `json:"time"`
			Description string `json:"description"`
			IsBillable  bool   `json:"isBillable"`
		}{
			Hours:       timeLog.logTimeMetaData.hours,
			Minutes:     timeLog.logTimeMetaData.minutes,
			UserID:      timeLog.userId,
			Date:        timeLog.date,
			Time:        timeLog.time,
			Description: timeLog.logTimeMetaData.description,
			IsBillable:  timeLog.isBillable,
		},
	}

	jsonData, jsonMarshalErr := json.Marshal(data)
	if jsonMarshalErr != nil {
		return false, jsonMarshalErr
	}

	_, handlerErr := handler(url, "POST", envVariables.ApiKey, bytes.NewBuffer(jsonData))
	if handlerErr != nil {
		fmt.Printf("Error in request handler: %s", handlerErr)
		return false, handlerErr
	}

	return true, nil
}
