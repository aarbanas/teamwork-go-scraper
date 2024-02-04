package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func validateAction(action string) error {
	if action != "tag" && action != "projectId" {
		return errors.New(" Action must be \"tag\" or \"projectId")
	}

	return nil
}

func validateDate(date string) bool {
	_, err := time.Parse(time.DateOnly, date)
	return err == nil
}

func validateInputParams(action string, startDate string, endDate string) {
	if err := validateAction(action); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := validateDate(startDate); !err {
		fmt.Fprintf(os.Stderr, "error: Wrong startDate format!\n")
		os.Exit(1)
	}

	if err := validateDate(endDate); !err {
		fmt.Fprintf(os.Stderr, "error: Wrong endDate format!\n")
		os.Exit(1)
	}
}

func getDefaultDates() (string, string) {
	now := time.Now().UTC()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)

	return firstOfMonth.Format(time.DateOnly), now.Format(time.DateOnly)
}

func convertStringFormatToDate(dateString string) (time.Time, error) {
	date, err := time.Parse(time.DateOnly, dateString)
	if err != nil {
		fmt.Println("Failed to parse date:", err)
		return time.Time{}, err
	}

	return date, nil
}

func convertDateToString(date time.Time) string {
	return date.Format(time.DateOnly)
}

func getWorkingDays(startDate, endDate time.Time) []time.Time {
	var workdays []time.Time
	for d := startDate; d.Before(endDate.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			workdays = append(workdays, d)
		}
	}
	return workdays
}

func removeNoneWorkingDays(workDays *[]time.Time, nonWorkingDays *CroatianNoneWorkingDays) {
	if nonWorkingDays == nil || len(*workDays) < 1 {
		return
	}

	for index, workday := range *workDays {
		for _, nonWorkingDay := range *nonWorkingDays {
			if workday.Format(time.DateOnly) == nonWorkingDay.Date {
				*workDays = append((*workDays)[:index], (*workDays)[index+1:]...)
			}
		}
	}
}

func prepareWorkDays(startDate, endDate string, includeCroHolidays bool) (*[]time.Time, error) {
	workingStartDate, err := convertStringFormatToDate(startDate)
	if err != nil {
		return nil, err
	}

	workingEndDate, err := convertStringFormatToDate(endDate)
	if err != nil {
		return nil, err
	}

	workDays := getWorkingDays(workingStartDate, workingEndDate)

	if includeCroHolidays {
		croNoWorkingDays, err := getCroatianNoneWorkingDays(workingStartDate.Year())
		if err != nil {
			return nil, err
		}
		removeNoneWorkingDays(&workDays, croNoWorkingDays)
	}

	if len(workDays) < 1 {
		return nil, errors.New("There are no workdays in selected period")
	}

	return &workDays, nil
}

func areDatesTheSame(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()

	return y1 == y2 && m1 == m2 && d1 == d2
}
