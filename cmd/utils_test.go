package main

import (
	"testing"
	"time"
)

func TestValidateAction(t *testing.T) {
	tests := []struct {
		name            string
		action          string
		isErrorExpected bool
	}{
		{"valid action tag", "tag", false},
		{"valid action projectId", "projectId", false},
		{"invalid action", "invalidAction", true},
		{"empty action", "", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateAction(&test.action)
			if (err != nil) != test.isErrorExpected {
				t.Errorf("validateAction() error = %v, isErrorExpected %v", err, test.isErrorExpected)
			}
		})
	}
}

func TestValidateDate(t *testing.T) {
	tests := []struct {
		name            string
		date            string
		isErrorExpected bool
	}{
		{"valid date format", "2024-01-15", true},
		{"invalid format", "15.01.2024", false},
		{"empty date", "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateDate(&test.date)
			if err != test.isErrorExpected {
				t.Errorf("validateDate() error = %v, isErrorExpected %v", err, test.isErrorExpected)
			}
		})
	}
}

func TestConvertStringFormatToDate(t *testing.T) {
	tests := []struct {
		name            string
		date            string
		isErrorExpected bool
	}{
		{"ValidDate", "2022-11-11", false},
		{"InvalidDate", "invalid-date", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := convertStringFormatToDate(test.date)
			if err != nil && !test.isErrorExpected {
				t.Errorf("Unexpected error for input %s: %s", test.date, err.Error())
			}
		})
	}
}

func TestGetWorkingDays(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)

	workdays := getWorkingDays(startDate, endDate)

	if len(*workdays) != 5 {
		t.Errorf("Workdays count was incorrect, got: %d, want: %d.", len(*workdays), 5)
	}

	for _, d := range *workdays {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			t.Errorf("Workdays should not include weekend, but got: %s.", d.Weekday())
		}
	}
}
