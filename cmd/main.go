package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	defaultStartDate, defaultEndDate := getDefaultDates()
	absPath, _ := filepath.Abs("config.json")
	configuration := loadConfiguration(absPath)

	action := flag.String("action", "tag", "Action on which hours will be calculated (tag or projectId)")
	value := flag.String("value", "overtime", "Value for the specified action.")
	startDate := flag.String("startDate", defaultStartDate, "Must be in format YYYY-MM-DD")
	endDate := flag.String("endDate", defaultEndDate, "Must be in format YYYY-MM-DD")

	log := flag.Bool("l", false, "Enter the logging mode (default: reading logged hours)")
	projectMode := flag.Bool("p", false, "If selected hours will be logged by project id (default: log by task id )")

	flag.Parse()

	// Validate command-line arguments
	validateInputParams(*action, *startDate, *endDate)

	if *log {
		logHours(*startDate, *endDate, *projectMode, configuration)
		os.Exit(1)
	}

	// Send request to Teamwork
	response, responseError := getTimeLogs(*startDate, *endDate, configuration)
	if responseError != nil {
		os.Exit(1)
	}

	hours := 0.0
	switch *action {
	case "tag":
		context := NewContext(&CalculateByTag{})
		result := context.ExecuteStrategy(response, *value)
		hours = result
	case "projectId":
		context := NewContext(&CalculateByProjectId{})
		result := context.ExecuteStrategy(response, *value)
		hours = result
	}

	fmt.Println("\n*********************")
	fmt.Printf("\nTotal hours in period from %s to %s by \"%s\" with value \"%s\" are: \"%.1f\"\n", *startDate, *endDate, *action, *value, hours)
	fmt.Println("\n*********************")

}
