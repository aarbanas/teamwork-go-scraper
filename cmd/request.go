package main

import (
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

func getTimeLogs(startDate *string, endDate *string) (*Response, error) {

	envVariables := getEnvVariables()
	if envVariables == nil {
		return nil, errors.New("Can't load environment variables")
	}

	// URL
	url := fmt.Sprintf("%s?page=1&pageSize=50&getTotals=true&userId=%s&fromDate=%s&toDate=%s&sortBy=date&sortOrder=desc&matchAllTags=true", envVariables.TeamworkUrl, envVariables.UserId, *startDate, *endDate)

	// Create a new request with GET method
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Printf("Error in NewRequest: %s", err)
		return nil, err
	}

	// Convert string to bytes
	dataBytes := []byte(envVariables.ApiKey)

	// Encode to base64
	encodedData := base64.StdEncoding.EncodeToString(dataBytes)
	token := fmt.Sprintf("Basic %s", encodedData)

	// Add headers to the request
	req.Header.Add("Authorization", token)

	// Create a new HTTP client and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error in client.Do: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	// Read the response body
	body, _ := io.ReadAll(resp.Body)

	var result Response
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
		os.Exit(1)
	}

	return &result, nil
}
