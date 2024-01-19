package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	YYYYMMDD = "2006-01-02"
)

func validateAction(action *string) error {
	if *action != "tag" && *action != "projectId" {
		return errors.New(" Action must be \"tag\" or \"projectId")
	}

	return nil
}

func validateDate(date *string) bool {
	_, err := time.Parse(YYYYMMDD, *date)
	return err == nil
}

func validateInputParams(action *string, startDate *string, endDate *string) {
	if err := validateAction(action); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := validateDate(startDate); err != true {
		fmt.Fprintf(os.Stderr, "error: Wrong startDate format!\n")
		os.Exit(1)
	}

	if err := validateDate(endDate); err != true {
		fmt.Fprintf(os.Stderr, "error: Wrong endDate format!\n")
		os.Exit(1)
	}
}

func getDefaultDates() (string, string) {
	now := time.Now().UTC()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)

	return firstOfMonth.Format(YYYYMMDD), now.Format(YYYYMMDD)
}
