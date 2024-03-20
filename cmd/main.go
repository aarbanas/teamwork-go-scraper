package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/aarbanas/teamwork-go-scraper/config"
	"github.com/aarbanas/teamwork-go-scraper/teamwork"
)

func main() {
	defaultStartDate, defaultEndDate := teamwork.GetDefaultDates()
	absPath, _ := filepath.Abs("config.json")
	configuration := config.LoadConfiguration(absPath)

	action := flag.String("action", "tag", "Action on which hours will be calculated (tag or projectId)")
	value := flag.String("value", "overtime", "Value for the specified action.")
	startDate := flag.String("startDate", defaultStartDate, "Must be in format YYYY-MM-DD")
	endDate := flag.String("endDate", defaultEndDate, "Must be in format YYYY-MM-DD")
	startTime := flag.String("t", "09:00", "Start time when to log hours from (HH:mm)")

	log := flag.Bool("l", false, "Enter the logging mode (default: reading logged hours)")
	projectMode := flag.Bool("p", false, "If selected hours will be logged by project id (default: log by task id )")
	includeCroatianHolidays := flag.Bool("h", false, "Use for including Croatian national holidays in the calculations")
	checkMissingHours := flag.Bool("c", false, "Use for checking if there are some days where hours are not logged")
	nonBillable := flag.Bool("n", false, "Non billable hours in log mode (default: isBillable)")
	teamHours := flag.Bool("th", false, "Check time logs for team members")

	flag.Parse()

	// Validate command-line arguments
	teamwork.ValidateInputParams(*action, *startDate, *endDate)

	if *teamHours {
		teamwork.CheckTeamMembersHours(configuration, *includeCroatianHolidays, *startDate, *endDate)
		os.Exit(0)
	}

	if *log {
		teamwork.LogHours(*startDate, *endDate, *startTime, *projectMode, *includeCroatianHolidays, *nonBillable, configuration)
		os.Exit(0)
	}

	// Send request to Teamwork
	response, err := teamwork.GetTimeLogs(*startDate, *endDate, configuration)
	if err != nil {
		os.Exit(1)
	}

	if *checkMissingHours {
		if len(response.TimeEntries) == 0 {
			fmt.Println("No time entries found")
		}

		missingHours := teamwork.ValidateMissingHours(*response, *includeCroatianHolidays, *startDate, *endDate)
		if len(missingHours) > 0 {
			for _, missingHour := range missingHours {
				fmt.Printf("Missing \"%.1f\" hours for date \"%s\"\n", missingHour.HoursDecimal, missingHour.Date.Format(time.DateOnly))
			}
		} else {
			fmt.Println("All hours logged correctly!")
		}
		os.Exit(0)
	}

	hours := 0.0
	switch *action {
	case "tag":
		context := teamwork.NewContext(&teamwork.CalculateByTag{})
		result := context.ExecuteStrategy(response, *value)
		hours = result
	case "projectId":
		context := teamwork.NewContext(&teamwork.CalculateByProjectId{})
		result := context.ExecuteStrategy(response, *value)
		hours = result
	}

	fmt.Println("\n*********************")
	fmt.Printf("\nTotal hours in period from %s to %s by \"%s\" with value \"%s\" are: \"%.1f\"\n", *startDate, *endDate, *action, *value, hours)
	fmt.Println("\n*********************")

}
