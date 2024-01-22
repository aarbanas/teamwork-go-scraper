package main

import (
	"bufio"
	"fmt"
	"os"
)

type LogTimeMetaData struct {
	taskId      string
	hours       int16
	minutes     int16
	description string
}

func prepareData() (*LogTimeMetaData, error) {
	var taskId string
	var hours int16
	var minutes int16
	fmt.Println("Enter TaskId: ")
	_, taskError := fmt.Scan(&taskId)
	if taskError != nil {
		return nil, taskError
	}

	fmt.Println("Enter hours and minutes (separate with space): ")
	_, hoursMinutesError := fmt.Scan(&hours, &minutes)
	if taskError != nil {
		return nil, hoursMinutesError
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter description: ")
	description, _ := reader.ReadString('\n')

	// To remove the new line character at the end
	description = description[:len(description)-1]

	return &LogTimeMetaData{taskId: taskId, hours: hours, minutes: minutes, description: description}, nil
}

func logHours(startDate *string, endDate *string) {
	logMetadata, prepareDateErr := prepareData()
	if prepareDateErr != nil {
		fmt.Printf("Error: %s", prepareDateErr)
	}

	fmt.Printf("Log metadata %v", logMetadata.description)

	workingStartDate, err := convertStringFormatToDate(*startDate)
	workingEndDate, _err := convertStringFormatToDate(*endDate)
	if err != nil || _err != nil {
		os.Exit(1)
	}

	workDays := getWorkingDays(workingStartDate, workingEndDate)
	fmt.Println(*workDays)
}
