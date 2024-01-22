package main

import "testing"

func TestValidateAction(t *testing.T) {
	tests := []struct {
		name    string
		action  string
		wantErr bool
	}{
		{"valid action tag", "tag", false},
		{"valid action projectId", "projectId", false},
		{"invalid action", "invalidAction", true},
		{"empty action", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAction(&tt.action)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func validateDate(date *string) bool {
// 	_, err := time.Parse(YYYYMMDD, *date)
// 	return err == nil
// }

func TestValidateDate(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		wantErr bool
	}{
		{"valid date format", "2024-01-15", true},
		{"invalid format", "15.01.2024", false},
		{"empty date", "", false},
	}

	for _, testSuite := range tests {
		t.Run(testSuite.name, func(t *testing.T) {
			err := validateDate(&testSuite.date)
			if err != testSuite.wantErr {
				t.Errorf("validateDate() error = %v, wantErr %v", err, testSuite.wantErr)
			}
		})
	}
}
