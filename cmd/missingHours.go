package main

import (
	"fmt"
	"os"
	"time"
)

type MissingHours struct {
	Date         time.Time
	HoursDecimal float64
}

var (
	minimumWorkingHours = 7.5
)

func ValidateMissingHours(loggedHours Response, includeCroHolidays bool, startDate, endDate string) []MissingHours {
	workDays, err := prepareWorkDays(startDate, endDate, includeCroHolidays)
	if err != nil {
		fmt.Printf("Error: %s", err)
		os.Exit(1)
	}

	totalMissingHours := []MissingHours{}
	areHoursLogged := false

	for _, workDay := range *workDays {
		for _, loggedHour := range loggedHours.TimeEntries {
			areHoursLogged = false
			y1, m1, d1 := workDay.Date()
			y2, m2, d2 := loggedHour.Date.Date()

			if y1 == y2 && m1 == m2 && d1 == d2 {
				areHoursLogged = true
				if loggedHour.HoursDecimal < minimumWorkingHours {
					hoursDif := minimumWorkingHours - loggedHour.HoursDecimal
					totalMissingHours = append(totalMissingHours, MissingHours{Date: workDay, HoursDecimal: hoursDif})
				}
			}
		}

		if !areHoursLogged {
			totalMissingHours = append(totalMissingHours, MissingHours{Date: workDay, HoursDecimal: minimumWorkingHours})
		}
	}

	return totalMissingHours
}
