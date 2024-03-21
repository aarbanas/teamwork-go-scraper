package teamwork

import (
	"fmt"
	"strings"
	"time"

	"github.com/aarbanas/teamwork-go-scraper/config"
)

type UserMissingDates struct {
	FirstName    string
	LastName     string
	MissingDates []string
}

func (u UserMissingDates) String() string {
	return fmt.Sprintf("\033[1;31m[%s %s]\033[0m: %v", u.FirstName, u.LastName, u.MissingDates)
}

func CheckTeamMembersHours(configuration config.Config, includeCroHolidays bool, startDate, endDate string) {
	var userMissingDates []UserMissingDates

	// Iterate through all user IDs
	for _, userId := range configuration.UserIds {

		// Get user logged hours
		timeLogs, err := GetTimeLogs(startDate, endDate, configuration, userId)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		if len(timeLogs.TimeEntries) == 0 {
			fmt.Printf("\033[1;31m[%s]:\033[0m No time entries found\n", userId)
			continue
		}

		missingHours := ValidateMissingHours(*timeLogs, includeCroHolidays, startDate, endDate)
		if len(missingHours) > 0 {
			userMissingDates = append(userMissingDates, UserMissingDates{FirstName: timeLogs.TimeEntries[0].UserFirstName, LastName: timeLogs.TimeEntries[0].UserLastName})
			index := len(userMissingDates) - 1
			for _, missingHour := range missingHours {
				userMissingDates[index].MissingDates = append(userMissingDates[index].MissingDates, missingHour.Date.Format(time.DateOnly))
			}
		}
	}

	if len(userMissingDates) > 0 {
		reportOutputStrings := make([]string, len(userMissingDates))
		for i, user := range userMissingDates {
			reportOutputStrings[i] = user.String()
		}

		fmt.Println(strings.Join(reportOutputStrings, ",\n"))
	}
}
