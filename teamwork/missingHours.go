package teamwork

import (
	"fmt"
	"os"
	"time"
)

type DayHours struct {
	Date         time.Time
	HoursDecimal float64
}

var (
	minimumWorkingHours = 7.5
)

func getDayHours(loggedHours Response) []DayHours {
	totalDayHours := []DayHours{}
	isDateExisting := false
	for _, loggedHour := range loggedHours.TimeEntries {
		if len(totalDayHours) == 0 {
			totalDayHours = append(totalDayHours, DayHours{Date: loggedHour.Date, HoursDecimal: loggedHour.HoursDecimal})
			continue
		}

		isDateExisting = false
		for index, totalDayHour := range totalDayHours {
			if areDatesTheSame(loggedHour.Date, totalDayHour.Date) {
				isDateExisting = true
				totalDayHours[index].HoursDecimal += loggedHour.HoursDecimal
				break
			}
		}

		if !isDateExisting {
			totalDayHours = append(totalDayHours, DayHours{Date: loggedHour.Date, HoursDecimal: loggedHour.HoursDecimal})
		}
	}

	return totalDayHours
}

func ValidateMissingHours(loggedHours Response, includeCroHolidays bool, startDate, endDate string) []DayHours {
	workDays, err := prepareWorkDays(startDate, endDate, includeCroHolidays)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	loggedHoursPerDay := getDayHours(loggedHours)

	totalMissingHours := []DayHours{}
	areHoursLogged := false

	for _, workDay := range *workDays {

		for _, loggedHour := range loggedHoursPerDay {
			areHoursLogged = false

			if areDatesTheSame(workDay, loggedHour.Date) {
				areHoursLogged = true
				if loggedHour.HoursDecimal < minimumWorkingHours {
					hoursDif := minimumWorkingHours - loggedHour.HoursDecimal
					totalMissingHours = append(totalMissingHours, DayHours{Date: workDay, HoursDecimal: hoursDif})
				}
				break
			}
		}

		if !areHoursLogged {
			totalMissingHours = append(totalMissingHours, DayHours{Date: workDay, HoursDecimal: minimumWorkingHours})
		}
	}

	return totalMissingHours
}
