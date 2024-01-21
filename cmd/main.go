package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	defaultStartDate, defaultEndDate := getDefaultDates()

	actionType := flag.String("actionType", "read", "The type of action to preform (read or log)")
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
	if *actionType == "read" {
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
	} else if *actionType == "log" {
		var taskId string
		var hours int
		var minutes int
		var description string
		fmt.Println("Enter TaskId: ")
		fmt.Scan(&taskId)

		fmt.Println("Enter hours and minutes (separate with space): ")
		fmt.Scan(&hours, &minutes)

		fmt.Println("Enter description: ")
		fmt.Scan(&description)
	} else {
		fmt.Println("Wrong action type selected. Needs to be \"read\" or \"log\"")
	}

}
