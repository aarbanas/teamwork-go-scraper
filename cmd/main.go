package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	defaultStartDate, defaultEndDate := getDefaultDates()

	action := flag.String("action", "tag", "Action on which hours will be calculated (tag or projectId)")
	value := flag.String("value", "overtime", "Value for the specified action.")
	startDate := flag.String("startDate", defaultStartDate, "Must be in format YYYY-MM-DD")
	endDate := flag.String("endDate", defaultEndDate, "Must be in format YYYY-MM-DD")

	flag.Parse()

	// Validate command-line arguments
	validateInputParams(action, startDate, endDate)

	// Send request to Teamtailor
	response, responseError := getTimeLogs(startDate, endDate)
	if responseError != nil {
		os.Exit(1)
	}

	hours := 0.0
	switch *action {
	case "tag":
		context := NewContext(&CalculateByTag{})
		result := context.ExecuteStrategy(response, *value)
		hours = result
		break
	case "projectId":
		context := NewContext(&CalculateByProjectId{})
		result := context.ExecuteStrategy(response, *value)
		hours = result
		break
	}

	fmt.Printf("Total overtime hours: %.1f\n", hours)
}
