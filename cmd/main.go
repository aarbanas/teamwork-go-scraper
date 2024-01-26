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

	log := flag.Bool("l", false, "Enter the logging mode (default: reading logged hours)")
	projectMode := flag.Bool("p", false, "If selected hours will be logged by project id (default: log by task id )")

	flag.Parse()

	// Validate command-line arguments
	validateInputParams(action, startDate, endDate)

	// Send request to Teamtailor
	response, responseError := getTimeLogs(startDate, endDate)
	if responseError != nil {
		os.Exit(1)
	}

	hours := 0.0
	if *log == true {
		logHours(startDate, endDate, projectMode)
		os.Exit(1)
	}

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
	default:
		fmt.Println("Action not found")
		break
	}

	fmt.Println("\n*********************")
	fmt.Printf("\nTotal hours by \"%s\" with value \"%s\" are: \"%.1f\"\n", *action, *value, hours)
	fmt.Println("\n*********************")

}
