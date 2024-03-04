package teamwork

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
			err := validateAction(test.action)
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
			err := validateDate(test.date)
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

func TestConvertDateToStringFormat(t *testing.T) {
	date := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	stringFormatValue := "20241231"

	stringDate := convertDateToString(date)
	if stringDate != stringFormatValue {
		t.Errorf("Invalid date value, got: %s, want: %s", stringDate, stringFormatValue)
	}
}

func TestGetWorkingDays(t *testing.T) {
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)

	workdays := getWorkingDays(startDate, endDate)

	if len(workdays) != 6 {
		t.Errorf("Workdays count was incorrect, got: %d, want: %d.", len(workdays), 5)
	}

	for _, d := range workdays {
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			t.Errorf("Workdays should not include weekend, but got: %s.", d.Weekday())
		}
	}
}

func TestRemoveNonWorkingDays(t *testing.T) {
	workdays := []time.Time{
		time.Date(2023, 3, 12, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 3, 13, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 3, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC),
	}

	var nonWorkingDays CroatianNoneWorkingDays
	nonWorkingDays = append(nonWorkingDays, struct {
		Date string `json:"date"`
	}{
		Date: "2023-03-14",
	})

	// call the function
	removeNoneWorkingDays(&workdays, &nonWorkingDays)

	if len(workdays) != 3 {
		t.Errorf("Expected 3 workdays left, but got %d", len(workdays))
	}
}

func TestAreDatesTheSame(t *testing.T) {
	tests := []struct {
		date1    time.Time
		date2    time.Time
		expected bool
	}{
		{
			date1:    time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
			date2:    time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: true,
		},
		{
			date1:    time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC),
			date2:    time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			expected: false,
		},
		// Add more test cases as needed
	}

	for _, test := range tests {
		result := areDatesTheSame(test.date1, test.date2)
		if result != test.expected {
			t.Errorf("For dates %v and %v, expected %t, but got %t", test.date1, test.date2, test.expected, result)
		}
	}
}
